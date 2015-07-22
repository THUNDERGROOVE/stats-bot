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

const lookup = `({{.Faction.Name.En}}) {{if .Outfit.Alias }}[{{.Outfit.Alias}}]{{end}} {{.Name.First}}
Kills: {{.GetKills}} Deaths: {{.GetDeaths}} KDR: {{.KDR}}
{{if .Outfit.Name}} Outfit: {{.Outfit.Name}} with {{.Outfit.MemberCount}} members {{end}}
{{if .Dev}} Cached: {{.IsCached}}{{end}}
`

const helpText = `Hi.  I'm stats-bot.  You can ask me to '!lookup <name>' or '!lookupeu <name>'`

type Global struct {
	*census.Character
	Dev bool
}

var lookupTmpl *template.Template

func init() {
	var err error
	lookupTmpl = template.New("")
	lookupTmpl, err = lookupTmpl.Parse(lookup)
	if err != nil {
		log.Fatalf("Template failed to compile: [%v]", err.Error())
	}
	RegisterCommand("help", func(bot *slack.Slack, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
		Respond(helpText, out, ev)
	})

	RegisterCommand("lookup", func(bot *slack.Slack, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
		LookupWith(Census, CensusEU,bot, out, ev)
	})
	RegisterCommand("lookupeu", func(bot *slack.Slack, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
		LookupWith(CensusEU,Census, bot, out, ev)
	})
	/*
		RegisterCommand("invite", func(bot *slack.Slack, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
			log.Printf("invite triggered")
			args := strings.Split(ev.Text, " ")
			if len(args) <= 4 {
				Respond("usage: invite <email> <firstName> <lastName> <channel>", out, ev)
				return
			}

			email := args[1]

			if strings.Contains(email, "|") {
				email = strings.Split(email, "|")[1]
				email = strings.Replace(email, ">", "", -1)
			}

			firstName := args[2]
			lastName := args[3]
			channel := args[4]
			var cid string
			t, _ := bot.AuthTest()
			log.Printf("Checking channels")
			for _, v := range bot.GetInfo().Channels {
				if v.Name == channel {
					cid = v.Id
				} else {
					log.Printf("Channel wasn't '%v': [%v]", channel, v.Name)
				}
			}

			if cid == "" {
				Respond("Channel wasn't resolved", out, ev)
				return
			}

			log.Printf("inviting user")
			if err := bot.InviteUser(parseURL(t.Url), cid, firstName, lastName, email); err != nil {
				log.Printf("Error: %v", err.Error())
				Respond("Error: "+err.Error(), out, ev)
				return
			}

			Respond("User successfully invited", out, ev)
		})
	*/
}

func LookupWith(c *census.Census, fallbackc *census.Census, bot *slack.Slack, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
	args := strings.Split(ev.Text, " ")
	if len(args) <= 1 {
		Respond("Do you really expect me to lookup nothing?", out, ev)
		return
	}
	name := args[1]
	char, err := c.QueryCharacterByExactName(name)
	if err != nil {
		if strings.Contains(err.Error(), "Get") {
			Respond("ERROR: The server closed the connection on us.  The API is either down or we are being rate-limited", out, ev)
			return
		}
		log.Printf("Error getting character info: [%v]", err.Error())
		Respond(fmt.Sprintf("Error: %v",
			err.Error()), out, ev)
		return
	}
	buff := bytes.NewBuffer([]byte(""))
	if err := lookupTmpl.Execute(buff, Global{Character: char, Dev: Dev}); err != nil {
		buff.WriteString("\nerror encountered" + err.Error())
	}
	Respond(buff.String(), out, ev)
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
	if ev.Text[0] == '!' {

		if v, ok := Commands[strings.TrimLeft(c, "!")]; ok {

			//log.Printf("[Dispatch] Sending to %v", v.name)
			v.handler(bot, out, ev)

		} else {
			log.Printf("[Dispatch] Unhandled command '%v' ev: '%v'", c, ev.Text)
			Respond("!I don't know what you want from me :( do 'help'?", out, ev)
		}
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

func parseURL(url string) string {
	url = strings.Split(url, "//")[1]
	url = strings.Split(url, ".slack.com/")[0]
	return url
}
