package main

import (
	hbot "github.com/whyrusleeping/hellabot"
)

func newBot(serv, nick *string, tls *bool) (*hbot.Bot, error) {
	options := func(bot *hbot.Bot) {
		if *tls {
			bot.HijackSession = false
		} else {
			bot.HijackSession = true
		}
		bot.SSL = *tls
		bot.Channels = []string{"#test"}
	}

	return hbot.NewBot(*serv, *nick, options)
}
