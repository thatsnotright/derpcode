## Your task

You are reviewing a teammate's PR. Pair with them to get it PR-worthy:

- Find and fix the bugs in the code.
- Make the tests pass reliably with `go vet ./... && go test -race -count=1 ./...`.
- Note the best-practice feedback you would give in review (naming, docs, structure, tests).

---

## PR: Add bounded worker pool for image-processing jobs

### What it does

Adds `pool.Run`, which processes a slice of jobs with a fixed number of worker goroutines and returns
one `Result` per job. Each worker pulls jobs from a channel, calls the supplied `Handler`, and sends
the result back. Callers pass a `context.Context` so they can cancel a long batch.

I wanted something simpler than a full job queue, so it is a single function: start workers, feed the
jobs, wait for the workers, collect the results.

### How I tested it

Unit tests for the normal case, handler errors, an empty job list, and cancellation.

The tests just sit there for a long time when I run them and I have to Ctrl-C. I assumed my laptop
was slow. Can you take a look?
