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
	"log"
	"strconv"
	"strings"
	"text/template"

	"github.com/THUNDERGROOVE/census"
	"github.com/THUNDERGROOVE/stats-bot/db"
)

var searchTmpl *template.Template

func init() {
	searchTmpl = parseTemplate("search.tmpl")

	RegisterCommand("reportpsn", cmdReportPSN, CMD_READY)
	RegisterCommand("report", cmdReport, CMD_READY)

	RegisterCommand("clearreport", cmdClearReport, CMD_READY)
	RegisterCommand("deletereport", cmdDeleteReport, CMD_READY)

	RegisterCommand("searchreport", cmdSearchReports, CMD_ADMIN)
	RegisterCommand("searchreportpsn", cmdSearchReportsPSN, CMD_ADMIN)
	RegisterCommand("searchreportoutfit", cmdSearchReportsOutfit, CMD_ADMIN)

	RegisterCommand("allreports", cmdAllReports, CMD_ADMIN)
	RegisterCommand("myreports", cmdMyReports, CMD_READY)

	RegisterCommand("isadmin", cmdIsAdmin, CMD_READY)
}

func cmdAllReports(ctx *Context) {
	reports := []*db.Report{}
	db.DB.Where("cleared = 0").Find(&reports)
	g := map[string]interface{}{}
	g["Reports"] = reports
	g["Search"] = "cleared = 0"
	ctx.RenderTemplate(searchTmpl, g)
}

func cmdReport(ctx *Context) {
	args := strings.Split(ctx.Ev.Text, " ")
	if len(args) < 2 {
		ctx.Respond("report <name> <additional information>")
	}
	name := args[1]
	info := strings.Join(args[2:], " ")

	report(name, "not specified", info, ctx)
}

func cmdReportPSN(ctx *Context) {
	args := strings.Split(ctx.Ev.Text, " ")
	if len(args) < 3 {
		ctx.Respond("report <name> <psn> <additional information>")
	}
	name := args[1]
	psn := args[2]
	info := strings.Join(args[3:], " ")

	report(name, psn, info, ctx)
}

func cmdClearReport(ctx *Context) {
	args := strings.Split(ctx.Ev.Text, " ")
	if len(args) < 2 {
		ctx.Respond("clearreport <id>")
		return
	}
	sid := args[1]
	id, err := strconv.Atoi(sid)
	if err != nil {
		ctx.Respond("couldn't parse ID")
		return
	}

	r := db.GetReport(id)
	if r == nil {
		ctx.Respond("That report doesn't exist")
		return
	}

	if !isAdmin(ctx) && !(r.Reporter == ctx.Ev.User) {
		ctx.Respond("You do not have permission to do that")
		return
	}

	r.ToggleClear()
	if r.Cleared {
		ctx.Respond("The issue was marked as resolved")
	} else {
		ctx.Respond("The issue was re-opened")
	}
}

func cmdDeleteReport(ctx *Context) {
	args := strings.Split(ctx.Ev.Text, " ")
	if len(args) < 2 {
		ctx.Respond("deletereport <id>")
		return
	}
	sid := args[1]
	id, err := strconv.Atoi(sid)
	if err != nil {
		ctx.Respond("couldn't parse ID")
		return
	}

	r := db.GetReport(id)
	if r == nil {
		ctx.Respond("That report doesn't exist")
	}

	if !isAdmin(ctx) && !(r.Reporter == ctx.Ev.User) {
		log.Printf("! %v == %v", r.Reporter, ctx.Ev.User)
		ctx.Respond("You do not have permission to do that")
		return
	}

	db.DB.Delete(r)
	ctx.Respond("That report was successfully deleted")
}

func cmdSearchReportsPSN(ctx *Context) {
	args := strings.Split(ctx.Ev.Text, " ")
	search := args[1]

	reports := []*db.Report{}
	g := map[string]interface{}{}

	db.DB.Where("psn_name LIKE ?", "%"+search+"%").Find(&reports)
	g["Reports"] = reports

	ctx.RenderTemplate(searchTmpl, g)
}

func cmdSearchReportsOutfit(ctx *Context) {
	args := strings.Split(ctx.Ev.Text, " ")
	search := strings.Join(args[1:], " ")
	log.Printf("search: '%v'", search)

	var err error
	var outfit *census.Outfit
	if outfit, err = Census.GetOutfitByName(search); err != nil {
		err = nil
		if outfit, err = CensusEU.GetOutfitByName(search); err != nil {
			ctx.Respond("The outfit you're looking for doesn't exist.")
			return
		}
	}

	reports := []*db.Report{}
	g := map[string]interface{}{}
	db.DB.Where("outfit_cid = ?", outfit.ID).Find(&reports)
	g["Reports"] = reports
	g["Search"] = search
	ctx.RenderTemplate(searchTmpl, g)
}

func cmdSearchReports(ctx *Context) {
	args := strings.Split(ctx.Ev.Text, " ")
	search := args[1]

	reports := []*db.Report{}
	g := map[string]interface{}{}

	db.DB.Where("name like ?", "%"+search+"%").Find(&reports)
	g["Search"] = search
	g["Reports"] = reports

	ctx.RenderTemplate(searchTmpl, g)
}

func cmdMyReports(ctx *Context) {

	reports := []*db.Report{}
	g := map[string]interface{}{}

	db.DB.Where("reporter = ?", ctx.Ev.User).Find(&reports)

	g["Reports"] = reports

	ctx.RenderTemplate(searchTmpl, g)

}

func cmdIsAdmin(ctx *Context) {
	if isAdmin(ctx) {
		ctx.Respond("You are an admin")
	} else {
		ctx.Respond("You are not an admin")
	}
}

func isAdmin(ctx *Context) bool {
	u, err := ctx.RTM.GetUserInfo(ctx.Ev.User)
	if err != nil {
		log.Printf("isAdmin: %v", err.Error())
		return false
	}
	return (u.IsAdmin || u.IsOwner)
}

// Helpers

func report(name, psn, info string, ctx *Context) {
	var char *census.Character
	var err error
	if strings.Contains(name, ":") {
		region := strings.Split(name, ":")[0]
		name = strings.Split(name, ":")[1]

		switch strings.ToLower(region) {
		case "eu":
			char, err = CensusEU.GetCharacterByName(name)
		case "us":
			char, err = Census.GetCharacterByName(name)
		default:
			ctx.Respond("Unknown region code")
			return
		}
	} else {

		char, err = Census.GetCharacterByName(name)

		if err != nil {
			err = nil
			char, err = CensusEU.GetCharacterByName(name)
			if err != nil {
				ctx.Respond("That name didn't exist in either US or EU")
				return
			}
		} else {
			_, err := CensusEU.GetCharacterByName(name)
			if err == nil {
				ctx.Respond("The given character name matches on US and EU\nPlease specify with us:name or eu:name")
				return
			}
		}
	}
	if char == nil {
		ctx.Respond("That name didn't exist in either US or EU")
		return
	}
	if err := db.NewReport(ctx.Ev.User, char.Name.First, psn, info, char.Parent); err != nil {
		ctx.Respond("An internal error has occured: " + err.Error())
		return
	}
	ctx.Respond("Thank you for your report")
}
