package config

// EmitterSecrets holds all Emitter information and channel keys
// 2D4B615064526755

var PubSubNATSHome = map[string]string{
	"host":    "tcp://10.0.0.22:8080",
	"license": "c1aVVz_sTmIi_FTcugWsjzTsQ4kJrslAAAAAAAAAAAI",
	"secret":  "1ZPCl42pPIyq6ZsZbaV4OUexWw97cZvf",
}

var PubSubNATSWork = map[string]string{
	"host":    "tcp://127.0.0.1:8080",
	"license": "6YUwQVizikSOTIMfKtAvrcW5hwFBLFL2AAAAAAAAAAI",
	"secret":  "BQQ1M7WIVGhWzjEilfV5ENHwYekj3T2z",
}

var PubSubMQTTHome = map[string]string{
	"mqtt.broker.host":     "10.0.0.58",
	"mqtt.broker.port":     "1883",
	"mqtt.broker.clientId": "gohome",
	"mqtt.broker.username": "gohome",
	"mqtt.broker.password": "gohome",
}

//var PubSubCfg = PubSubNATSWork
var PubSubCfg = PubSubMQTTHome

var InfluxSecretsHome = map[string]string{
	"host":     "http://10.0.0.22:8086",
	"username": "influxdb",
	"password": "password",
	"database": "gohome",
}

var InfluxSecretsWork = map[string]string{
	"host":     "http://localhost:8086",
	"username": "influxdb",
	"password": "password",
	"database": "gohome",
}

var InfluxSecrets = InfluxSecretsWork
