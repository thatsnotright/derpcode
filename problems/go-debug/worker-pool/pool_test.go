package pool

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
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

func sorted(rs []Result) []Result {
	sort.Slice(rs, func(i, j int) bool { return rs[i].JobID < rs[j].JobID })
	return rs
}

func TestRunProcessesAllJobs(t *testing.T) {
	upper := func(_ context.Context, j Job) (string, error) {
		return strings.ToUpper(j.Payload), nil
	}
	got := sorted(Run(context.Background(), 3, makeJobs(10), upper))
	if len(got) != 10 {
		t.Fatalf("got %d results; want 10", len(got))
	}
	for i, r := range got {
		want := fmt.Sprintf("JOB-%d", i)
		if r.JobID != i || r.Output != want || r.Err != nil {
			t.Errorf("result %d = %+v; want output %q", i, r, want)
		}
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
	got := sorted(Run(context.Background(), 2, makeJobs(6), h))
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
