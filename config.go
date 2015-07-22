package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type _conf struct {
	Token string `json:"token"`
}

var Config *_conf

func init() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Printf("Error opening config file: %v", err.Error())
		return
	}
	Config = new(_conf)
	if err := json.Unmarshal(data, Config); err != nil {
		log.Printf("Error unmarshaling config file: %v", err.Error())
	}
}
