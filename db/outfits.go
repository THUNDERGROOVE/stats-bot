package db

import (
	"github.com/THUNDERGROOVE/census"
	"github.com/jinzhu/gorm"
)

// Outfit...
type Outfit struct {
	gorm.Model
	Name        string
	CID         string // CID is the ID given to it by census.
	LeaderSlack string // This is the username of the user in Slack
	// More?
}

func GetOutfit(CID string) *Outfit {
	outfit := new(Outfit)

	DB.Where(Outfit{CID: CID}).First(outfit)
	return outfit
}

func NewOutfit(Name, CID, LeaderSlack string, c *census.Census) {
	outfit := Outfit{
		Name:        Name,
		CID:         CID,
		LeaderSlack: LeaderSlack,
	}

	DB.Create(&outfit)
}
