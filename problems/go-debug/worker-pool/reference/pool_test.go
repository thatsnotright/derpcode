package pool

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func makeJobs(n int) []Job {
	jobs := make([]Job, n)
	for i := range jobs {
		jobs[i] = Job{ID: i, Payload: fmt.Sprintf("job-%d", i)}
	}
	return jobs
}

func TestRunProcessesAllJobs(t *testing.T) {
	upper := func(_ context.Context, j Job) (string, error) {
		return strings.ToUpper(j.Payload), nil
	}
	for _, workers := range []int{0, 1, 3, 50} {
		t.Run(fmt.Sprintf("workers=%d", workers), func(t *testing.T) {
			got := Run(context.Background(), workers, makeJobs(10), upper)
			if len(got) != 10 {
				t.Fatalf("got %d results; want 10", len(got))
			}
			for i, r := range got {
				want := fmt.Sprintf("JOB-%d", i)
				if r.JobID != i || r.Output != want || r.Err != nil {
					t.Errorf("result %d = %+v; want output %q", i, r, want)
				}
			}
		})
	}
}

func TestRunReportsHandlerErrors(t *testing.T) {
	boom := errors.New("boom")
	h := func(_ context.Context, j Job) (string, error) {
		if j.ID%2 == 1 {
			return "", boom
		}
		return "ok", nil
	}
	got := Run(context.Background(), 2, makeJobs(6), h)
	if len(got) != 6 {
		t.Fatalf("got %d results; want 6", len(got))
	}
	for _, r := range got {
		if wantErr := r.JobID%2 == 1; wantErr != errors.Is(r.Err, boom) {
			t.Errorf("job %d: err = %v", r.JobID, r.Err)
		}
	}
}

func TestRunNoJobs(t *testing.T) {
	h := func(context.Context, Job) (string, error) { return "", nil }
	if got := Run(context.Background(), 3, nil, h); len(got) != 0 {
		t.Fatalf("got %d results; want 0", len(got))
	}
}

func TestRunStopsOnCancel(t *testing.T) {
	slow := func(ctx context.Context, j Job) (string, error) {
		select {
		case <-time.After(100 * time.Millisecond):
			return "done", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	start := time.Now()
	got := Run(ctx, 2, makeJobs(20), slow)
	if elapsed := time.Since(start); elapsed > 400*time.Millisecond {
		t.Fatalf("Run took %v after cancel; want it to return promptly", elapsed)
	}
	if len(got) != 20 {
		t.Fatalf("got %d results; want one per job (20)", len(got))
	}
	canceled := 0
	for _, r := range got {
		if errors.Is(r.Err, context.DeadlineExceeded) {
			canceled++
		}
	}
	if canceled == 0 {
		t.Fatal("expected some results to report the context error")
	}
}

func TestRunRespectsWorkerLimit(t *testing.T) {
	var running, peak atomic.Int32
	h := func(_ context.Context, _ Job) (string, error) {
		n := running.Add(1)
		for {
			p := peak.Load()
			if n <= p || peak.CompareAndSwap(p, n) {
				break
			}
		}
		time.Sleep(5 * time.Millisecond)
		running.Add(-1)
		return "", nil
	}
	Run(context.Background(), 3, makeJobs(30), h)
	if p := peak.Load(); p > 3 {
		t.Fatalf("peak concurrency = %d; want <= 3", p)
	}
}

func TestRunDoesNotLeakGoroutines(t *testing.T) {
	before := runtime.NumGoroutine()
	h := func(context.Context, Job) (string, error) { return "", nil }
	for i := 0; i < 20; i++ {
		Run(context.Background(), 4, makeJobs(10), h)
	}
	deadline := time.Now().Add(2 * time.Second)
	for runtime.NumGoroutine() > before && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if after := runtime.NumGoroutine(); after > before {
		t.Fatalf("goroutines: before=%d after=%d", before, after)
	}
}
