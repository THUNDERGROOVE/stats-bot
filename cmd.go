package main

import (
	"bytes"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"strings"
	"text/template"
)

var Commands = make(map[string]*Cmd)

const lookup = `({{.Faction.Name.En}}) {{if .Outfit.Alias }}[{{.Outfit.Alias}}]{{end}} {{.Name.First}}
Kills: {{.GetKills}} Deaths: {{.GetDeaths}} KDR: {{.KDR}}
{{if .Outfit.Name}} Outfit: {{.Outfit.Name}} with {{.Outfit.MemberCount}} members {{end}}
Cached: {{.IsCached}}
`

var lookupTmpl *template.Template

func init() {
	var err error
	lookupTmpl = template.New("")
	lookupTmpl, err = lookupTmpl.Parse(lookup)
	if err != nil {
		log.Fatalf("Template failed to compile: [%v]", err.Error())
	}
	RegisterCommand("help", func(bot *slack.Slack, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
		Respond("Help? With what?", out, ev)
	})

	RegisterCommand("lookup", func(bot *slack.Slack, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
		args := strings.Split(ev.Text, " ")
		if len(args) <= 1 {
			Respond("Do you really expect me to lookup nothing?", out, ev)
			return
		}
		name := args[1]
		char, err := Census.QueryCharacterByExactName(name)
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				Respond("ERROR: The server closed the connection on us.", out, ev)
				return
			}
			log.Printf("Error getting character info: [%v]", err.Error())
			Respond(fmt.Sprintf("Error: %v",
				err.Error()), out, ev)
			return
		}
		buff := bytes.NewBuffer([]byte(""))
		if err := lookupTmpl.Execute(buff, char); err != nil {
			buff.WriteString("\nerror encountered" + err.Error())
		}
		Respond(buff.String(), out, ev)
	})
}

type Cmd struct {
	name    string
	handler func(*slack.Slack, chan slack.OutgoingMessage, *slack.MessageEvent)
}

func RegisterCommand(name string, handler func(*slack.Slack, chan slack.OutgoingMessage, *slack.MessageEvent)) {
	cmd := new(Cmd)
	cmd.name = name
	cmd.handler = handler
	Commands[name] = cmd
}

func Dispatch(bot *slack.Slack, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
	if bot.GetInfo().User.Id == ev.UserId {
		return
	}
	c := strings.ToLower(strings.Split(ev.Text, " ")[0])
	if v, ok := Commands[c]; ok {
		log.Printf("[Dispatch] Sending to %v", v.name)
		v.handler(bot, out, ev)
	} else {
		//@TODO: Handle undhandled commandd
		log.Printf("[Dispatch] Unhandled command %v", c)
	}
}

func Respond(s string, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
	lines := strings.Split(s, "\n")
	for _, v := range lines {
		o := slack.OutgoingMessage{}
		o.Text = v
		o.ChannelId = ev.ChannelId
		o.Type = ev.Type
		out <- o
	}
}
