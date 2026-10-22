// Package pool runs batches of jobs on a fixed number of worker goroutines.
package pool

import (
	"context"
	"sort"
	"sync"
)

// Job is a unit of work.
type Job struct {
	ID      int
	Payload string
}

// Result is the outcome of a single Job.
type Result struct {
	JobID  int
	Output string
	Err    error
}

// Handler processes a single job. It should return promptly once ctx is done.
type Handler func(ctx context.Context, j Job) (string, error)

// Run processes jobs using the given number of workers and returns exactly
// one Result per job, ordered by JobID. Jobs that were not started because
// ctx was canceled get a Result whose Err is ctx.Err(). A workers value below
// one is treated as one.
func Run(ctx context.Context, workers int, jobs []Job, h Handler) []Result {
	if workers < 1 {
		workers = 1
	}
	jobCh := make(chan Job)
	// Buffered so a worker never blocks sending a result: every job
	// produces exactly one Result.
	resCh := make(chan Result, len(jobs))

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := range jobCh {
				out, err := h(ctx, j)
				resCh <- Result{JobID: j.ID, Output: out, Err: err}
			}
		}()
	}

feed:
	for i, j := range jobs {
		select {
		case jobCh <- j:
		case <-ctx.Done():
			for _, rest := range jobs[i:] {
				resCh <- Result{JobID: rest.ID, Err: ctx.Err()}
			}
			break feed
		}
	}
	close(jobCh) // the sender closes the channel, exactly once
	wg.Wait()
	close(resCh)

	results := make([]Result, 0, len(jobs))
	for r := range resCh {
		results = append(results, r)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].JobID < results[j].JobID })
	return results
}
