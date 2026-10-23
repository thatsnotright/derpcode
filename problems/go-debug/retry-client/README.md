## Your task

You are reviewing a teammate's PR. Pair with them to get it PR-worthy:

- Find and fix the bugs in the code.
- Make the tests pass reliably with `go vet ./... && go test -race -count=1 ./...`.
- Note the best-practice feedback you would give in review (naming, docs, structure, tests).

---

## PR: Add retrying HTTP client wrapper

### What it does

Adds `retry.Client`, a thin wrapper around `*http.Client` that retries requests which fail
transiently, so every service does not have to hand-roll its own loop.

- `New(maxRetries, baseDelay)` builds a client; `MaxRetries` is the number of retries *after* the
  first attempt, so `MaxRetries: 3` means at most 4 requests.
- Network errors and 5xx/429 responses are retried with exponential backoff (`baseDelay`, 2x, 4x, ...).
- `Do(ctx, req)` returns the final response (even if it is still a 5xx) or the final error.
- `Client.Backoff` can be set to override the delay so tests do not have to wait.

We will use it for calls to the billing service, including `POST /charges`.

### How I tested it

Tests run against `httptest.NewServer` handlers that fail a few times and then succeed, with a 1ms
backoff. I have not tried cancelling a request, and CI shows some odd test failures I cannot
reproduce. Can you take a look?
