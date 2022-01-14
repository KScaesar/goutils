package message

import (
	"context"

	"github.com/Min-Feng/goutils"
)

type Topic = string

type Message interface {
	Topic() Topic
}

type Event interface {
	Message
	Event()
}

type CommandQuery interface {
	Message
	CommandQuery()
}

type Response interface {
	Response()
}

type request struct {
	ctx     context.Context
	payload Message

	blocking bool
	ready    chan Topic
	reply    chan result
}

func (req *request) complete(topic Topic) {
	req.ready <- topic
}

type result struct {
	payload Response
	err     error
}

type InfoBase struct {
	CorrelationID string       `json:"correlation_id"`
	MessageID     string       `json:"message_id"`
	MessageTopic  Topic        `json:"topic"`
	OccurredAt    goutils.Time `json:"occurred_at"`
}

func (b *InfoBase) Topic() Topic {
	return b.MessageTopic
}
