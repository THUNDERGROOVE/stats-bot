package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// not idiomatic but I don't want it in my namespace switch to anonymous struct?
type _conf struct {
	Token string `json:"token"`
}

var Config *_conf

func init() {
	Config = new(_conf)
	if _, err := os.Stat("config.json"); err == nil {
		data, err := ioutil.ReadFile("config.json")
		if err != nil {
			log.Printf("Error opening config file: %v", err.Error())
			return
		}
		if err := json.Unmarshal(data, Config); err != nil {
			log.Printf("Error unmarshaling config file: %v", err.Error())
		}
	} else {
		tok := os.Getenv("slack_token")
		if tok == "" {
			log.Printf("Failed to get token from config AND from env\n")
		} else {
			Config.Token = tok
		}
	}
}
