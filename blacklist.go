package main

import (
	"bytes"
	"fmt"
	"github.com/THUNDERGROOVE/stats-bot/db"
	"log"
	"text/template"

	"strings"
)

var searchTmpl *template.Template

func init() {
	log.Println("Registering blacklist commands")
	searchTmpl = parseTemplate("search.tmpl")

	RegisterCommand("reportpsn", cmdReportPSN, CMD_DEV)
	RegisterCommand("searchreport", cmdSearchReports, CMD_DEV)
}

func cmdSearchReports(ctx *Context) {
	args := strings.Split(ctx.Ev.Text, " ")
	search := args[1]

	reports := []*db.Report{}
	g := map[string]interface{}{}

	db.DB.Where("name like ?", fmt.Sprintf("%%%v%%", search)).Find(&reports)
	g["Reports"] = reports
	buf := bytes.NewBufferString("")
	searchTmpl.Execute(buf, g)
	ctx.Respond(buf.String())
}

func cmdReportPSN(ctx *Context) {
	args := strings.Split(ctx.Ev.Text, " ")
	if len(args) < 3 {
		ctx.Respond("report <name> <psn> <additional information>")
	}
	name := args[1]
	psn := args[2]

	info := strings.Join(args[len(args)-3:], " ")

	// Lookup Character to know what region they play on

	char, err := Census.GetCharacterByName(name)

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
			ctx.Respond("The given character name matches on US and EU")
		}
	}
	db.NewReport(ctx.Ev.Name, char.Name.First, psn, info, char.Parent)
	ctx.Respond(fmt.Sprintf("Reported: %v for %v.", char.Name.First, info))
}

func cmdIsAdmin(ctx *Context) {
	if isAdmin(ctx) {
		ctx.Respond("You are an admin")
	} else {
		ctx.Respond("You are not an admin")
	}
}

func isAdmin(ctx *Context) bool {
	u, err := ctx.Bot.GetUserInfo(ctx.Ev.Username)
	if err != nil {
		return false
	}
	return u.IsAdmin
}
