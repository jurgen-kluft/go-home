// To parse and unparse this JSON data, add this code to your project and do:
//
//    r, err := UnmarshalCaqiResponse(bytes)
//    bytes, err = r.Marshal()

package aqi

import "encoding/json"

func unmarshalCaqiResponse(data []byte) (CaqiResponse, error) {
	var r CaqiResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CaqiResponse) marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CaqiResponse struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	Aqi          int64         `json:"aqi"`
	Idx          int64         `json:"idx"`
	Attributions []Attribution `json:"attributions"`
	City         City          `json:"city"`
	Dominentpol  string        `json:"dominentpol"`
	Iaqi         Iaqi          `json:"iaqi"`
	Time         Time          `json:"time"`
}

type Attribution struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type City struct {
	Name string    `json:"name"`
	URL  string    `json:"url"`
	Geo  []float64 `json:"geo"`
}

type Iaqi struct {
	Co   Var `json:"co"`
	D    Var `json:"d"`
	H    Var `json:"h"`
	No2  Var `json:"no2"`
	O3   Var `json:"o3"`
	P    Var `json:"p"`
	Pm10 Var `json:"pm10"`
	Pm25 Var `json:"pm25"`
	So2  Var `json:"so2"`
	T    Var `json:"t"`
	W    Var `json:"w"`
	Wd   Var `json:"wd"`
}

type Var struct {
	V float64 `json:"v"`
}

type Time struct {
	S  string `json:"s"`
	Tz string `json:"tz"`
	V  int64  `json:"v"`
}
