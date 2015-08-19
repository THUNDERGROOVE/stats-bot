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

package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/census"
	"strconv"
	"strings"
)

// true = US
var Worlds = map[string]bool{
	"ceres":    false,
	"lithcorp": false,
	"xelas":    true,
	"palos":    true,
	"crux":     true,
	"genudine": true,
	"searhus":  true,
}

var (
	USEvents *census.EventStream
	EUEvents *census.EventStream

	USPop *census.PopulationSet
	EUPop *census.PopulationSet
)

// percPad adds padding spaces to make sure the length is always three
func percPad(perc int) string {
	s := strconv.Itoa(perc)
	switch len(s) {
	case 1:
		return "  " + s
	case 2:
		return " " + s
	}
	return s
}

func PopResp(pop *census.PopulationSet, server string) string {
	for k, v := range pop.Servers {
		if strings.ToLower(k) == strings.ToLower(server) {
			var out string

			out += fmt.Sprintf(":vanu: %%%v | %v\\",
				percPad(v.VSPercent()),
				strings.Repeat("*", v.VSPercent()/2))
			out += fmt.Sprintf(":tr: %%%v | %v\\",
				percPad(v.TRPercent()),
				strings.Repeat("*", v.TRPercent()/2))
			out += fmt.Sprintf(":nc: %%%v | %v \\",
				percPad(v.NCPercent()),
				strings.Repeat("*", v.NCPercent()/2))
			return out
		}
	}
	return "That server doesn't exist.  Really"
}

// TODO: See if moving this to an init would break anything
// would look a lot cleaner
func StartPopGathering() {
	go DoUSPop()
	go DoEUPop()
}

func DoUSPop() {
	USPop = Census.NewPopulationSet()
	for {
		USEvents = Census.NewEventStream()

		sub := census.NewEventSubscription()
		sub.Worlds = []string{"all"}
		sub.Characters = []string{"all"}
		sub.EventNames = []string{"PlayerLogin", "PlayerLogout"}

		if err := USEvents.Subscribe(sub); err != nil {
			fmt.Printf("FAIL: Couldn't subscribe to events: [%v]\n", err.Error())
		}

		parseEventsInto(Census, USEvents, USPop)
	}
}

func DoEUPop() {
	EUPop = CensusEU.NewPopulationSet()
	for {
		EUEvents = CensusEU.NewEventStream()

		sub := census.NewEventSubscription()
		sub.Worlds = []string{"all"}
		sub.Characters = []string{"all"}
		sub.EventNames = []string{"PlayerLogin", "PlayerLogout"}

		if err := EUEvents.Subscribe(sub); err != nil {
			fmt.Printf("FAIL: Couldn't subscribe to events: [%v]\n", err.Error())
		}

		parseEventsInto(CensusEU, EUEvents, EUPop)
	}
}

func parseEventsInto(c *census.Census, events *census.EventStream, pop *census.PopulationSet) {
loop:
	for {
		select {
		case err := <-events.Err:
			if strings.Contains(err.Error(), census.ErrCharDoesNotExist.Error()) {
				fmt.Printf("Events: error: %v\n", err.Error())
			}
		case <-events.Closed:
			fmt.Printf("Events: websocket closed\n")
			break loop
		case event := <-events.Events:
			switch event.Payload.EventName {
			case "PlayerLogin":
				ch, err := c.GetCharacterByID(event.Payload.CharacterID)
				if err != nil {
					if err == census.ErrCharDoesNotExist {
						continue
					}
					fmt.Printf("Events: ERROR: Failed to get character from ID: '%v' [%v]\n",
						event.Payload.CharacterID, err.Error())
					continue
				}
				server := c.GetServerByID(event.Payload.WorldID)
				pop.PlayerLogin(server.Name.En, ch.FactionID)
			case "PlayerLogout":
				ch, err := c.GetCharacterByID(event.Payload.CharacterID)
				if err != nil {
					fmt.Printf("ERROR: Failed to get character from ID: '%v' [%v]\n",
						event.Payload.CharacterID, err.Error())
					continue
				}
				server := c.GetServerByID(event.Payload.WorldID)
				pop.PlayerLogin(server.Name.En, ch.FactionID)
			}
		}
	}
	events.Close()
}
