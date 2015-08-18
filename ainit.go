package main

import (
	"log"
	"os"
)

// This file exits because inits are done alphabetically.  We want Dev to be
// set as early as possible.  If this continues being a proble make a function
// isDev and assign it to Dev at global scope:
// var Dev = isDev()
func init() {
	if _, err := os.Stat(".git"); err == nil {
		log.Println("Git data found.  Running in development mode")
		Dev = true
	}
}
