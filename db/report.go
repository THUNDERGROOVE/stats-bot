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

package db

import (
	"github.com/THUNDERGROOVE/census"
	"github.com/jinzhu/gorm"
)

// Report ...
type Report struct {
	gorm.Model
	Reporter       string
	Name           string
	PSNName        string
	AdditionalInfo string
	OutfitCID      string // The outfit the user was in at the time of the report
	Cleared        bool
}

// GetReport gets a report by the ID
func GetReport(ID int) *Report {
	report := new(Report)

	q := Report{}
	q.ID = uint(ID)

	DB.Where(q).First(report)
	return report
}

// NewReport creates a new report with the provided information
func NewReport(reporter, player, PSN, info string, c *census.Census) {
	report := Report{
		Reporter:       reporter,
		Name:           player,
		PSNName:        PSN,
		AdditionalInfo: info,
	}

	// TODO Figure out when to get our CID.  Here seams like a smart place?
	DB.Create(&report)
}

// ToggleClear toggles the Cleared value in the database.
func (r *Report) ToggleClear() {
	if r.Cleared {
		r.Cleared = false
		DB.Save(r)
	} else {
		r.Cleared = true
		DB.Save(r)
	}
}
