// Package pool runs batches of jobs on a fixed number of worker goroutines.
package pool

import (
	"context"
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

// Handler processes a single job.
type Handler func(ctx context.Context, j Job) (string, error)

// Run processes jobs using the given number of workers and returns one
// Result per job.
func Run(ctx context.Context, workers int, jobs []Job, h Handler) []Result {
	jobCh := make(chan Job)
	resCh := make(chan Result)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			defer close(resCh)
			for j := range jobCh {
				out, err := h(context.Background(), j)
				resCh <- Result{JobID: j.ID, Output: out, Err: err}
			}
		}()
	}

	for _, j := range jobs {
		jobCh <- j
	}
	close(jobCh)
	wg.Wait()

	var results []Result
	for r := range resCh {
		results = append(results, r)
	}
	return results
}
