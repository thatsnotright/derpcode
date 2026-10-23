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
type Bus struct {
	mu     sync.RWMutex
	subs   map[string]map[int]chan Event
	nextID int
}

// New returns an empty Bus.
func New() *Bus {
	return &Bus{subs: map[string]map[int]chan Event{}}
}

// Subscribe registers a subscriber for topic. The returned channel has the given buffer size and
// is closed by the returned unsubscribe function.
func (b *Bus) Subscribe(topic string, buffer int) (<-chan Event, func()) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.subs[topic] == nil {
		b.subs[topic] = map[int]chan Event{}
	}
	id := b.nextID
	b.nextID++
	ch := make(chan Event, buffer)
	b.subs[topic][id] = ch

	unsubscribe := func() {
		b.mu.RLock()
		defer b.mu.RUnlock()
		delete(b.subs[topic], id)
		close(ch)
	}
	return ch, unsubscribe
}

// Publish delivers data to every subscriber of topic.
func (b *Bus) Publish(topic, data string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	ev := Event{Topic: topic, Data: data}
	for _, ch := range b.subs[topic] {
		ch <- ev
	}
}

// Once returns a channel that receives the first event published on topic and is then closed.
// It is also closed when ctx is cancelled.
func (b *Bus) Once(ctx context.Context, topic string) <-chan Event {
	out := make(chan Event, 1)
	ch, unsubscribe := b.Subscribe(topic, 1)
	go func() {
		defer unsubscribe()
		out <- <-ch
		close(out)
	}()
	return out
}

// Len returns the number of subscribers of topic.
func (b *Bus) Len(topic string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subs[topic])
}

// Close shuts the bus down and closes all subscriber channels.
func (b *Bus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, subs := range b.subs {
		for _, ch := range subs {
			close(ch)
		}
	}
}
