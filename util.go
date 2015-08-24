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
	"bytes"
	"github.com/THUNDERGROOVE/census"
	"github.com/THUNDERGROOVE/stats-bot/db"
	"github.com/nlopes/slack"
	"log"
	"strings"
)

// LookupWith looks for a character given a several paramaters
//
// TODO: Just cleaned up a bit.  Anything else we can do?
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
func lookupStatsChar(c *census.Census, name string) (string, error) {
	char, err := c.GetCharacterByName(name)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString("")

	// This is ugly and terrible and idc
	var cleared int
	var reported int
	reports := []db.Report{}
	db.DB.Where("name = ?", char.Name.First).Find(&reports)
	for _, v := range reports {
		if v.Cleared {
			cleared += 1
		} else {
			reported += 1
		}
	}

	if err := lookupTmpl.Execute(buf, Global{
		Character:     char,
		Dev:           Dev,
		TimesReported: reported,
		TimesCleared:  cleared,
		Reports:       reports}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getUsername(bot *slack.Client, uid string) string {
	u, err := bot.GetUserInfo(uid)
	if err != nil {
		log.Printf("Error getting user info: %v", err.Error())
		return "Unknown user"
	}
	return u.Name
}

// parseURL gets the slack domain from your team given the URL to it
func parseURL(url string) string {
	url = strings.Split(url, "//")[1]
	url = strings.Split(url, ".slack.com/")[0]
	return url
}
