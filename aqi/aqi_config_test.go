package aqi

import (
	"bytes"
	"testing"
)

func TestCorrectJsonConfig(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.WriteString(testCorrectJsonConfig)

	cr, err := unmarshalcaqi(buf.Bytes())
	if err != nil {
		t.Error("Cannot unmarshall correct json configuration", err)
	}
	if cr.City != "@1437" {
		t.Error("Json config unmarshall results in wrong city value")
	}
	if cr.Token != "this is a correct token" {
		t.Error("Json config unmarshall results in wrong token value")
	}
	if cr.URL != "https://api.waqi.info/feed/${CITY}/?token=${TOKEN}" {
		t.Error("Json config unmarshall results in wrong URL value")
	}
}

var testCorrectJsonConfig = `
{
    "token": "this is a correct token",
    "city": "@1437",
    "url": "https://api.waqi.info/feed/${CITY}/?token=${TOKEN}"
}
`
