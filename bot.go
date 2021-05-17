package main

import (
	hbot "github.com/whyrusleeping/hellabot"
)

func newBot(serv, nick *string, tls *bool, channels []string) (*hbot.Bot, error) {
	options := func(bot *hbot.Bot) {
		if *tls {
			bot.HijackSession = false
		} else {
			bot.HijackSession = true
		}
		bot.SSL = *tls
		bot.Channels = channels
	}

	return hbot.NewBot(*serv, *nick, options)
}
