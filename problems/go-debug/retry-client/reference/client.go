// Package retry provides an HTTP client wrapper that retries transient failures.
package retry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// maxDelay caps the exponential backoff.
	maxDelay = 30 * time.Second
	// maxDrainBytes bounds how much of a discarded response body we read to reuse the connection.
	maxDrainBytes = 4 << 10
	// defaultTimeout is the per-attempt timeout of the default HTTP client.
	defaultTimeout = 10 * time.Second
)

// Client wraps an http.Client and retries failed requests with exponential backoff.
type Client struct {
	// HTTP is the underlying client.
	HTTP *http.Client
	// MaxRetries is the number of retries after the first attempt.
	MaxRetries int
	// BaseDelay is the delay before the first retry; it doubles on each retry.
	BaseDelay time.Duration
	// Backoff, if set, overrides the delay computed from BaseDelay.
	Backoff func(attempt int) time.Duration
}

// New returns a Client that retries up to maxRetries times.
func New(maxRetries int, baseDelay time.Duration) *Client {
	return &Client{
		HTTP:       &http.Client{Timeout: defaultTimeout},
		MaxRetries: maxRetries,
		BaseDelay:  baseDelay,
	}
}

// Do sends req, retrying on network errors and retryable status codes.
//
// Only idempotent requests (or requests carrying an Idempotency-Key header) whose body can be
// replayed are retried. When retries are exhausted, the last response is returned with a nil error.
// The caller must close the returned response body.
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	attempts := 1
	if c.MaxRetries > 0 && canRetry(req) {
		attempts += c.MaxRetries
	}

	for attempt := 0; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		r, err := cloneRequest(ctx, req)
		if err != nil {
			return nil, err
		}

		resp, err := c.HTTP.Do(r)
		if attempt == attempts-1 || !shouldRetry(resp, err) {
			if err != nil {
				return nil, fmt.Errorf("%s %s: attempt %d of %d: %w", req.Method, req.URL, attempt+1, attempts, err)
			}
			return resp, nil
		}
		if resp != nil {
			drainAndClose(resp.Body)
		}
		if err := sleep(ctx, c.delay(attempt)); err != nil {
			return nil, err
		}
	}
}

func (c *Client) delay(attempt int) time.Duration {
	if c.Backoff != nil {
		return c.Backoff(attempt)
	}
	return backoff(c.BaseDelay, attempt)
}

// backoff returns base * 2^attempt, capped at maxDelay without overflowing.
func backoff(base time.Duration, attempt int) time.Duration {
	if base <= 0 {
		return 0
	}
	if attempt < 0 {
		attempt = 0
	}
	if attempt >= 62 || base > maxDelay>>attempt {
		return maxDelay
	}
	return base << attempt
}

// sleep waits for d or until ctx is done, whichever comes first.
func sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// cloneRequest returns a copy of req bound to ctx with a fresh copy of the body.
func cloneRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	r := req.Clone(ctx)
	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return nil, fmt.Errorf("rewind request body: %w", err)
		}
		r.Body = body
	}
	return r, nil
}

// canRetry reports whether req is safe to send more than once.
func canRetry(req *http.Request) bool {
	if req.Body != nil && req.Body != http.NoBody && req.GetBody == nil {
		return false
	}
	switch req.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodPut, http.MethodDelete, http.MethodTrace:
		return true
	}
	return req.Header.Get("Idempotency-Key") != ""
}

func shouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	return resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
}

func drainAndClose(body io.ReadCloser) {
	io.Copy(io.Discard, io.LimitReader(body, maxDrainBytes))
	body.Close()
}
