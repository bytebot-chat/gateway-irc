package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/bytebot-chat/gateway-irc/model"
	"github.com/go-redis/redis/v8"
	"github.com/satori/go.uuid"
	hbot "github.com/whyrusleeping/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	ctx context.Context
	rdb *redis.Client

	channels  stringArrayFlags
	serv      = flag.String("server", "localhost:6667", "hostname and port for irc server to connect to")
	redisAddr = flag.String("redis", "localhost:6379", "Address and port of redis host")
	nick      = flag.String("nick", "bytebot", "nickname for the bot")
	id        = flag.String("id", "irc", "ID to use when publishing messages")
	inbound   = flag.String("inbound", "irc-inbound", "Pubsub queue to publish inbound messages to")
	outbound  = flag.String("outbound", *id, "Pubsub to subscribe to for sending outbound messages. Defaults to being equivalent to `id`")
	tls       = flag.Bool("tls", false, "Use TLS when connecting to IRC server")
)

func main() {

	flag.Var(&channels, "channel", "Channel to join at startup. May be repeated. Example: -channel=\\#foo -channel=\\#bar")

	flag.Parse()
	parseEnv()

	rdb = rdbConnect(*redisAddr)
	ctx = context.Background()

	log.Info("Connecting to IRC", "server", *serv, "nick", *nick, "tls", *tls)
	irc, _ := newBot(serv, nick, tls, channels)
	irc.AddTrigger(relayMessages)

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
		msg.From = m.From
		msg.To = m.To
		msg.Content = m.Content
		stringMsg, _ := json.Marshal(msg)
		log.Info("dispatching message", "topic", *inbound, "from", m.From, "to", m.To, "content", m.Content)
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
			log.Error("Failed to unmarshal message", "error", err.Error())
			fmt.Println(err)
		}
		if m.Metadata.Dest == *id {
			log.Info("dispatching message", "topic", sub, "from", m.From, "to", m.To, "content", m.Content)
			irc.Msg(m.To, m.Content)
		}
	}
}
