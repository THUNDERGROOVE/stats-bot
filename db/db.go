// db sets up our database connection and exposes functions and methods to make
// working with data easier
package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"

	"log"
)

// TODO we need Update functions for our types

var DB *gorm.DB

func init() {
	if err := LoadSQLite("/data/sqlite.db"); err != nil {
		log.Fatalf("[DB] error: %v", err.Error())
	}
	DB.CreateTable(Report{})
}

// LoadSQLite loads the database as SQLite this should be used for local testing
func LoadSQLite(filename) error {
	var err error
	DB, err = gorm.Open("sqlite3", filename)
	return err
}
