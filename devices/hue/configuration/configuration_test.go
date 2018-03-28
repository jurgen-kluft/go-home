package configuration

import (
	"os"
	"testing"
)

var (
	testUsername       string
	testHostname       string
	testDeleteUsername string
	testConfiguration  *Configuration
)

func init() {
	testUsername = os.Getenv("HUE_TEST_USERNAME")
	testHostname = os.Getenv("HUE_TEST_HOSTNAME")
	testHostname = "10.0.0.27"
	testUsername = "a6vdEfxKE2RD6VoLl4SHLUeM5hCkywRJnnMPfjX4"
	testDeleteUsername = ""
	testConfiguration = New(testHostname, testUsername)
}

func TestCreateUser(t *testing.T) {
	testApplicationName := "go-home"
	testDeviceType := "home"
	apiResponse, err := testConfiguration.CreateUser(testApplicationName, testDeviceType)
	if err != nil {
		t.Log("TestCreateUser Error: ", err)
		t.Fail()
	} else {
		t.Log(apiResponse[0].Success["username"].(string))
	}
}

func TestGetFullState(t *testing.T) {
	_, err := testConfiguration.GetFullState(testUsername)
	if err != nil {
		t.Log("TestGetFullState Error: ", err)
		t.Fail()
	}
}

func TestGetConfiguration(t *testing.T) {
	_, err := testConfiguration.GetConfiguration(testUsername)
	if err != nil {
		t.Log("TestGetConfiguration Error: ", err)
		t.Fail()
	}
}

func TestDeleteUser(t *testing.T) {
	_, err := testConfiguration.DeleteUser(testUsername, testDeleteUsername)
	if err != nil {
		t.Log("TestDeleteUser Error: ", err)
		t.Fail()
	}
}
