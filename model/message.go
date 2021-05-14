package model

import (
	"github.com/satori/go.uuid"
	hbot "github.com/whyrusleeping/hellabot"
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
