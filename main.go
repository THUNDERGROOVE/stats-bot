package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/census"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

var Census *census.Census
var CensusEU *census.Census
var Dev bool

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic! [%v]", r)
			log.Printf("Restarting bot")
			StartBot()
		}
	}()

	if _, err := os.Stat(".git"); err == nil {
		log.Println("Git data found.  Running in development mode")
		Dev = true
	}

	if Dev {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	StartBot()
}

func StartBot() {
	log.Printf("Setting up slack bot")

	go StartHTTPServer()

	bot := slack.New(Config.Token)

	log.Printf("Setting up census client")
	Census = census.NewCensus("s:maximumtwang", "ps2ps4us:v2")

	CensusEU = census.NewCensus("s:maximumtwang", "ps2ps4eu:v2")

	StartPopGathering()

	t, err := bot.AuthTest()

	if err != nil {
		log.Printf("Error in auth test: ", err.Error())
		return
	}

	log.Printf("Auth: %v on team %v", t.User, t.Team)

	log.Printf("Starting slack events websocket\n")
	api, err := bot.StartRTM("", fmt.Sprintf("http://%v:8080/", "http://localhost/"))

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
				Dispatch(bot, sender, m)
			}
		}
	}

}

func sendMessages(api *slack.WS, sender chan slack.OutgoingMessage) {
	for {
		select {
		case msg := <-sender:
			if err := api.SendMessage(&msg); err != nil {
				log.Printf("Error sending message: %v", err.Error())
			}
		}
	}
}
