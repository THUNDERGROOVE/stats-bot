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
	"github.com/THUNDERGROOVE/census"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
)

// CommandData represents the data sent to our handlers from Slack
type CommandData struct {
	Token       string
	TeamID      string
	TeamDomain  string
	ChannelName string
	Command     string
	Text        string
}

// ParseCommandData gets our data from the request.
//
// @TODO: Make this work for GET requests as well?
func ParseCommandData(req *http.Request) *CommandData {
	c := new(CommandData)
	req.ParseForm()
	c.Token = req.FormValue("token")
	c.TeamID = req.FormValue("team_id")
	c.TeamDomain = req.FormValue("team_domain")
	c.ChannelName = req.FormValue("channel")
	c.Command = req.FormValue("command")
	c.Text = req.FormValue("text")
	return c
}

// StartHTTPServer starts an http server with handlers for all of statsbot's
// commands
func StartHTTPServer() {
	log.Printf("Starting command connection handler!\n")
	r := mux.NewRouter()

	r.HandleFunc("/pop", handlePop)
	r.HandleFunc("/lookup", handleLookup)

	err := http.ListenAndServe(":1339", handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		log.Printf("Why did I die? %v", err.Error())
	}
}

func handleLookup(rw http.ResponseWriter, req *http.Request) {
	c := ParseCommandData(req)

	switch c.Command {
	case "/lookup":
		out, err := lookupStatsChar(Census, c.Text)
		// @TODO: Refactor this?
		if err != nil {
			rw.Write([]byte(err.Error()))
		} else {
			rw.Write([]byte(out))
		}
	case "/lookupeu":
		out, err := lookupStatsChar(CensusEU, c.Text)
		if err != nil {
			rw.Write([]byte(err.Error()))
		} else {
			rw.Write([]byte(out))
		}
	default:
		log.Printf("lookup handler called with wrong command?: %v\n", c.Command)
		rw.Write([]byte("The command given wasn't sent correctly"))
	}
}

func handlePop(rw http.ResponseWriter, req *http.Request) {
	c := ParseCommandData(req)

	var pop *census.PopulationSet

	if v, ok := Worlds[strings.ToLower(c.Text)]; ok {
		if v {
			pop = USPop
		} else {
			pop = EUPop
		}
		rw.Write([]byte(strings.Replace(PopResp(pop, c.Text), "\\", "\n", -1)))
	} else {
		rw.Write([]byte("I don't know about that server.  I'm sorry :("))
	}
}
