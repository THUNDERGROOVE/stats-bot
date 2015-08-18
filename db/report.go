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
