// Package bus implements a small in-process publish/subscribe event bus.
package bus

import (
	"context"
	"sync"
)

// Event is a message published on a topic.
type Event struct {
	Topic string
	Data  string
}

// Bus delivers published events to the subscribers of a topic.
//
// Delivery is best effort: Publish never blocks, so an event is dropped for a subscriber whose
// buffer is full.
type Bus struct {
	mu     sync.RWMutex // guards every field below
	subs   map[string]map[int]chan Event
	nextID int
	closed bool
}

// New returns an empty Bus.
func New() *Bus {
	return &Bus{subs: map[string]map[int]chan Event{}}
}

// Subscribe registers a subscriber for topic. The returned channel has the given buffer size (a
// negative size is treated as zero) and is closed by the returned unsubscribe function or by Close.
// The unsubscribe function is safe to call more than once.
func (b *Bus) Subscribe(topic string, buffer int) (<-chan Event, func()) {
	if buffer < 0 {
		buffer = 0
	}
	ch := make(chan Event, buffer)

	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		close(ch)
		return ch, func() {}
	}

	if b.subs[topic] == nil {
		b.subs[topic] = map[int]chan Event{}
	}
	id := b.nextID
	b.nextID++
	b.subs[topic][id] = ch

	unsubscribe := func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		// The subscriber may already have been removed by an earlier call or by Close.
		if _, ok := b.subs[topic][id]; !ok {
			return
		}
		delete(b.subs[topic], id)
		if len(b.subs[topic]) == 0 {
			delete(b.subs, topic)
		}
		close(ch)
	}
	return ch, unsubscribe
}

// Publish delivers data to every subscriber of topic without blocking. Subscribers whose buffer is
// full miss the event. Publishing on a closed Bus is a no-op.
func (b *Bus) Publish(topic, data string) {
	ev := Event{Topic: topic, Data: data}

	// Sends are non-blocking, so holding the read lock is cheap, and it guarantees that no
	// channel is closed (by unsubscribe or Close) while we send on it.
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subs[topic] {
		select {
		case ch <- ev:
		default:
		}
	}
}

// Once returns a channel that receives the first event published on topic and is then closed.
// It is also closed, without an event, when ctx is cancelled or the Bus is closed.
func (b *Bus) Once(ctx context.Context, topic string) <-chan Event {
	out := make(chan Event, 1)
	ch, unsubscribe := b.Subscribe(topic, 1)
	go func() {
		defer close(out)
		defer unsubscribe()
		select {
		case ev, ok := <-ch:
			if ok {
				out <- ev
			}
		case <-ctx.Done():
		}
	}()
	return out
}

// Len returns the number of subscribers of topic.
func (b *Bus) Len(topic string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subs[topic])
}

// Close shuts the bus down and closes all subscriber channels. It is safe to call more than once.
func (b *Bus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for _, subs := range b.subs {
		for _, ch := range subs {
			close(ch)
		}
	}
	b.subs = map[string]map[int]chan Event{}
}
