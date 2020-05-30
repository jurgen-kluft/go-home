package config

import (
	"bytes"
	"testing"
)

func TestCorrectJsonAqiConfig(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.WriteString(testCorrectJSONAqiConfig)

	cr, err := AqiConfigFromJSON(buf.Bytes())
	if err != nil {
		t.Error("Cannot unmarshall correct json configuration", err)
	}
	if cr.City != "@1437" {
		t.Error("Json config unmarshall results in wrong city value")
	}
	if cr.Token.String != "this is a correct token" {
		t.Error("Json config unmarshall results in wrong token value")
	}
	if cr.URL != "https://api.waqi.info/feed/${CITY}/?token=${TOKEN}" {
		t.Error("Json config unmarshall results in wrong URL value")
	}
}

var testCorrectJSONAqiConfig = `
{
    "token": "8AiNYmJTWzzcIm0g-e0ZUhwGjYBxogDV_PyGFIDPjBzE1fSdPm6U",
    "city": "@1437",
    "url": "https://api.waqi.info/feed/${CITY}/?token=${TOKEN}"
}
`
