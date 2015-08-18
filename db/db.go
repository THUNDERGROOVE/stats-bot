// db sets up our database connection and exposes functions and methods to make
// working with data easier
package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"os"

	"log"
)

// TODO we need Update functions for our types

var DB *gorm.DB

func init() {
	if os.Getenv("DOCKER") == "" {
		if err := LoadSQLite("sqlite.db"); err != nil {
			log.Printf("[DB] error: %v\n", err.Error())
		}
	} else {
		log.Printf("[DB] Opening in Docker")
		if err := LoadSQLite("/data/sqlite.db"); err != nil {
			log.Fatalf("[DB] error: %v", err.Error())
		}
	}

	DB.CreateTable(Report{})
	DB.CreateTable(Outfit{})
}

// LoadSQLite loads the database as SQLite this should be used for local testing
func LoadSQLite(filename string) error {
	db, err := gorm.Open("sqlite3", filename)
	DB = &db
	return err
}
