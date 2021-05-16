package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/bbriggs/bytebot-irc/model"
	"github.com/go-redis/redis/v8"
	"github.com/satori/go.uuid"
	hbot "github.com/whyrusleeping/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	ctx context.Context
	rdb *redis.Client
)

var serv = flag.String("server", "localhost:6667", "hostname and port for irc server to connect to")
var nick = flag.String("nick", "bytebot", "nickname for the bot")
var id = flag.String("id", "irc", "ID to use when publishing messages")
var inbound = flag.String("inbound", "irc-inbound", "Pubsub queue to publish inbound messages to")
var outbound = flag.String("outbound", *id, "Pubsub to subscribe to for sending outbound messages. Defaults to being equivalent to `id`")
var tls = flag.Bool("tls", false, "Use TLS when connecting to IRC server")

func main() {
	flag.Parse()

	rdb = rdbConnect()
	ctx = context.Background()

	irc, _ := newBot(serv, nick, tls)
	irc.AddTrigger(relayMessages)
	irc.Logger.SetHandler(log.StreamHandler(os.Stdout, log.JsonFormat()))
	go handleOutbound(*outbound, rdb, irc)
	irc.Run()
	fmt.Println("Bot shutting down.")
}

var relayMessages = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return true
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		msg := new(model.Message)
		msg.Metadata.ID = uuid.Must(uuid.NewV4(), *new(error))
		msg.Metadata.Source = *id
		msg.Message = m
		stringMsg, _ := json.Marshal(msg)
		rdb.Publish(ctx, *inbound, stringMsg)
		return false
	},
}

func handleOutbound(sub string, rdb *redis.Client, irc *hbot.Bot) {
	ctx := context.Background()
	topic := rdb.Subscribe(ctx, sub)
	channel := topic.Channel()
	for msg := range channel {
		m := &model.Message{}
		err := m.Unmarshal([]byte(msg.Payload))
		if err != nil {
			fmt.Println(err)
		}
		if m.Metadata.Dest == *id {
			irc.Msg(m.To, m.Content)
		}
	}
}
