package aqi

import (
	"bytes"
	"testing"
)

func TestCorrectJsonResponse(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.WriteString(testCorrectJsonResponse)

	cr, err := unmarshalCaqiResponse(buf.Bytes())
	if err != nil {
		t.Error("Cannot unmarshall correct json text (%s)", err)
	}
	if cr.Data.Aqi != 169 {
		t.Error("Json unmarshall results in wrong Aqi value")
	}
}

var testCorrectJsonResponse = `
{
	"status":"ok",
	"data": {
		"aqi":169,
		"idx":1437,
		"attributions":[
			{
				"url":"http://www.semc.gov.cn/",
				"name":"Shanghai Environment Monitoring Center(上海市环境监测中心)"
			},
			{
				"url":"http://113.108.142.147:20035/emcpublish/",
				"name":"China National Urban air quality real-time publishing platform (全国城市空气质量实时发布平台)"
			},
			{
				"url":"http://shanghai.usembassy-china.org.cn/airmonitor.html",
				"name":"U.S. Consulate Shanghai Air Quality Monitor"
			}
		],
		"city":	{
			"geo": [ 31.2047372, 121.4489017 ],
			"name":"Shanghai",
			"url":"http://aqicn.org/city/shanghai/"
		},
		"dominentpol": "pm25",
		"iaqi":{
			"co": {"v":16.3},
			"d": {"v":-4},
			"h": {"v":65},
			"no2": {"v":28.8},
			"o3":{"v":14.3},
			"p":{"v":1027},
			"pm10":{"v":84},
			"pm25":{"v":169},
			"so2":{"v":9.2},
			"t":{"v":2},
			"w":{"v":4},
			"wd":{"v":360}},
			"time":{
				"s":"2018-02-08 10:00:00",
				"tz":"+08:00",
				"v":1518084000
			}
		}
	}
`
