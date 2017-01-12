package netgear

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// netgear: Go API to retrieve devices attached to modern Netgear routers

const sessionID = "A7D88AE69687E58D9A00"

const soapActionLogin = "urn:NETGEAR-ROUTER:service:ParentalControl:1#Authenticate"
const soapLoginMessage = `\
<?xml version="1.0" encoding="utf-8" ?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
<SOAP-ENV:Header>
<SessionID xsi:type="xsd:string" xmlns:xsi="http://www.w3.org/1999/XMLSchema-instance">%s</SessionID>
</SOAP-ENV:Header>
<SOAP-ENV:Body>
<Authenticate>
  <NewUsername>%s</NewUsername>
  <NewPassword>%s</NewPassword>
</Authenticate>
</SOAP-ENV:Body>
</SOAP-ENV:Envelope>
`

const soapActionGetAttachedDevices = "urn:NETGEAR-ROUTER:service:DeviceInfo:1#GetAttachDevice"
const soapGetAttachedDevicesMesssage = `\
<?xml version="1.0" encoding="utf-8" standalone="no"?>
<SOAP-ENV:Envelope xmlns:SOAPSDK1="http://www.w3.org/2001/XMLSchema" xmlns:SOAPSDK2="http://www.w3.org/2001/XMLSchema-instance" xmlns:SOAPSDK3="http://schemas.xmlsoap.org/soap/encoding/" xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
<SOAP-ENV:Header>
<SessionID>%s</SessionID>
</SOAP-ENV:Header>
<SOAP-ENV:Body>
<M1:GetAttachDevice xmlns:M1="urn:NETGEAR-ROUTER:service:DeviceInfo:1">
</M1:GetAttachDevice>
</SOAP-ENV:Body>
</SOAP-ENV:Envelope>
`

// Router describes a modern Netgear router providing a SOAP interface at port 5000
type Router struct {
	host     string
	username string
	password string
	loggedIn bool
	regex    *regexp.Regexp
}

// IsLoggedIn returns true if the session has been authenticated against the
// Netgear Router or false otherwise.
func (netgear *Router) isLoggedIn() bool {
	return netgear.loggedIn
}

// Login authenticates the session against the Netgear router
// On success true and nil should be returned. Otherwise false and
// the related error are returned
func (netgear *Router) login() (bool, error) {
	message := fmt.Sprintf(soapLoginMessage, sessionID, netgear.username, netgear.password)

	resp, err := netgear.makeRequest(soapActionLogin, message)

	if strings.Contains(resp, "<ResponseCode>000</ResponseCode>") {
		netgear.loggedIn = true
	} else {
		netgear.loggedIn = false
	}
	return netgear.loggedIn, err
}

func (netgear *Router) getURL() string {
	return fmt.Sprintf("http://%s:5000/soap/server_sa/", netgear.host)
}

func (netgear *Router) makeRequest(action string, message string) (string, error) {
	client := &http.Client{}

	url := netgear.getURL()

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(message)))
	if err != nil {
		return "", err
	}
	req.Header.Add("SOAPAction", action)

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}

// Get queries the Netgear router for attached network
// devices and returns a list of them. If an error occures an empty list
// and the respective error is returned.
func (netgear *Router) Get(mac map[string]bool) error {
	message := fmt.Sprintf(soapGetAttachedDevicesMesssage, sessionID)
	resp, err := netgear.makeRequest(soapActionGetAttachedDevices, message)

	if strings.Contains(resp, "<ResponseCode>000</ResponseCode>") {
		re := netgear.regex.FindStringSubmatch(resp)
		if len(re) < 2 {
			err = fmt.Errorf("Invalid response code")
			return err
		}

		filteredDevicesStr := strings.Replace(re[1], "&lt;unknown&gt;", "unknown", -1)

		deviceStrs := strings.Split(filteredDevicesStr, "@")
		for _, deviceStr := range deviceStrs {
			fields := strings.Split(deviceStr, ";")
			if len(fields) >= 4 {
				macStr := fields[3]
				_, macExists := mac[macStr]
				if macExists {
					mac[macStr] = true
				}
			}
		}
	}
	return err
}

// newNetgearRouter returns a new and already initialized Netgear router instance
// However, the Netgear SOAP session has not been authenticated at this point.
// Use Login() to authenticate against the router
func New(host, username, password string) *Router {
	router := &Router{
		host:     host,
		username: username,
		password: password,
		regex:    regexp.MustCompile("<NewAttachDevice>(.*)</NewAttachDevice>"),
	}
	return router
}
