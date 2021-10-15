package message

import "context"

type Publisher interface {
	BlockingPublishEvent(ctx context.Context, events ...Event) error
	NonBlockingPublishEvent(ctx context.Context, events ...Event)
}
