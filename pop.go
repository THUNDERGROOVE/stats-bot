package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/census"
	"log"
	"strconv"
	"strings"
)

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
				strings.Repeat("*", v.VSPercent()))
			out += fmt.Sprintf(":tr: %%%v | %v\\",
				percPad(v.TRPercent()),
				strings.Repeat("*", v.TRPercent()))
			out += fmt.Sprintf(":nc: %%%v | %v \\",
				percPad(v.NCPercent()),
				strings.Repeat("*", v.NCPercent()))
			return out
		}
	}
	return "That server doesn't exist.  Really"
}

func StartPopGathering() {
	log.Printf("Starting to gather population information")
	USEvents = Census.NewEventStream()
	EUEvents = CensusEU.NewEventStream()

	USPop = Census.NewPopulationSet()
	EUPop = CensusEU.NewPopulationSet()

	sub := census.NewEventSubscription()
	sub.Worlds = []string{"all"}
	sub.Characters = []string{"all"}
	sub.EventNames = []string{"PlayerLogin", "PlayerLogout"}

	if err := USEvents.Subscribe(sub); err != nil {
		fmt.Printf("FAIL: Couldn't subscribe to events: [%v]\n", err.Error())
	}

	if err := EUEvents.Subscribe(sub); err != nil {
		fmt.Printf("FAIL: Couldn't subscribe to events: [%v]\n", err.Error())
	}
	parseEventsInto(Census, USEvents, USPop)
	parseEventsInto(CensusEU, EUEvents, EUPop)
}

func parseEventsInto(c *census.Census, events *census.EventStream, pop *census.PopulationSet) {
	log.Printf("Starting event parsing routine")
	go func() {
		for {
			select {
			case err := <-events.Err:
				if !strings.Contains(err.Error(), "EOF") {
					fmt.Printf("Events: error: %v\n", err.Error())
				}
			case <-events.Closed:
				fmt.Printf("Events: websocket closed\n")
				break
			case event := <-events.Events:
				switch event.Payload.EventName {
				case "PlayerLogin":
					ch, err := c.GetCharacterByID(event.Payload.CharacterID)
					if err != nil {
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
	}()
}
