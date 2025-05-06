package yeel

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Device struct {
	CacheControl    string
	Location        string
	ID              string
	Model           string // mono, color, stripe
	FirmwareVersion string
	Support         map[string]bool
	Power           string
	Flowing         int
	FlowParams      string

	Brightness      int
	ColorMode       int
	ColorTemprature int
	RGB             int
	Hue             int
	Saturation      int

	Name string
	// Type, Online, LQI, R, G, B, Level, Effect int

}

func (d *Device) LocationAddr() string {
	return strings.TrimPrefix(d.Location, "yeelight://")

}

func ParseDeviceFromHeader(header http.Header) (Device, error) {
	var d Device

	for k, v := range map[string]*string{
		"cache-control": &d.CacheControl,
		"id":            &d.ID,
		"location":      &d.Location,
		"fw_ver":        &d.FirmwareVersion,
		"name":          &d.Name,
		"model":         &d.Model,
		"power":         &d.Power,
	} {
		str := header.Get(k)
		*v = str
	}

	for k, v := range map[string]*int{
		"bright":     &d.Brightness,
		"rgb":        &d.RGB,
		"hue":        &d.Hue,
		"color_mode": &d.ColorMode,
		"sat":        &d.Saturation,
		"ct":         &d.ColorTemprature,
	} {
		str := header.Get(k)
		if str != "" {
			n, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				log.Fatal(str, err)
			}
			*v = int(n)
		}
	}
	{
		supportList := strings.Split(header.Get("support"), " ")
		supportMap := make(map[string]bool, len(supportList))
		for _, v := range supportList {
			supportMap[v] = true
		}
		d.Support = supportMap
	}
	// spew.Dump(header,d)
	return d, nil
}
