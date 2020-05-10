package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
)

var errUnsupportedConfig = errors.New("unsupported config file type")

func loadJSON(path string, v interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func loadFile(path string, v interface{}) error {
	var err error
	switch ext := filepath.Ext(path); ext {
	case ".json":
		err = loadJSON(path, v)
	}
	return err
}

type configuration struct {
	HTTP struct {
		Enabled    bool   `json:"enabled"`
		ListenAddr string `json:"listen"`
	} `json:"http"`
	LogDirectory string `json:"log_directory"`
}

func loadConfig(path string) (*configuration, error) {
	var c configuration
	err := loadFile(path, &c)
	return &c, err
}
