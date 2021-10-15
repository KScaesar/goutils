package message

import "context"

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
