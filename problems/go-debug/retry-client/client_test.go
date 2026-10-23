package retry

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// flaky returns a server that answers with the given statuses in order, then 200 forever.
func flaky(t *testing.T, statuses ...int) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := int(hits.Add(1))
		if n <= len(statuses) {
			w.WriteHeader(statuses[n-1])
			io.WriteString(w, "busy")
			return
		}
		io.WriteString(w, "ok")
	}))
	t.Cleanup(srv.Close)
	return srv, &hits
}

func newClient(retries int) *Client {
	c := New(retries, time.Millisecond)
	c.Backoff = func(int) time.Duration { return time.Millisecond }
	return c
}

func get(t *testing.T, c *Client, url string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Do(context.Background(), req)
	if err != nil || resp == nil {
		t.Fatalf("Do = %v, %v", resp, err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func TestRetriesUntilSuccess(t *testing.T) {
	srv, hits := flaky(t, 503, 503)
	resp := get(t, newClient(3), srv.URL)
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d; want 200", resp.StatusCode)
	}
	if got := hits.Load(); got != 3 {
		t.Fatalf("server saw %d requests; want 3", got)
	}
}

func TestAttemptCountIsOnePlusMaxRetries(t *testing.T) {
	srv, hits := flaky(t, 503, 503, 503, 503, 503, 503)
	var sleeps atomic.Int32
	c := newClient(2)
	c.Backoff = func(int) time.Duration { sleeps.Add(1); return time.Millisecond }

	resp := get(t, c, srv.URL)
	if resp.StatusCode != 503 {
		t.Fatalf("status = %d; want the final 503 to be returned", resp.StatusCode)
	}
	if got := hits.Load(); got != 3 {
		t.Fatalf("server saw %d requests; want 3 (1 + 2 retries)", got)
	}
	if got := sleeps.Load(); got != 2 {
		t.Fatalf("slept %d times; want 2 (no sleep after the last attempt)", got)
	}
}

func TestZeroRetriesStillSendsOnce(t *testing.T) {
	srv, hits := flaky(t)
	resp := get(t, newClient(0), srv.URL)
	if resp.StatusCode != 200 || hits.Load() != 1 {
		t.Fatalf("status = %d, hits = %d; want 200 and 1", resp.StatusCode, hits.Load())
	}
}

func TestClientErrorsAreNotRetried(t *testing.T) {
	for _, status := range []int{400, 401, 404, 422} {
		srv, hits := flaky(t, status, status, status)
		resp := get(t, newClient(3), srv.URL)
		if resp.StatusCode != status || hits.Load() != 1 {
			t.Errorf("status %d: got status %d after %d requests; want 1 request", status, resp.StatusCode, hits.Load())
		}
	}
}

func TestTooManyRequestsIsRetried(t *testing.T) {
	srv, hits := flaky(t, 429)
	resp := get(t, newClient(2), srv.URL)
	if resp.StatusCode != 200 || hits.Load() != 2 {
		t.Fatalf("status = %d, hits = %d; want 200 and 2", resp.StatusCode, hits.Load())
	}
}

func TestPostIsNotRetriedWithoutIdempotencyKey(t *testing.T) {
	srv, hits := flaky(t, 503, 503)
	req, _ := http.NewRequest(http.MethodPost, srv.URL, strings.NewReader(`{"amount":5}`))
	resp, err := newClient(3).Do(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if got := hits.Load(); got != 1 {
		t.Fatalf("POST was sent %d times; a charge must not be duplicated", got)
	}
}

func TestPostWithIdempotencyKeyIsRetried(t *testing.T) {
	srv, hits := flaky(t, 503)
	req, _ := http.NewRequest(http.MethodPost, srv.URL, strings.NewReader(`{"amount":5}`))
	req.Header.Set("Idempotency-Key", "abc-123")
	resp, err := newClient(3).Do(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 || hits.Load() != 2 {
		t.Fatalf("status = %d, hits = %d; want 200 and 2", resp.StatusCode, hits.Load())
	}
}

func TestRequestBodyIsResentOnRetry(t *testing.T) {
	var mu sync.Mutex
	var bodies []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		bodies = append(bodies, string(b))
		n := len(bodies)
		mu.Unlock()
		if n == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodPut, srv.URL, strings.NewReader("payload"))
	resp, err := newClient(2).Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer resp.Body.Close()

	mu.Lock()
	defer mu.Unlock()
	if len(bodies) != 2 || bodies[0] != "payload" || bodies[1] != "payload" {
		t.Fatalf("server saw bodies %q; want the full payload twice", bodies)
	}
}

type trackedBody struct {
	io.ReadCloser
	mu     sync.Mutex
	eof    bool
	closed bool
}

func (b *trackedBody) Read(p []byte) (int, error) {
	n, err := b.ReadCloser.Read(p)
	if err == io.EOF {
		b.mu.Lock()
		b.eof = true
		b.mu.Unlock()
	}
	return n, err
}

func (b *trackedBody) Close() error {
	b.mu.Lock()
	b.closed = true
	b.mu.Unlock()
	return b.ReadCloser.Close()
}

type trackingTransport struct {
	mu     sync.Mutex
	bodies []*trackedBody
}

func (tt *trackingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	tb := &trackedBody{ReadCloser: resp.Body}
	resp.Body = tb
	tt.mu.Lock()
	tt.bodies = append(tt.bodies, tb)
	tt.mu.Unlock()
	return resp, nil
}

func TestDiscardedResponsesAreDrainedAndClosed(t *testing.T) {
	srv, _ := flaky(t, 503, 502)
	tt := &trackingTransport{}
	c := newClient(3)
	c.HTTP = &http.Client{Transport: tt}

	resp := get(t, c, srv.URL)
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d; want 200", resp.StatusCode)
	}

	tt.mu.Lock()
	defer tt.mu.Unlock()
	if len(tt.bodies) != 3 {
		t.Fatalf("saw %d responses; want 3", len(tt.bodies))
	}
	for i, b := range tt.bodies[:2] {
		b.mu.Lock()
		closed, eof := b.closed, b.eof
		b.mu.Unlock()
		if !closed || !eof {
			t.Errorf("response %d: closed=%v drained=%v; want both (leaks the connection otherwise)", i, closed, eof)
		}
	}
}

func TestContextCancelStopsBackoff(t *testing.T) {
	srv, _ := flaky(t, 503, 503, 503)
	c := newClient(3)
	c.Backoff = func(int) time.Duration { return 10 * time.Second }

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	time.AfterFunc(50*time.Millisecond, cancel)

	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	done := make(chan error, 1)
	go func() {
		resp, err := c.Do(ctx, req)
		if resp != nil {
			resp.Body.Close()
		}
		done <- err
	}()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %v; want context.Canceled", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Do kept sleeping after the context was cancelled")
	}
}

func TestBackoffGrowsAndNeverOverflows(t *testing.T) {
	prev := time.Duration(0)
	for attempt := 0; attempt < 100; attempt++ {
		d := backoff(500*time.Millisecond, attempt)
		if d <= 0 || d > time.Minute {
			t.Fatalf("backoff(attempt=%d) = %v; want 0 < d <= 1m", attempt, d)
		}
		if d < prev {
			t.Fatalf("backoff(attempt=%d) = %v is smaller than the previous %v", attempt, d, prev)
		}
		prev = d
	}
	if got := backoff(100*time.Millisecond, 2); got != 400*time.Millisecond {
		t.Fatalf("backoff(100ms, 2) = %v; want 400ms", got)
	}
}
