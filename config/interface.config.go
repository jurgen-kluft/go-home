package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
)

type Config interface {
	FromJSON(json []byte) error
	ToJSON() ([]byte, error)
}

var errUnsupportedFileExtension = errors.New("unsupported file extension")

func loadJSON(path string, v interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func LoadFile(path string, v interface{}) error {
	var err error
	switch ext := filepath.Ext(path); ext {
	case ".json":
		err = loadJSON(path, v)
	default:
		err = errUnsupportedFileExtension
	}
	return err
}
