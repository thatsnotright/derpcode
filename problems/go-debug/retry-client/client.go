// Package retry provides an HTTP client wrapper that retries transient failures.
package retry

import (
	"context"
	"net/http"
	"time"
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
		HTTP:       http.DefaultClient,
		MaxRetries: maxRetries,
		BaseDelay:  baseDelay,
	}
}

// Do sends req, retrying on network errors and retryable status codes.
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	for attempt := 0; attempt < c.MaxRetries; attempt++ {
		resp, err = c.HTTP.Do(req.WithContext(ctx))
		if !shouldRetry(resp, err) {
			return resp, err
		}
		time.Sleep(c.delay(attempt))
	}
	return resp, err
}

func (c *Client) delay(attempt int) time.Duration {
	if c.Backoff != nil {
		return c.Backoff(attempt)
	}
	return backoff(c.BaseDelay, attempt)
}

func backoff(base time.Duration, attempt int) time.Duration {
	return base * time.Duration(1<<attempt)
}

func shouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	return resp.StatusCode >= 400
}
