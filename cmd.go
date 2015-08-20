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

// cmd.go contains all of our code to parse !commands they will be going away
// soon though so expect this file to shorten up or go away completely.  That
// is if I can get Slack to get back at me about /command timeouts.

package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/census"
	"github.com/nlopes/slack"
	"log"
	"runtime/debug"
	"strings"
	"text/template"
)

type cmdType uint8

const (
	CMD_READY cmdType = iota
	CMD_ADMIN
	CMD_DEV
)

const helpText = `Hi.\
I'm stats-bot.  I have serveral commands!\
!lookup   <name>\
!lookupeu <name>\
!pop      <server>\

!report <name> <additional info>
!reportpsn <name> <psn> <additional info>

!clearreport <id>
!deletereport <id>

!searchreport <name>
!searchreportpsn <psn name>
!searchreportoutfit <outfit name>

!isadmin
`

// TODO: Rip out after we convert over to all /commands
// May not happen after all :(
var Commands = make(map[string]*Cmd)

// Global is the struct given to the lookup template for responses
type Global struct {
	*census.Character
	Dev bool
}

var lookupTmpl *template.Template

// Cmd is a command handler struct
type Cmd struct {
	name       string
	handler    func(*Context)
	adminCheck bool
}

// RegisterCommand registers a command for the bot to dispatch
func RegisterCommand(name string, handler func(*Context), state cmdType) {
	cmd := new(Cmd)
	cmd.name = name
	switch state {
	case CMD_DEV:
		if !Dev {
			cmd.handler = notReadyYet
		} else {
			cmd.handler = handler
		}
	case CMD_READY:
		cmd.handler = handler
	case CMD_ADMIN:
		cmd.adminCheck = true
		cmd.handler = handler

	}
	Commands[name] = cmd
}

// Dispatch sends a message to the bot
func Dispatch(ctx *Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in dispatch")
			ctx.Respond("Something in that command scared me :( could almost say I was paniced")
			debug.PrintStack()
		}
	}()

	c := strings.ToLower(strings.Split(ctx.Ev.Text, " ")[0])
	if len(ctx.Ev.Text) == 0 {
		log.Printf("Got blank message")
		return
	}
	if ctx.Ev.Text[0] == '!' {
		if v, ok := Commands[strings.TrimLeft(c, "!")]; ok {
			if !v.adminCheck || isAdmin(ctx) {
				v.handler(ctx)
			} else {
				ctx.Respond("You do not have permission to do that.  Sorry :(")
			}
		} else {
			ctx.Respond("I don't know what you want from me :( do !help?")
		}
	}
}

// Respond is a helper function to send text responses to the slack server.
func Respond(s string, rtm *slack.RTM, ev *slack.MessageEvent) {
	//lines := strings.Split(s, "\\")
	text := strings.Replace(s, "\\", "\n", -1)
	if rtm == nil {
		log.Printf("RTM nil?")
	}
	out := rtm.NewOutgoingMessage(text, ev.Channel)
	rtm.SendMessage(out)
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

	RegisterCommand("version", cmdVersion, CMD_READY)
}

func cmdHelp(ctx *Context) {
	ctx.Respond(helpText)
}

func cmdVersion(ctx *Context) {
	ctx.Respond(fmt.Sprintf("stats-bot: v%v running commit %v", Version, Commit))
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
