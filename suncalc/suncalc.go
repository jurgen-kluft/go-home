package main

// sun calculations are based on http://aa.quae.nl/en/reken/zonpositie.html formulas

import (
	"math"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
)

const (
	pi  = math.Pi
	rad = pi / 180.0
)

var sin = math.Sin
var cos = math.Cos
var tan = math.Tan
var asin = math.Asin
var atan = math.Atan2
var acos = math.Acos

// date/time constants and conversions
const (
	daySeconds = float64(60.0 * 60.0 * 24.0)
	j1970      = float64(2440588)
	j2000      = float64(2451545)
)

func toJulian(t time.Time) float64 {
	return (float64(t.Unix()) / daySeconds) - 0.5 + j1970
}

func fromJulian(j float64) time.Time {
	return time.Unix(int64((j+0.5-j1970)*daySeconds), 0)
}

func toDays(t time.Time) float64 {
	return toJulian(t) - j2000
}

// general calculations for position

var e = rad * 23.4397 // obliquity of the Earth

func rightAscension(l float64, b float64) float64 {
	return atan(sin(l)*cos(e)-tan(b)*sin(e), cos(l))
}

func declination(l float64, b float64) float64 {
	return asin(sin(b)*cos(e) + cos(b)*sin(e)*sin(l))
}

func azimuth(H float64, phi float64, dec float64) float64 {
	return atan(sin(H), cos(H)*sin(phi)-tan(dec)*cos(phi))
}
func altitude(H float64, phi float64, dec float64) float64 {
	return asin(sin(phi)*sin(dec) + cos(phi)*cos(dec)*cos(H))
}

func siderealTime(d float64, lw float64) float64 {
	return rad*(280.16+360.9856235*d) - lw
}

func astroRefraction(h float64) float64 {
	if h < 0 { // the following formula works for positive altitudes only.
		h = 0.0 // if h = -0.08901179 a div/0 would occur.
	}

	// formula 16.4 of "Astronomical Algorithms" 2nd edition by Jean Meeus (Willmann-Bell, Richmond) 1998.
	// 1.02 / tan(h + 10.26 / (h + 5.10)) h in degrees, result in arc minutes -> converted to rad:
	return 0.0002967 / tan(h+0.00312536/(h+0.08901179))
}

// general sun calculations

func solarMeanAnomaly(d float64) float64 {
	return rad * (357.5291 + 0.98560028*d)
}

func eclipticLongitude(M float64) float64 {

	// equation of center
	var C = rad * (1.9148*sin(M) + 0.02*sin(2*M) + 0.0003*sin(3*M))

	// perihelion of the Earth
	var P = rad * 102.9372

	return M + C + P + pi
}

func sunCoords(d float64) (dec float64, ra float64) {
	M := solarMeanAnomaly(d)
	L := eclipticLongitude(M)

	dec = declination(L, 0)
	ra = rightAscension(L, 0)
	return
}

// calculates sun position for a given date and latitude/longitude
func getPosition(date time.Time, lat float64, lng float64) (outAzimuth float64, outAltitude float64) {
	lw := rad * -lng
	phi := rad * lat
	d := toDays(date)

	cra, cdec := sunCoords(d)
	H := siderealTime(d, lw) - cra

	outAzimuth = azimuth(H, phi, cdec)
	outAltitude = altitude(H, phi, cdec)
	return
}

// calculations for sun times

const (
	j0 = 0.0009
)

func julianCycle(d float64, lw float64) float64 {
	return math.Floor(d - j0 - lw/(2*pi))
}
func approxTransit(Ht float64, lw float64, n float64) float64 {
	return j0 + (Ht+lw)/(2*pi) + n
}
func solarTransitJ(ds float64, M float64, L float64) float64 {
	return j2000 + ds + 0.0053*sin(M) - 0.0069*sin(2*L)
}
func hourAngle(h float64, phi float64, d float64) float64 {
	return acos((sin(h) - sin(phi)*sin(d)) / (cos(phi) * cos(d)))
}

// returns set time for the given sun altitude
func getSetJ(h float64, lw float64, phi float64, dec float64, n float64, M float64, L float64) float64 {

	var w = hourAngle(h, phi, dec)
	var a = approxTransit(w, lw, n)
	return solarTransitJ(a, M, L)
}

// Cmoment is a time period (from - to) with a title and description
// For example "sunrise", "top edge of the sun appears on the horizon until it is fully visible"
// with a 'begin' (time.Time) and an 'end' (time.Time).
// The whole array spans a full day starting at midnight (0:00) the morning, noon, evening until midnight.
type Cmoment struct {
	title string
	descr string
	start time.Time
	end   time.Time
}

// GetMoments returns the current day of Cmoment items (see comment on Cmoment)
func (s *instance) getMoments(date time.Time, lat float64, lng float64) (result []Cmoment) {
	lw := rad * -lng
	phi := rad * lat

	d := toDays(date)
	n := julianCycle(d, lw)
	ds := approxTransit(0, lw, n)

	M := solarMeanAnomaly(ds)
	L := eclipticLongitude(M)
	dec := declination(L, 0)

	Jnoon := solarTransitJ(ds, M, L)
	//fmt.Printf("Suncalc, Noon %v\n", fromJulian(Jnoon))

	mtimes := map[string]time.Time{}
	mtimes["today.begin"] = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	mtimes["today.end"] = time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())

	mtimes["night.darkest"] = fromJulian(Jnoon - 0.5)

	noon := fromJulian(Jnoon)
	mtimes["sun.noon.begin"] = hoursLater(noon, -0.15)
	mtimes["sun.noon.end"] = hoursLater(noon, +0.15)
	for _, a := range s.config.Anglecfg {
		Jset := getSetJ(a.Angle*rad, lw, phi, dec, n, M, L)
		Jrise := Jnoon - (Jset - Jnoon)
		mtimes[a.Rise] = fromJulian(Jrise)
		mtimes[a.Set] = fromJulian(Jset)
	}

	// type Cmoment struct {
	// 	title string
	// 	descr string
	// 	from  time.Time
	// 	to    time.Time
	// }

	moments := []Cmoment{}
	for _, m := range s.config.Moments {
		t0 := mtimes[m.Begin]
		t1 := mtimes[m.End]

		moment := Cmoment{}
		moment.descr = m.Descr
		moment.title = m.Title
		moment.start = t0
		moment.end = t1
		moments = append(moments, moment)
	}

	return moments
}

// moon calculations, based on http://aa.quae.nl/en/reken/hemelpositie.html formulas

func moonCoords(d float64) (float64, float64, float64) { // geocentric ecliptic coordinates of the moon

	L := rad * (218.316 + 13.176396*d) // ecliptic longitude
	M := rad * (134.963 + 13.064993*d) // mean anomaly
	F := rad * (93.272 + 13.229350*d)  // mean distance

	l := L + rad*6.289*sin(M)   // longitude
	b := rad * 5.128 * sin(F)   // latitude
	dt := 385001 - 20905*cos(M) // distance to the moon in km

	ra := rightAscension(l, b)
	dec := declination(l, b)
	dist := dt

	return ra, dec, dist
}

type moonpos struct {
	azimuth          float64
	altitude         float64
	distance         float64
	parallacticAngle float64
}

func getMoonPosition(date time.Time, lat float64, lng float64) moonpos {

	lw := rad * -lng
	phi := rad * lat
	d := toDays(date)

	cra, cdec, cdist := moonCoords(d)
	H := siderealTime(d, lw) - cra
	h := altitude(H, phi, cdec)
	// formula 14.1 of "Astronomical Algorithms" 2nd edition by Jean Meeus (Willmann-Bell, Richmond) 1998.
	pa := atan(sin(H), tan(phi)*cos(cdec)-sin(cdec)*cos(H))

	h = h + astroRefraction(h) // altitude correction for refraction

	p := moonpos{}
	p.azimuth = azimuth(H, phi, cdec)
	p.altitude = h
	p.distance = cdist
	p.parallacticAngle = pa
	return p
}

// calculations for illumination parameters of the moon,
// based on http://idlastro.gsfc.nasa.gov/ftp/pro/astro/mphase.pro formulas and
// Chapter 48 of "Astronomical Algorithms" 2nd edition by Jean Meeus (Willmann-Bell, Richmond) 1998.

func getMoonIllumination(date time.Time) (fraction float64, phase float64, angle float64) {

	d := toDays(date)
	sra, sdec := sunCoords(d)
	mra, mdec, mdist := moonCoords(d)

	sdist := 149598000.0 // distance from Earth to Sun in km

	phi := acos(sin(sdec)*sin(mdec) + cos(sdec)*cos(mdec)*cos(sra-mra))
	inc := atan(sdist*sin(phi), mdist-sdist*cos(phi))
	angle = atan(cos(sdec)*sin(sra-mra), sin(sdec)*cos(mdec)-cos(sdec)*sin(mdec)*cos(sra-mra))

	fraction = (1.0 + cos(inc)) / 2.0
	if angle < 0.0 {
		phase = 0.5 + 0.5*inc*-1.0/pi
	} else {
		phase = 0.5 + 0.5*inc*1.0/pi
	}
	return
}

func hoursLater(date time.Time, h float64) time.Time {
	return time.Unix(date.Unix()+int64(h*float64(daySeconds)/24.0), 0)
}

// calculations for moon rise/set times are based on http://www.stargazing.net/kepler/moonrise.html article

func getMoonTimes(date time.Time, lat float64, lng float64, inUTC bool) (moonrise bool, moonriseTime time.Time, moonset bool, moonsetTime time.Time, alwaysUp bool, alwaysDown bool) {
	t := date

	hc := 0.133 * rad
	mp := getMoonPosition(t, lat, lng)
	h0 := mp.altitude - hc

	// go in 2-hour chunks, each time seeing if a 3-point quadratic curve crosses zero (which means rise or set)
	i := float64(1)
	brise := false
	bset := false
	rise := 0.0
	set := 0.0
	x1 := 0.0
	x2 := 0.0
	ye := 0.0
	for i <= 24 {
		h1 := getMoonPosition(hoursLater(t, i), lat, lng).altitude - hc
		h2 := getMoonPosition(hoursLater(t, i+1), lat, lng).altitude - hc

		a := (h0+h2)/2 - h1
		b := (h2 - h0) / 2
		xe := -b / (2 * a)
		ye := (a*xe+b)*xe + h1
		d := b*b - 4*a*h1
		roots := 0

		if d >= 0 {
			dx := math.Sqrt(d) / (math.Abs(a) * 2)
			x1 = xe - dx
			x2 = xe + dx
			if math.Abs(x1) <= 1 {
				roots++
			}
			if math.Abs(x2) <= 1 {
				roots++
			}
			if x1 < -1 {
				x1 = x2
			}
		}

		if roots == 1 {
			if h0 < 0 {
				brise = true
				rise = i + x1
			} else {
				bset = true
				set = i + x1
			}

		} else if roots == 2 {
			if ye < 0 {
				brise = true
				bset = true
				rise = i + x2
				set = i + x1
			} else {
				brise = true
				bset = true
				rise = i + x1
				set = i + x2
			}
		}

		if brise && bset {
			break
		}

		h0 = h2
		i += 2
	}

	moonrise = brise
	if brise {
		moonriseTime = hoursLater(t, rise)
	} else {
		moonriseTime = time.Now()
	}

	moonset = bset
	if bset {
		moonsetTime = hoursLater(t, set)
	} else {
		moonsetTime = time.Now()
	}

	alwaysUp = false
	alwaysDown = false
	if !brise && !bset {
		if ye > 0 {
			alwaysUp = true
		} else {
			alwaysDown = true
		}
	}
	return
}

type instance struct {
	config *config.SuncalcConfig
}

func new() *instance {
	s := &instance{}
	return s
}

func (s *instance) initialize(jsondata []byte) error {
	config, err := config.SuncalcConfigFromJSON(jsondata)
	if err == nil {
		s.config = config
	}
	return err
}

func (s *instance) buildJSONMessage() ([]byte, error) {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.Local)

	lat := s.config.Geo.Latitude
	lng := s.config.Geo.Longitude

	moments := s.getMoments(now, lat, lng)

	sunstate := config.NewSensorState("state.sensor.sun", "sun")
	for _, m := range moments {
		sunstate.AddTimeWndAttr(m.title, m.start, m.end)
	}

	_, moonPhase, _ := getMoonIllumination(now)
	sunstate.AddFloatAttr("moon.illumination", moonPhase)

	jsonbytes, err := sunstate.ToJSON()
	return jsonbytes, err
}

func main() {
	suncalc := new()

	register := []string{"state/sensor/sun/", "config/request/"}
	subscribe := []string{"config/suncalc/"}

	m := microservice.New("suncalc")
	m.RegisterAndSubscribe(register, subscribe)

	m.RegisterHandler("config/suncalc/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo(m.Name, "received configuration")
		err := suncalc.initialize(msg)
		if err != nil {
			m.Logger.LogError(m.Name, err.Error())
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if tickCount%30 == 0 { // every 30 seconds
			if suncalc.config == nil {
				m.Pubsub.PublishStr("config/request/", m.Name)
			}
		}
		if tickCount%300 == 0 { // every 5 minutes
			if suncalc.config != nil {
				jsonbytes, err := suncalc.buildJSONMessage()
				if err == nil {
					m.Pubsub.Publish("state/sensor/sun/", jsonbytes)
				} else {
					m.Logger.LogError(m.Name, err.Error())
				}
			}
		}
		tickCount++
		return true
	})

	m.Loop()
}
