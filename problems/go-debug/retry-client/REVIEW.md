# Review notes: retry-client

## Planted issues

1. **`client.go: Do` loops `attempt < c.MaxRetries`.** `MaxRetries` is documented as retries *after* the first attempt, so the loop sends one request too few, and with `MaxRetries == 0` it sends none and returns `nil, nil` (a nil-pointer panic for the caller). Fix: `attempts := MaxRetries + 1` (at least 1). Caught by: `TestAttemptCountIsOnePlusMaxRetries`, `TestZeroRetriesStillSendsOnce`.
2. **`client.go: shouldRetry` retries every status `>= 400`.** 4xx means the request is wrong; retrying cannot help and hammers the server (and 401/422 retries can lock accounts or trip rate limits). Fix: retry only network errors, 429 and 5xx. Caught by: `TestClientErrorsAreNotRetried`.
3. **`client.go: Do` retries any method, including `POST /charges`.** A 503 after the server already processed the charge means a duplicate charge. Fix: only retry idempotent methods (GET/HEAD/PUT/DELETE/OPTIONS/TRACE) or requests carrying an `Idempotency-Key`. Caught by: `TestPostIsNotRetriedWithoutIdempotencyKey`.
4. **`client.go: Do` reuses `req.Body` on every attempt.** The first attempt consumes the reader, so retries send an empty body (or fail with a ContentLength mismatch). Fix: clone the request per attempt and get a fresh body from `req.GetBody`; if there is a body but no `GetBody`, do not retry. Caught by: `TestRequestBodyIsResentOnRetry`.
5. **`client.go: Do` drops the failed `resp` without draining and closing its body.** Every retried response leaks a connection (the transport cannot reuse it, file descriptors pile up under load). Fix: `io.Copy(io.Discard, io.LimitReader(body, N))` then `Close()` before the next attempt. Caught by: `TestDiscardedResponsesAreDrainedAndClosed`.
6. **`client.go: Do` uses `time.Sleep` for the backoff.** It ignores `ctx`: a cancelled request keeps sleeping (up to the full backoff) and then fires another attempt with a dead context. Also it sleeps after the last attempt, delaying the final result for nothing. Fix: `select` on a `time.Timer` and `ctx.Done()`, check `ctx.Err()` at the top of each attempt, no sleep after the last attempt. Caught by: `TestContextCancelStopsBackoff`, `TestAttemptCountIsOnePlusMaxRetries` (sleep count).
7. **`client.go: backoff` is `base * time.Duration(1<<attempt)` with no cap.** Overflows to negative/zero around attempt 30+ with a 1s base (a negative `Sleep` returns immediately: a retry storm), and an uncapped delay is useless anyway. Fix: cap at a `maxDelay`, check `base > maxDelay>>attempt` before shifting. Caught by: `TestBackoffGrowsAndNeverOverflows`.
8. **Review only: robustness and API design.** `New` uses `http.DefaultClient` (no timeout, so a hung server blocks an attempt forever; the reference uses a client with a timeout); the last error is returned without saying which attempt or request failed (wrap with `%w` and context); no jitter, so many clients that fail together retry together (thundering herd); `Retry-After` is ignored for 429/503; exported mutable `Backoff` field and `HTTP` field make the zero value unsafe (`Client{}` panics on a nil `HTTP`).

## Coaching / best-practice talking points

- Retries are only safe for idempotent operations. Ask "what happens to the customer if the first POST actually succeeded?". Idempotency keys are the standard answer.
- Always handle the request body on retry: `GetBody`, `req.Clone`, or buffer small bodies. Mention that `http.NewRequest` sets `GetBody` for `bytes.Reader`, `bytes.Buffer` and `strings.Reader` only.
- The "close the body" rule: `defer resp.Body.Close()` for the response you return, drain-and-close for those you discard. Explain why draining matters (keep-alive reuse) and why you cap the drain.
- Context hygiene: every blocking wait (sleep, channel, lock-wait) should also select on `ctx.Done()`. `time.Sleep` is almost always wrong in library code. Prefer `time.NewTimer` + `Stop` over `time.After` in loops.
- Exponential backoff needs a cap, overflow protection and jitter ("full jitter"). Inject the delay function so tests are fast; here the test for the real function calls it directly instead of sleeping.
- Name the contract precisely: is `MaxRetries` retries or attempts? Document it, and test the boundary (0, 1, N).
- Make failure observable: wrap errors with context (`%w`), consider logging/metrics per attempt, and keep returning the final response for the caller to interpret.
- Tests: use `httptest.NewServer` with a scripted handler and `atomic` counters (handlers run on other goroutines, so plain ints would race); a tracking `RoundTripper` is a neat way to assert on connection hygiene; `t.Cleanup` over `defer` in helpers.
- Possible extensions to discuss: `RoundTripper` middleware instead of a wrapper client, circuit breaker, per-attempt timeouts vs overall deadline, `Retry-After` support.

## How to run the conversation

- Ask the engineer to run the tests and talk through the first failure before touching code. Start with `TestZeroRetriesStillSendsOnce`: it is a quick win and reveals the loop bound.
- For the POST issue, ask the business question ("what does this do to a charge?") before the technical one.
- When the body test fails, ask what `req.Body` is after the first `Do`, and let them find `GetBody` in the docs.
- For the leak, ask how they would prove it (a tracking transport, `lsof`, `netstat`) rather than telling them to add a `Close`.
- Ask "what happens if the caller cancels while we are backing off for 30 seconds?" and let them design the select.
- Keep default-client timeout, jitter and `Retry-After` for the "what else would you flag?" part, and the `RoundTripper` design as a stretch discussion.
