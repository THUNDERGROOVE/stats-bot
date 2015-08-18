package db

import (
	"github.com/THUNDERGROOVE/census"
	"log"
	"testing"
)

// Really not real tests.  Just here to verify gorm works as I thought

func TestNewOutfit(t *testing.T) {
	c := census.NewCensus("maximumtwang", "ps2ps4us:v2")

	NewOutfit("Testing Outfit", "r2877654", "THUNDERGROOVE", c)
}

func TestGetOutfit(t *testing.T) {
	o := GetOutfit("12877654")
	if o.Name != "Testing Outfit" {
		t.Fatalf("Failed to get outfit :(")
	}
}
