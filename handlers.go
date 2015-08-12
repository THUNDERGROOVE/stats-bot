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

type CommandData struct {
	Token       string
	TeamID      string
	TeamDomain  string
	ChannelName string
	Command     string
	Text        string
}

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

func StartHTTPServer() {
	log.Printf("Starting command connection handler!\n")
	r := mux.NewRouter()

	r.HandleFunc("/pop", handlePop)

	err := http.ListenAndServe(":1339", handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		log.Printf("Why did I die? %v", err.Error())
	}
}

func handlePop(rw http.ResponseWriter, req *http.Request) {
	c := ParseCommandData(req)

	log.Printf("handlePop: Got server: %v", c.Text)

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
