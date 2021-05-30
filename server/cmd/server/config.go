package main

import (
	"encoding/json"
	"io/ioutil"
)

type InfluxDBConfig struct {
	Address      string `json:"address"`
	Organization string `json:"org"`
	AuthToken    string `json:"token"`
	Bucket       string `json:"bucket"`
	Measurement  string `json:"measurement"`
}

func loadConfig(file string, v interface{}) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}
