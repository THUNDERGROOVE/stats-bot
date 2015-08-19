// The MIT License (MIT)
//
// Copyright (c) 2015 Nick Powell
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
