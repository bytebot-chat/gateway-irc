package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bbriggs/bytebot-irc/model"
	"github.com/go-redis/redis/v8"
	"github.com/satori/go.uuid"
)

// Flags and their default values.
var addr = flag.String("redis", "localhost:6379", "Redis server address")
var inbound = flag.String("inbound", "irc-inbound", "Pubsub queue to listen for new messages")
var outbound = flag.String("outbound", "irc", "Pubsub queue for sending messages outbound")

func main() {
	flag.Parse()
	ctx := context.Background()

	// We connect to the redis server
	rdb := redis.NewClient(&redis.Options{
		Addr: *addr,
		DB:   0,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		time.Sleep(3 * time.Second)
		err := rdb.Ping(ctx).Err()
		if err != nil {
			panic(err)
		}
	}

	// Reading the incoming messages into a channel
	topic := rdb.Subscribe(ctx, *inbound)
	channel := topic.Channel()

	// We iterate over each new message and reply to it appropriately.
	for msg := range channel {
		// Unmarshaling the message to be able to use it.
		// A message looks like this:
		//  type Message struct {
		//  	From     string
		//  	To       string
		//  	Content  string
		//  	Metadata Metadata // Source, Dest and ID.
		//  }

		m := &model.Message{}
		err := m.Unmarshal([]byte(msg.Payload))
		if err != nil {
			fmt.Println(err)
		}

		// Very simple logging
		fmt.Printf("%+v\n", m)

		// Here starts the app specific part.
		// This being a reaction, trigger, we want to reply
		// a different asciimoji to triggers.
		// Let's set the reaction with a switch.
		reactionContent := ""
		switch m.Content {
		case "!shrug":
			reactionContent = "¯\\_(ツ)_/¯"
		case "!lenny":
			reactionContent = "( ͡° ͜ʖ ͡°)"
		case "!tableflip":
			reactionContent = "(╯°□°)╯︵ ┻━┻"
		case "!tablefix":
			reactionContent = "┬─┬ノ( º _ ºノ)"
		}

		// If the message was a call to this app, we reply.
		if reactionContent != "" {
			reply(ctx, *m, rdb, reactionContent)
		}
	}
}

// This function creates a new message, and adds it to bytebot's inbound queue.
func reply(ctx context.Context, m model.Message, rdb *redis.Client, replyContent string) {
	// First, the IRC message, setting the destination to the original message's
	// source.
	if !strings.HasPrefix(m.To, "#") { // DMs go back to source, channel goes back to channel
		m.To = m.From
	}
	m.From = "" // No need to add a From field, Hellabot takes care of it for us.

	m.Content = replyContent // Setting the message to the one provided.

	// Now the message's metadata.
	m.Metadata.Dest = m.Metadata.Source                  // Setting the original message's source as a dest
	m.Metadata.Source = "reactions app"                  // The app's name.
	m.Metadata.ID = uuid.Must(uuid.NewV4(), *new(error)) // Creating a new ID that will be used for logging.

	// Finally, we marshal the message and push it to bytebot's inbound queue.
	stringMsg, _ := json.Marshal(m)
	rdb.Publish(ctx, *outbound, stringMsg)
}
