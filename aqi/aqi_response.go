package main

// To parse and unparse this JSON data, add this code to your project and do:
//
//    r, err := UnmarshalCaqiResponse(bytes)
//    bytes, err = r.Marshal()

import "encoding/json"

func unmarshalCaqiResponse(data []byte) (CaqiResponse, error) {
	var r CaqiResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CaqiResponse) marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Example response
// {
// 	"status": "ok",
// 	"data": {
// 	  "aqi": 87,
// 	  "idx": 1437,
// 	  "attributions": [
// 		{
// 		  "url": "http://106.37.208.233:20035/emcpublish/",
// 		  "name": "China National Urban air quality real-time publishing platform (全国城市空气质量实时发布平台)"
// 		},
// 		{
// 		  "url": "https://china.usembassy-china.org.cn/embassy-consulates/shanghai/air-quality-monitor-stateair/",
// 		  "name": "U.S. Consulate Shanghai Air Quality Monitor"
// 		},
// 		{
// 		  "url": "https://sthj.sh.gov.cn/",
// 		  "name": "Shanghai Environment Monitoring Center(上海市环境监测中心)"
// 		},
// 		{
// 		  "url": "https://waqi.info/",
// 		  "name": "World Air Quality Index Project"
// 		}
// 	  ],
// 	  "city": {
// 		"geo": [
// 		  31.2047372,
// 		  121.4489017
// 		],
// 		"name": "Shanghai (上海)",
// 		"url": "https://aqicn.org/city/shanghai"
// 	  },
// 	  "dominentpol": "pm25",
// 	  "iaqi": {
// 		"co": {
// 		  "v": 5.5
// 		},
// 		"h": {
// 		  "v": 21.3
// 		},
// 		"no2": {
// 		  "v": 11.9
// 		},
// 		"p": {
// 		  "v": 1007.7
// 		},
// 		"pm10": {
// 		  "v": 57
// 		},
// 		"pm25": {
// 		  "v": 87
// 		},
// 		"so2": {
// 		  "v": 3.1
// 		},
// 		"t": {
// 		  "v": 29.3
// 		},
// 		"w": {
// 		  "v": 0.1
// 		}
// 	  },
// 	  "time": {
// 		"s": "2020-05-20 17:00:00",
// 		"tz": "+08:00",
// 		"v": 1589994000
// 	  },
// 	  "debug": {
// 		"sync": "2020-05-20T18:29:53+09:00"
// 	  }
// 	}
//}

// CaqiResponse is the response we get from the HTTP request
type CaqiResponse struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

// Data is the actual interesting AQI data
type Data struct {
	Aqi          int64         `json:"aqi"`
	Idx          int64         `json:"idx"`
	Attributions []Attribution `json:"attributions"`
	City         City          `json:"city"`
	Dominentpol  string        `json:"dominentpol"`
	Iaqi         Iaqi          `json:"iaqi"`
	Time         Time          `json:"time"`
}

// Attribution shows where the AQI data is actually coming from
type Attribution struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// The City we have requested the AQI of
type City struct {
	Name string    `json:"name"`
	URL  string    `json:"url"`
	Geo  []float64 `json:"geo"`
}

// Iaqi is the detailed AQI information
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

// Var is a float64 value
type Var struct {
	V float64 `json:"v"`
}

// Time holds the time information of the AQI data
type Time struct {
	S  string `json:"s"`
	Tz string `json:"tz"`
	V  int64  `json:"v"`
}
