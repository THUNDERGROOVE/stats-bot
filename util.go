package main

import (
	"github.com/nlopes/slack"
	"log"
)

func getUsername(bot *slack.Slack, uid string) string {
	u, err := bot.GetUserInfo(uid)
	if err != nil {
		log.Printf("Error getting user info: %v", err.Error())
		return "Unknown user"
	}
	return u.Name
}
