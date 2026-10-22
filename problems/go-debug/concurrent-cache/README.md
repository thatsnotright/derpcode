## Your task

You are reviewing a teammate's PR. Pair with them to get it PR-worthy:

- Find and fix the bugs in the code.
- Make the tests pass reliably with `go vet ./... && go test -race -count=1 ./...`.
- Note the best-practice feedback you would give in review (naming, docs, structure, tests).

---

## PR: Add TTL cache for user sessions

### What it does

Adds a small in-memory string cache (`cache.Cache`) with a fixed per-entry TTL. It will be shared by
our HTTP handlers, so many goroutines call `Set`/`Get` at the same time.

- `Set(key, value)` stores a value that expires after the TTL.
- `Get(key)` returns the value if present and not expired, and counts hits/misses.
- `Delete(key)` removes a key.
- `Cleanup()` removes expired entries; a background goroutine calls it every TTL.
- `Len()` and `Stats()` are there for metrics.

I used a `sync.RWMutex` so reads can happen in parallel.

### How I tested it

Added unit tests for set/get, expiry, and a test that hammers the cache from 20 goroutines.

I think it works, but the tests are flaky sometimes and once the whole test binary just hung in CI.
I'm not sure why. Can you take a look?
