package deconz

import (
	"strconv"
	"testing"

	"github.com/jurgen-kluft/go-home/conbee/deconz/event"
)

const smokeDetectorNoFireEventPayload = `{	"e": "changed",	"id": "5",	"r": "sensors",	"state": {	  "fire": false,	  "lastupdated": "2018-03-13T19:46:03",	  "lowbattery": false,	  "tampered": false	},	"t": "event"  }`

type testLookup struct {
}

func (t *testLookup) LookupSensor(i int) (*Device, error) {
	return &Device{Name: "Test Sensor", Type: "ZHAFire"}, nil
}

func (t *testLookup) SupportsResource(_type string) bool {
	return true
}

func (t *testLookup) LookupType(i int) (string, error) {
	return "ZHAFire", nil
}

type testReader struct {
}

func (t testReader) ReadEvent() (*event.Event, error) {
	d := event.Decoder{TypeStore: &testLookup{}}
	return d.Parse([]byte(smokeDetectorNoFireEventPayload))
}
func (t testReader) Dial() error {
	return nil
}
func (t testReader) Close() error {
	return nil
}
func TestSensorEventReader(t *testing.T) {

	r := DeviceEventReader{reader: testReader{}}
	channel := make(chan *DeviceEvent)
	err := r.Start(channel)
	if err != nil {
		t.Fail()
	}
	e := <-channel
	if strconv.Itoa(e.Event.ID) != "5" {
		t.Fail()
	}
	fields, err := e.Fields()
	if err != nil {
		t.Logf(err.Error())
		t.FailNow()
	}

	if fields["fire"] != false {
		t.Fail()
	}

}
