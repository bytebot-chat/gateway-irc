package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/satori/go.uuid"
	hbot "github.com/whyrusleeping/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	ctx context.Context
	rdb *redis.Client
)

type Message struct {
	*hbot.Message
	Metadata Metadata
}

type Metadata struct {
	Source string
	Dest   string
	ID     uuid.UUID
}

var serv = flag.String("server", "localhost:6667", "hostname and port for irc server to connect to")
var nick = flag.String("nick", "bytebot", "nickname for the bot")
var id = flag.String("id", "irc", "ID to use when publishing messages")
var inbound = flag.String("inbound", "irc-inbound", "Pubsub queue to publish inbound messages to")
var outbound = flag.String("outbound", *id, "Pubsub to subscribe to for sending outbound messages. Defaults to being equivalent to `id`")

func main() {
	flag.Parse()

	rdb = rdbConnect()
	ctx = context.Background()

	irc, _ := newBot(serv, nick)
	irc.AddTrigger(sayInfoMessage)
	irc.Logger.SetHandler(log.StreamHandler(os.Stdout, log.JsonFormat()))

	irc.Run()
	fmt.Println("Bot shutting down.")
}

var sayInfoMessage = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return true
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		msg := new(Message)
		msg.Metadata.ID = uuid.Must(uuid.NewV4(), *new(error))
		msg.Metadata.Source = *id
		msg.Message = m
		stringMsg, _ := json.Marshal(msg)
		rdb.Publish(ctx, *inbound, stringMsg)
		return false
	},
}
