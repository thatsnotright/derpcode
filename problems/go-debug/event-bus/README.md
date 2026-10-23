## Your task

You are reviewing a teammate's PR. Pair with them to get it PR-worthy:

- Find and fix the bugs in the code.
- Make the tests pass reliably with `go vet ./... && go test -race -count=1 ./...`.
- Note the best-practice feedback you would give in review (naming, docs, structure, tests).

---

## PR: Add in-process event bus

### What it does

Adds `bus.Bus`, a small topic-based pub/sub used to decouple our order, email and audit components.
Many goroutines publish and subscribe at the same time.

- `Subscribe(topic, buffer)` returns a receive-only channel of `Event`s plus an unsubscribe function.
  Unsubscribing closes the channel.
- `Publish(topic, data)` delivers an event to every current subscriber of the topic. A slow
  subscriber must not hold up publishers or other subscribers.
- `Once(ctx, topic)` returns a channel that receives the first event on a topic and is then closed;
  it is also closed if `ctx` is cancelled first.
- `Len(topic)` reports the number of subscribers (handy for metrics and tests).
- `Close()` shuts the bus down and closes every subscriber channel.

I used an `RWMutex` so publishes can run in parallel.

### How I tested it

Unit tests for subscribe/publish and topic isolation. Everything passes when I run it once, but our
audit component sometimes hangs in staging and I have seen a `panic: send on closed channel` once.
I am not sure why. Can you take a look?
