package main

import (
	hbot "github.com/whyrusleeping/hellabot"
)

func newBot(serv, nick *string) (*hbot.Bot, error) {
	hijackSession := func(bot *hbot.Bot) {
		bot.HijackSession = true
	}

	channels := func(bot *hbot.Bot) {
		bot.Channels = []string{"#test"}
	}

	return hbot.NewBot(*serv, *nick, hijackSession, channels)
}
