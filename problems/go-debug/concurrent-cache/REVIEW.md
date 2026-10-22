# Review notes: concurrent-cache

## Planted issues

1. **`cache.go: Set` uses `RLock` while writing the map.** Multiple goroutines can hold the read lock at once, so concurrent map writes happen (`fatal error: concurrent map writes` / race report). Fix: `Lock` + `defer Unlock`. Caught by: `TestConcurrentAccess` (and `-race`).
2. **`cache.go: Get` returns on the miss path without unlocking.** The next `Set`/`Get` deadlocks; this is the "CI hung once" symptom. Fix: `defer c.mu.Unlock()` immediately after `Lock`. Caught by: `TestMissesAreCounted` (second `Get` hangs until the timeout) once issue 1 is fixed.
3. **`cache.go: expired` uses `time.Since(e.expiresAt) > c.ttl`.** `expiresAt` is already `now+ttl`, so entries effectively live 2x TTL. Fix: `!now.Before(e.expiresAt)` (or `time.Now().After(...)`); decide on `>=` vs `>` boundary deliberately and pass `now` in so Cleanup uses one clock reading. Caught by: `TestExpiry`.
4. **`cache.go: Cleanup` deletes from the map while holding only `RLock`.** Deleting during `range` is legal in Go, but it is a write, so it needs the exclusive lock. Races with `Get`/`Set` and with the janitor goroutine. Fix: `Lock`. Caught by: `-race` (and sometimes concurrent map writes fatal).
5. **`cache.go: Len` and `Stats` read shared state without a lock.** Data race with writers (the janitor and other goroutines). Fix: `RLock` + `defer RUnlock`. Caught by: `-race` in `TestCleanupRemovesExpired`; also an easy review catch ("every field the mutex guards must be accessed under it").
6. **Review only: `New` starts a goroutine with `time.Tick` that can never be stopped.** Goroutine and ticker leak for every cache created (every test leaks one), no way to shut down, and `time.Tick(0)` panics if `ttl <= 0`. Fix: `Close()` with a stop channel + `sync.Once`, `time.NewTicker` with `defer Stop()`, validate or default the interval.

## Coaching / best-practice talking points

- Lock discipline: `mu.Lock(); defer mu.Unlock()` as a pair on adjacent lines. Comment which fields the mutex guards (or group them under it).
- `RWMutex` is only worth it for read-heavy workloads with a real read path; here `Get` mutates counters and deletes, so it needs the write lock. Start with `sync.Mutex`.
- Run `go test -race` locally and in CI always; "flaky" tests are often real races.
- Time handling: inject a clock (`now func() time.Time`) so expiry tests do not need `time.Sleep`. The reference test still sleeps to stay close to the original; mention the trade-off.
- Poll with a deadline instead of a fixed sleep when waiting on a background goroutine.
- Lifecycle of goroutines: anything started in a constructor needs a stop path (`Close`), and tests should `t.Cleanup(c.Close)`.
- Doc comments on exported identifiers; `gofmt`/`go vet` clean; table-driven tests with `t.Run`; `t.Helper()` in helpers.
- Possible extensions to discuss: per-entry TTL, max size/eviction, `sync.Map`, sharding for contention, generics (`Cache[K, V]`).

## How to run the conversation

- Start by asking the engineer to run the tests with `-race` and tell you what they see before you point at anything.
- Ask "what does RLock guarantee? who else can be in this section at the same time?" rather than naming the bug.
- Let them drive the fixes; ask them to predict what a test will do before running it.
- When the hang appears, ask how they would find the stuck goroutine (the stack dump printed on timeout).
- Leave the goroutine-leak and clock-injection topics for last or as discussion if time is short.
