package main

import (
	"bytes"
	"fmt"
	"github.com/THUNDERGROOVE/census"
	"github.com/nlopes/slack"
	"log"
	"strings"
	"text/template"
)

var Commands = make(map[string]*Cmd)

const helpText = `Hi.\
I'm stats-bot.  I have serveral commands!\
!lookup   <name>\
!lookupeu <name>\
!pop      <server>\
!popeu    <server>`

// Global is the struct given to any template parsed for responses
type Global struct {
	*census.Character
	Dev bool
}

var lookupTmpl *template.Template

func init() {

	lookupTmpl = template.Must(template.ParseFiles("lookup_template.tmpl"))

	// !help
	RegisterCommand("help",
		func(bot *slack.Slack,
			out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
			Respond(helpText, out, ev)
		})

	// !lookup
	RegisterCommand("lookup",
		func(bot *slack.Slack,
			out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
			LookupWith(Census, CensusEU, bot, out, ev)
		})

	// !lookupeu
	RegisterCommand("lookupeu",
		func(bot *slack.Slack,
			out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
			LookupWith(CensusEU, Census, bot, out, ev)
		})

	// !pop
	RegisterCommand("pop",
		func(bot *slack.Slack, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
			args := strings.Split(ev.Text, " ")
			if len(args) <= 1 {
				Respond("pop requires an argument you dingus", out, ev)
			}
			Respond(PopResp(USPop, args[1]), out, ev)
		})
	// !popeu
	RegisterCommand("popeu",
		func(bot *slack.Slack, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
			args := strings.Split(ev.Text, " ")
			if len(args) <= 1 {
				Respond("pop requires an argument you dingus", out, ev)
			}
			Respond(PopResp(EUPop, args[1]), out, ev)
		})
}

// LookupWith looks for a character given a several paramaters
func LookupWith(c *census.Census, fallbackc *census.Census, bot *slack.Slack, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
	args := strings.Split(ev.Text, " ")
	if len(args) <= 1 {
		Respond("Do you really expect me to lookup nothing?", out, ev)
		return
	}
	name := args[1]
	char, err := c.GetCharacterByName(name)
	if err != nil {
		if strings.Contains(err.Error(), "Get") {
			Respond("ERROR: The server closed the connection on us.  The API is either down or we are being rate-limited", out, ev)
			return
		} else if err != nil {
			Respond(fmt.Sprintf("Couldn't find the character '%v'", name), out, ev)
			return
		}
		log.Printf("Error getting character info: [%v] trying fallback", err.Error())

		char, err = fallbackc.GetCharacterByName(name)
		if err != nil {
			if strings.Contains(err.Error(), "Get") {
				Respond("ERROR: The server closed the connection on us.  The API is either down or we are being rate-limited", out, ev)
				return
			} else if err != nil {
				Respond(fmt.Sprintf("Couldn't find the character '%v'", name), out, ev)
				return
			}
		}
	}
	if char == nil {
		log.Printf("Query didn't return any error but character was nil")
		Respond("Query didn't return any error but character was nil", out, ev)
		return
	}
	if !census.CheckCache(census.CACHE_CHARACTER_EVENTS, "kills"+char.ID) {
		Respond("This character has no kills cache.  May take some time to process kill information!", out, ev)
	}
	buff := bytes.NewBuffer([]byte(""))
	if err := lookupTmpl.Execute(buff, Global{Character: char, Dev: Dev}); err != nil {
		buff.WriteString("\nerror encountered" + err.Error())
	}
	Respond(buff.String(), out, ev)
}

// Cmd is a command handler struct
type Cmd struct {
	name    string
	handler func(*slack.Slack, chan slack.OutgoingMessage, *slack.MessageEvent)
}

// RegisterCommand registers a command for the bot to dispatch
func RegisterCommand(name string, handler func(*slack.Slack, chan slack.OutgoingMessage, *slack.MessageEvent)) {
	cmd := new(Cmd)
	cmd.name = name
	cmd.handler = handler
	Commands[name] = cmd
}

// Dispatch sends a message to the bot
func Dispatch(bot *slack.Slack, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in dispatch")
		}
	}()

	if bot.GetInfo().User.Name == ev.User {
		return
	}
	c := strings.ToLower(strings.Split(ev.Text, " ")[0])
	if len(ev.Text) == 0 {
		log.Printf("Got blank message")
		return
	}
	if ev.Text[0] == '!' {

		if v, ok := Commands[strings.TrimLeft(c, "!")]; ok {

			//log.Printf("[Dispatch] Sending to %v", v.name)
			v.handler(bot, out, ev)

		} else {
			log.Printf("[Dispatch] Unhandled command '%v' ev: '%v'", c, ev.Text)
			Respond("I don't know what you want from me :( do !help?", out, ev)
		}
	}
}

// Respond is a helper function to send text responses to the slack server.
func Respond(s string, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
	//lines := strings.Split(s, "\\")
	text := strings.Replace(s, "\\", "\n", -1)
	o := slack.OutgoingMessage{}
	o.Text = text
	o.Channel = ev.Channel
	o.Type = ev.Type
	out <- o

}

func parseURL(url string) string {
	url = strings.Split(url, "//")[1]
	url = strings.Split(url, ".slack.com/")[0]
	return url
}

// No longer used?
func TKPercent(char *census.Character) float64 {
	kills := char.TeamKillsInLast(150)
	return (float64(kills) / 1000) * 100
}
