package config

// EmitterSecrets holds all Emitter information and channel keys
// 2D4B615064526755

var PubSubHome = map[string]string{
	"host":    "tcp://10.0.0.22:8080",
	"license": "c1aVVz_sTmIi_FTcugWsjzTsQ4kJrslAAAAAAAAAAAI",
	"secret":  "1ZPCl42pPIyq6ZsZbaV4OUexWw97cZvf",
}

var PubSubWork = map[string]string{
	"host":    "tcp://127.0.0.1:8080",
	"license": "6YUwQVizikSOTIMfKtAvrcW5hwFBLFL2AAAAAAAAAAI",
	"secret":  "BQQ1M7WIVGhWzjEilfV5ENHwYekj3T2z",
}

var PubSubCfg = PubSubHome

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
