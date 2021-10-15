package message

import "context"

type EventHandlerFunc func(context.Context, Event) error

type CommandQueryHandlerFunc func(context.Context, CommandQuery) (Response, error)

type eventSubscription struct {
	topic   Topic
	handler EventHandlerFunc
}

type commandQuerySubscription struct {
	topic   Topic
	handler CommandQueryHandlerFunc
}
