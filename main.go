package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/census"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var Census *census.Census
var CensusEU *census.Census
var Dev bool

func getExternalIP() string {
	return "0.0.0.0"
	resp, err := http.Get("http://myexternalip.com/raw")

	if err != nil {
		log.Printf("Error getting external IP [%v]", err.Error())
	}
	data, _ := ioutil.ReadAll(resp.Body)
	o := string(data)
	o = strings.Replace(o, "\n", "", -1)
	return o
}
func main() {
	log.SetFlags(log.Lshortfile)
	/*defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic! [%v]", r)
		}
	}()*/

	if _, err := os.Stat(".git"); err == nil {
		log.Println("Git data found.  Running in development mode")
		Dev = true
	}
	log.Printf("Setting up slack bot")
	bot := slack.New(Config.Token)

	log.Printf("Setting up census client")
	Census = census.NewCensus("s:maximumtwang", "ps2ps4us:v2")

	CensusEU = census.NewCensus("s:maximumtwang", "ps2ps4eu:v2")

	//bot.SetDebug(true)

	t, err := bot.AuthTest()

	if err != nil {
		log.Printf("Error in auth test")
	}

	log.Printf("Auth: %v on team %v", t.User, t.Team)

	ip := getExternalIP()
	log.Printf("Starting RTM @[%v]", ip)
	api, err := bot.StartRTM("", fmt.Sprintf("http://%v:8080/", ip))

	if err != nil {
		log.Printf("Error setting up RTM [%v]", err.Error())
	}

	sender := make(chan slack.OutgoingMessage)
	receiver := make(chan slack.SlackEvent)

	go api.HandleIncomingEvents(receiver)
	go api.Keepalive(20 * time.Second)

	go sendMessages(api, sender)

	for {
		select {
		case msg := <-receiver:
			switch m := msg.Data.(type) {
			case *slack.MessageEvent:
				//log.Printf("Got message: %v -> %v ", getUsername(bot, m.UserId), m.Text)

				Dispatch(bot, sender, m)
			}
		}
	}
}

func sendMessages(api *slack.SlackWS, sender chan slack.OutgoingMessage) {
	for {
		select {
		case msg := <-sender:
			if err := api.SendMessage(&msg); err != nil {
				log.Printf("Error sending message: %v", err.Error())
			}
		}
	}
}
