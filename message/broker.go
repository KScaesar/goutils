package message

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/KScaesar/goutils/errors"
)

func NewBroker() *Broker {
	b := &Broker{
		eventTable: make(map[Topic][]EventHandlerFunc),
		eventSub:   make(chan eventSubscription),
		eventPub:   make(chan request),

		cqTable: make(map[Topic]CommandQueryHandlerFunc),
		cqSub:   make(chan commandQuerySubscription),
		cqPub:   make(chan request),

		stop: make(chan struct{}),
	}

	go b.run()
	return b
}

type Broker struct {
	eventTable map[Topic][]EventHandlerFunc
	eventSub   chan eventSubscription
	eventPub   chan request

	cqTable map[Topic]CommandQueryHandlerFunc
	cqSub   chan commandQuerySubscription
	cqPub   chan request

	closeMu   sync.RWMutex
	isClose   bool
	stop      chan struct{}
	onceClose sync.Once
}

func (b *Broker) run() {
	for {
		select {
		case sub := <-b.eventSub:
			b.eventTable[sub.topic] = append(b.eventTable[sub.topic], sub.handler)

		case req := <-b.eventPub:
			consumers := b.eventTable[req.payload.Topic()]
			go b.broadcastEvent(req, consumers)

		case sub := <-b.cqSub:
			_, ok := b.cqTable[sub.topic]
			if ok {
				panic(fmt.Sprintf("repeat subscribe topic: %v", sub.topic))
			}
			b.cqTable[sub.topic] = sub.handler

		case <-b.stop:
			return
		}
	}
}

func (b *Broker) broadcastEvent(req request, consumers []EventHandlerFunc) {
	var wg sync.WaitGroup
	for _, c := range consumers {
		consumer := c
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := consumer(req.ctx, req.payload.(Event))

			if req.blocking {
				req.reply <- result{payload: nil, err: err}
			}
		}()
	}

	wg.Wait()
	if req.blocking {
		req.complete(req.payload.Topic())
	}
}

func (b *Broker) Close() {
	b.onceClose.Do(func() {
		b.closeMu.Lock()
		b.isClose = true
		b.closeMu.Unlock()
		close(b.stop)
	})
}

func (b *Broker) SubscribeEvent(topic Topic, handlers ...EventHandlerFunc) {
	b.closeMu.RLock()
	defer b.closeMu.RUnlock()

	if b.isClose {
		return
	}
	for _, handler := range handlers {
		b.eventSub <- eventSubscription{topic: topic, handler: handler}
	}
}

func (b *Broker) BlockingPublishEvent(ctx context.Context, events ...Event) error {
	ready := make(chan Topic, len(events))
	reply := make(chan result)

	b.closeMu.RLock()
	if b.isClose {
		return errors.Wrap(errors.ErrSystem, "message broker closed")
	}

	for i := range events {
		req := request{
			ctx:      ctx,
			payload:  events[i],
			blocking: true,
			ready:    ready,
			reply:    reply,
		}
		b.eventPub <- req
	}
	b.closeMu.RUnlock()

	go waitMessageComplete(ready, reply, len(events))

	errCh := make(chan error, 1)
	go func() {
		var errSet []error
		for r := range reply {
			if r.err != nil {
				errSet = append(errSet, r.err)
			}
		}

		if len(errSet) == 0 {
			errCh <- nil
			return
		}

		err := errors.ErrSystem
		for _, e := range errSet {
			err = errors.WrapMessage(err, e.Error())
		}
		errCh <- err
	}()

	return <-errCh
}

func (b *Broker) NonBlockingPublishEvent(ctx context.Context, events ...Event) {
	b.closeMu.RLock()
	defer b.closeMu.RUnlock()

	if b.isClose {
		return
	}

	for i := range events {
		req := request{
			ctx:      ctx,
			payload:  events[i],
			blocking: false,
			ready:    nil,
			reply:    nil,
		}
		b.eventPub <- req
	}
}

func waitMessageComplete(ready <-chan Topic, reply chan result, count int) {
	waitList := make(map[Topic]bool)
	timeout := time.After(12 * time.Hour)
	for {
		select {
		case topic := <-ready:
			waitList[topic] = true
			if len(waitList) == count {
				close(reply)
				return
			}
		case <-timeout:
			close(reply)
			return
		}
	}
}

func (b *Broker) SubscribeCommand(topic Topic, handler CommandQueryHandlerFunc) {
	b.closeMu.RLock()
	defer b.closeMu.RUnlock()

	if b.isClose {
		return
	}
	b.cqSub <- commandQuerySubscription{topic: topic, handler: handler}
}
