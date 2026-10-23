package bus

import (
	"context"
	"sync"
	"testing"
	"time"
)

// within runs fn and fails the test if it panics or does not return in time.
func within(t *testing.T, d time.Duration, what string, fn func()) {
	t.Helper()
	done := make(chan any, 1)
	go func() {
		defer func() { done <- recover() }()
		fn()
	}()
	select {
	case p := <-done:
		if p != nil {
			t.Fatalf("%s: panic: %v", what, p)
		}
	case <-time.After(d):
		t.Fatalf("%s: did not finish within %v", what, d)
	}
}

func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(time.Millisecond)
	}
}

func recv(t *testing.T, ch <-chan Event) (Event, bool) {
	t.Helper()
	select {
	case ev, ok := <-ch:
		return ev, ok
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for a value on the channel")
		return Event{}, false
	}
}

func TestPublishReachesOnlyTopicSubscribers(t *testing.T) {
	b := New()
	defer b.Close()
	orders1, _ := b.Subscribe("orders", 4)
	orders2, _ := b.Subscribe("orders", 4)
	emails, _ := b.Subscribe("emails", 4)

	b.Publish("orders", "o1")
	b.Publish("orders", "o2")

	for _, ch := range []<-chan Event{orders1, orders2} {
		for _, want := range []string{"o1", "o2"} {
			if ev, _ := recv(t, ch); ev != (Event{Topic: "orders", Data: want}) {
				t.Fatalf("got %+v; want orders/%s", ev, want)
			}
		}
	}
	select {
	case ev := <-emails:
		t.Fatalf("emails subscriber received %+v", ev)
	default:
	}
}

func TestSlowSubscriberDoesNotBlockOthers(t *testing.T) {
	b := New()
	b.Subscribe("t", 1) // never read
	fast, _ := b.Subscribe("t", 16)

	within(t, 2*time.Second, "Publish with a stuck subscriber", func() {
		for i := 0; i < 5; i++ {
			b.Publish("t", "x")
		}
	})
	within(t, 2*time.Second, "Subscribe after a stuck Publish", func() {
		_, unsub := b.Subscribe("t", 1)
		unsub()
	})
	for i := 0; i < 5; i++ {
		if _, ok := recv(t, fast); !ok {
			t.Fatal("fast subscriber channel closed early")
		}
	}
}

func TestFullBufferDropsEventsInsteadOfBlocking(t *testing.T) {
	b := New()
	defer b.Close()
	ch, _ := b.Subscribe("t", 2)
	within(t, time.Second, "Publish into a full buffer", func() {
		for i := 0; i < 10; i++ {
			b.Publish("t", "x")
		}
	})
	if got := len(ch); got != 2 {
		t.Fatalf("buffered %d events; want 2", got)
	}
}

func TestNegativeBufferIsTreatedAsUnbuffered(t *testing.T) {
	b := New()
	defer b.Close()
	within(t, time.Second, "Subscribe with a negative buffer", func() {
		_, unsub := b.Subscribe("t", -1)
		unsub()
	})
}

func TestUnsubscribeClosesChannelAndStopsDelivery(t *testing.T) {
	b := New()
	defer b.Close()
	ch, unsub := b.Subscribe("t", 4)
	other, _ := b.Subscribe("t", 4)

	unsub()

	if _, ok := recv(t, ch); ok {
		t.Fatal("channel still open after unsubscribe")
	}
	if got := b.Len("t"); got != 1 {
		t.Fatalf("Len = %d; want 1", got)
	}
	within(t, time.Second, "Publish after unsubscribe", func() { b.Publish("t", "x") })
	if ev, _ := recv(t, other); ev.Data != "x" {
		t.Fatalf("remaining subscriber got %+v", ev)
	}
}

func TestUnsubscribeTwiceIsSafe(t *testing.T) {
	b := New()
	defer b.Close()
	_, unsub := b.Subscribe("t", 1)
	within(t, time.Second, "double unsubscribe", func() {
		unsub()
		unsub()
	})
}

func TestUnsubscribeAfterCloseIsSafe(t *testing.T) {
	b := New()
	_, unsub := b.Subscribe("t", 1)
	b.Close()
	within(t, time.Second, "unsubscribe after Close", unsub)
}

func TestConcurrentPublishSubscribeUnsubscribe(t *testing.T) {
	b := New()

	var wg sync.WaitGroup
	run := func(fn func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if p := recover(); p != nil {
					t.Errorf("panic: %v", p)
				}
			}()
			fn()
		}()
	}

	for i := 0; i < 4; i++ {
		run(func() {
			for j := 0; j < 200; j++ {
				b.Publish("t", "x")
			}
		})
		run(func() {
			for j := 0; j < 200; j++ {
				_, unsub := b.Subscribe("t", 4)
				unsub()
			}
		})
	}

	finished := make(chan struct{})
	go func() { wg.Wait(); close(finished) }()
	select {
	case <-finished:
	case <-time.After(10 * time.Second):
		t.Fatal("publishers and subscribers did not finish (deadlock?)")
	}
	if got := b.Len("t"); got != 0 {
		t.Fatalf("Len = %d; want 0", got)
	}
	b.Close()
}

func TestCloseClosesSubscribersAndIsIdempotent(t *testing.T) {
	b := New()
	a, _ := b.Subscribe("a", 1)
	c, _ := b.Subscribe("c", 1)

	within(t, time.Second, "Close twice", func() {
		b.Close()
		b.Close()
	})
	for _, ch := range []<-chan Event{a, c} {
		if _, ok := recv(t, ch); ok {
			t.Fatal("subscriber channel still open after Close")
		}
	}
}

func TestUseAfterClose(t *testing.T) {
	b := New()
	b.Close()

	within(t, time.Second, "Publish after Close", func() { b.Publish("t", "x") })

	var ch <-chan Event
	var unsub func()
	within(t, time.Second, "Subscribe after Close", func() { ch, unsub = b.Subscribe("t", 1) })
	if _, ok := recv(t, ch); ok {
		t.Fatal("channel from a closed bus should be closed")
	}
	within(t, time.Second, "unsubscribe on closed bus", unsub)
	if got := b.Len("t"); got != 0 {
		t.Fatalf("Len = %d; want 0", got)
	}
}

func TestOnceDeliversFirstEventThenCloses(t *testing.T) {
	b := New()
	defer b.Close()
	out := b.Once(context.Background(), "t")

	b.Publish("t", "first")
	b.Publish("t", "second")

	if ev, ok := recv(t, out); !ok || ev.Data != "first" {
		t.Fatalf("got %+v, %v; want the first event", ev, ok)
	}
	if _, ok := recv(t, out); ok {
		t.Fatal("Once channel should be closed after one event")
	}
	waitFor(t, "Once to unsubscribe", func() bool { return b.Len("t") == 0 })
}

func TestOnceStopsWhenContextIsCancelled(t *testing.T) {
	b := New()
	defer b.Close()
	ctx, cancel := context.WithCancel(context.Background())
	out := b.Once(ctx, "t")
	waitFor(t, "Once to subscribe", func() bool { return b.Len("t") == 1 })

	cancel()

	if ev, ok := recv(t, out); ok {
		t.Fatalf("got %+v; want the channel closed without an event", ev)
	}
	waitFor(t, "Once to unsubscribe after cancel", func() bool { return b.Len("t") == 0 })
}

func TestOnceOnClosedBusDeliversNothing(t *testing.T) {
	b := New()
	out := b.Once(context.Background(), "t")
	b.Close()
	if ev, ok := recv(t, out); ok {
		t.Fatalf("got %+v; a closed bus must not produce a zero-value event", ev)
	}
}
