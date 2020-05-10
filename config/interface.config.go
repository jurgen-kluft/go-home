package config

type Config interface {
	FromJSON(json []byte) error
	ToJSON() ([]byte, error)
}
