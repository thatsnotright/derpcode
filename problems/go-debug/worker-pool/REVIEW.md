# Review notes: worker-pool

## Planted issues

1. **`pool.go: Run` - unbuffered `resCh` with the collector running only after `wg.Wait()`.** Workers block on `resCh <- ...` because nobody is receiving, so they never finish; the producer then blocks on `jobCh <- j` once all workers are stuck; `wg.Wait()` is never reached. Classic deadlock for any batch larger than the number of workers. This is the "tests just sit there" symptom. Fix: buffer `resCh` to `len(jobs)` (each job yields exactly one result), or collect concurrently in another goroutine. Caught by: `TestRunProcessesAllJobs` (hangs until the 20s timeout; the stack dump shows workers on `chan send`).
2. **`wg.Add(1)` is called inside the goroutine.** `wg.Wait()` can run before any `Add` has happened and return early (or panic with "WaitGroup misuse"). `Add` must happen before `go`. Fix: `wg.Add(workers)` before the loop (or `Add(1)` right before each `go`). Recent Go versions' `go vet` flags this ("WaitGroup.Add called from inside new goroutine"), so the very first thing the candidate sees may be the vet error; older toolchains need a review catch. Tests do not reliably catch it.
3. **Wrong side closes the channel: every worker does `defer close(resCh)`.** With more than one worker the second close panics (`close of closed channel`), and a worker closing while others still send causes `send on closed channel`. Rule: only the sender (and only one sender, or a coordinator after `wg.Wait()`) closes a channel; receivers never close. Fix: close `resCh` once, after `wg.Wait()`. Caught by: `TestRunProcessesAllJobs` once issue 1 is fixed (panic), `TestRunNoJobs`.
4. **The `ctx` parameter is ignored and the handler is called with `context.Background()`.** Cancellation can never reach the handler, and the feeder keeps pushing jobs after cancel. Fix: pass `ctx` to the handler and `select` on `ctx.Done()` when feeding; report unstarted jobs with `ctx.Err()`. Caught by: `TestRunStopsOnCancel` (takes ~700ms instead of returning in ~30ms).
5. **Goroutine leak on early exit / stuck sends.** Workers blocked sending on `resCh`, and a feeder blocked on `jobCh` (if one were moved to its own goroutine or `Run` returned early) never exit. Anything that exits early must guarantee that each started goroutine has a way out (buffered channel, `select` with `ctx.Done()`, close of the job channel). Review discussion; reference test `TestRunDoesNotLeakGoroutines` checks it.
6. **Review only: no guard for `workers < 1`** (with 0 workers the feeder blocks forever), **no panic recovery** in workers (one panicking handler kills the process), and **result order is nondeterministic** (callers must sort; the reference sorts by `JobID` and documents it). Also `var results []Result` returns `nil` for empty input versus an empty slice; pre-size with `make([]Result, 0, len(jobs))`.

## Coaching / best-practice talking points

- Ownership rules for channels: who sends, who receives, who closes. Write it in a comment.
- Always ask "what unblocks this goroutine?" for every channel operation. Use `select` with `ctx.Done()` on blocking sends/receives.
- `go vet`/`-race` do not find deadlocks or WaitGroup misuse; read the goroutine dump printed on timeout (`-timeout 20s`).
- Prefer returning an `error`/documenting behaviour for invalid args vs silently clamping; discuss the trade-off.
- Consider `golang.org/x/sync/errgroup` (`SetLimit`) in real code; stdlib-only here.
- Test hygiene: always run with `-timeout` and `-race`; avoid `time.Sleep`-based synchronization where possible; assert on behaviour (peak concurrency with atomics), not timing alone; check goroutine leaks.
- `ctx` as the first parameter, never stored in structs; do not pass `context.Background()` when you have a ctx.
- gofmt, doc comments that state ordering/cancellation guarantees.

## How to run the conversation

- Let them run the tests and wait. Ask what they would do when a test hangs (answer: read the stack dump; `-timeout`).
- Ask them to trace by hand: 10 jobs, 3 workers; what is each goroutine doing after 3 jobs have been sent?
- For the WaitGroup and close issues, ask "what does `Wait` need to be true when it is called?" and "who is allowed to close this channel and why?"
- After deadlock is fixed, tests will panic (closed channel) and then the cancel test fails; this staging is intentional so they experience one fix at a time.
