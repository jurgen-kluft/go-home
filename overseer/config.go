/*
   Copyright 2014 Nick Saika

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

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
