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

// stats-bot is a bot written for the Slack platform that takes advantage of the
// census package to produce statistics for players.
//
// Additionally it contains other functionality to help maintain the community
package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/census"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var (
	Census   *census.Census
	CensusEU *census.Census

	Dev     bool
	Commit  string
	Version string
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if Dev {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	StartBot()
}

func StartBot() {
	log.Printf("stats-bot: v%v#%v", Version, Commit)

	go StartHTTPServer() // Has to run in a Goroutine.  Blocks
	bot := slack.New(Config.Token)

	Census = census.NewCensus("s:maximumtwang", "ps2ps4us:v2")
	CensusEU = census.NewCensus("s:maximumtwang", "ps2ps4eu:v2")

	StartPopGathering()

	t, err := bot.AuthTest()
	if err != nil {
		log.Printf("Error in auth test: [%v]", err.Error())
		return
	}

	log.Printf("Auth: %v on team %v", t.User, t.Team)

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
				Dispatch(&Context{Bot: bot, Ev: m, Out: sender})
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
