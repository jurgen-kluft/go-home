package mqtt

import (
	"fmt"
	"net/url"
)

// mqttconfig .
type mqttconfig struct {
	User     string
	Password string
	Broker   string
}

func parseConnStr(connurl string) (mqttconfig, error) {
	u, err := url.Parse(connurl)
	if err != nil {
		return mqttconfig{}, fmt.Errorf("buu")

	}
	host := u.Host

	if u.Port() == "" {
		host = fmt.Sprintf("%s:1883", host)
	}

	u2 := url.URL{
		Scheme: "tcp",
		Host:   host,
	}

	pw, _ := u.User.Password()
	return mqttconfig{
		User:     u.User.Username(),
		Password: pw,
		Broker:   u2.String(),
	}, nil

}
