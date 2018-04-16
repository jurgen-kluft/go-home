package config

type Config interface {
	FromJSON(json string) (Config, error)
	ToJSON() (string, error)
}
