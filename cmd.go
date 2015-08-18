// cmd.go contains all of our code to parse !commands they will be going away
// soon though so expect this file to shorten up or go away completely.  That
// is if I can get Slack to get back at me about /command timeouts.

package main

import (
	"bytes"
	"github.com/THUNDERGROOVE/census"
	"github.com/nlopes/slack"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"text/template"
)

type cmdType uint8

const (
	CMD_READY cmdType = iota
	CMD_DEV
)

const helpText = `Hi.\
I'm stats-bot.  I have serveral commands!\
!lookup   <name>\
!lookupeu <name>\
!pop      <server>\`

// TODO: Rip out after we convert over to all /commands
// May not happen after all :(
var Commands = make(map[string]*Cmd)
var lookupTmpl *template.Template

// Context is what's given to every command handler.  It should contain
// everything a command will need
type Context struct {
	Bot *slack.Slack
	Out chan slack.OutgoingMessage
	Ev  *slack.MessageEvent
}

func (c *Context) Respond(s string) {
	Respond(s, c.Out, c.Ev)
}

// Global is the struct given to any template parsed for responses
type Global struct {
	*census.Character
	Dev bool
}

func init() {
	log.Println("Registering main commands")
	lookupTmpl = parseTemplate("lookup_template.tmpl")

	// !help
	RegisterCommand("help", cmdHelp, CMD_READY)

	// !lookup
	RegisterCommand("lookup", cmdLookup, CMD_READY)

	// !lookupeu
	RegisterCommand("lookupeu", cmdLookupEU, CMD_READY)

	// !pop
	RegisterCommand("pop", cmdPop, CMD_READY)
}

func cmdHelp(ctx *Context) {
	ctx.Respond(helpText)
}

func cmdVersion(ctx *Context) {

}

func cmdLookup(ctx *Context) {
	LookupWith(Census, CensusEU, ctx)
}
func cmdLookupEU(ctx *Context) {
	LookupWith(CensusEU, Census, ctx)
}
func cmdPop(ctx *Context) {
	args := strings.Split(ctx.Ev.Text, " ")
	if len(args) <= 1 {
		ctx.Respond("pop requires an argument you dingus")
	}

	if v, ok := Worlds[strings.ToLower(args[1])]; ok {
		if v {
			ctx.Respond(PopResp(USPop, args[1]))
		} else {
			ctx.Respond(PopResp(EUPop, args[1]))
		}
	} else {
		ctx.Respond("I don't know about that server.  I'm sorry :(")
	}
}

func lookupStatsChar(c *census.Census, name string) (string, error) {
	char, err := c.GetCharacterByName(name)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString("")
	if err := lookupTmpl.Execute(buf, Global{Character: char, Dev: Dev}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// LookupWith looks for a character given a several paramaters
//
// @TODO: Just cleaned up a bit.  Anything else we can do?
func LookupWith(c *census.Census, fallbackc *census.Census, ctx *Context) {
	args := strings.Split(ctx.Ev.Text, " ")
	if len(args) <= 1 {
		ctx.Respond("Do you really expect me to lookup nothing?")
		return
	}

	var response string
	var err error

	name := args[1]

	response, err = lookupStatsChar(c, name)
	if err != nil {
		resp, err := lookupStatsChar(fallbackc, name)
		if err != nil {
			response = "The character wasn't found."
		}
		response = resp
	}
	ctx.Respond(response)
}

// Cmd is a command handler struct
type Cmd struct {
	name    string
	handler func(*Context)
}

// RegisterCommand registers a command for the bot to dispatch
func RegisterCommand(name string, handler func(*Context), state cmdType) {
	cmd := new(Cmd)
	cmd.name = name
	if state == CMD_DEV && !Dev {
		cmd.handler = notReadyYet
	} else {
		cmd.handler = handler
	}
	Commands[name] = cmd
}

func notReadyYet(ctx *Context) {
	ctx.Respond("That command isn't ready yet")

}

// Dispatch sends a message to the bot
func Dispatch(ctx *Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in dispatch")
			debug.PrintStack()
		}
	}()

	if ctx.Bot.GetInfo().User.Name == ctx.Ev.User {
		return
	}
	c := strings.ToLower(strings.Split(ctx.Ev.Text, " ")[0])
	if len(ctx.Ev.Text) == 0 {
		log.Printf("Got blank message")
		return
	}
	if ctx.Ev.Text[0] == '!' {

		if v, ok := Commands[strings.TrimLeft(c, "!")]; ok {

			//log.Printf("[Dispatch] Sending to %v", v.name)
			v.handler(ctx)

		} else {
			ctx.Respond("I don't know what you want from me :( do !help?")
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

func parseTemplate(filename string) *template.Template {
	// Default directory if we're in a Docker environment
	lookupName := filepath.Join("/assets", filename)

	// Sometimes it might just be in the current working directory
	if _, err := os.Stat(lookupName); err != nil {
		lookupName = filename
	}

	return template.Must(template.ParseFiles(lookupName))
}
