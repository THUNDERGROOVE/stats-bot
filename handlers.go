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
