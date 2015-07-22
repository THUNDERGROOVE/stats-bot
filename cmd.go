package main

import (
	"bytes"
	"github.com/THUNDERGROOVE/census"
	"github.com/nlopes/slack"
	"log"
	"strings"
	"text/template"
)

var Commands = make(map[string]*Cmd)

const lookup = `({{.Faction.Name.En}}) {{if .Outfit.Alias }}[{{.Outfit.Alias}}]{{end}} {{.Name.First}} BR: {{.Battlerank.Rank}} :cert: {{.GetCerts}}
Kills: {{.GetKills}} Deaths: {{.GetDeaths}} KDR: {{.KDR}}
{{if .Outfit.Name}} Outfit: {{.Outfit.Name}} with {{.Outfit.MemberCount}} members {{end}}
Defended: {{.GetFacilitiesDefended}} Captured: {{.GetFacilitiesCaptured}}
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
		LookupWith(Census, CensusEU, bot, out, ev)
	})
	RegisterCommand("lookupeu", func(bot *slack.Slack, out chan slack.OutgoingMessage, ev *slack.MessageEvent) {
		LookupWith(CensusEU, Census, bot, out, ev)
	})
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
		log.Printf("Error getting character info: [%v] trying fallback", err.Error())

		char, err = fallbackc.QueryCharacterByExactName(name)
		if err != nil {
			if strings.Contains(err.Error(), "Get") {
				Respond("ERROR: The server closed the connection on us.  The API is either down or we are being rate-limited", out, ev)
				return
			}
		}
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

func TKPercent(char *census.Character) float64 {
	kills := char.TeamKillsInLast(150)
	return (float64(kills) / 1000) * 100
}
