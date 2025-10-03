package utils

import (
	"context"
	"time"
)

// RetryWithContext retries a function with context cancellation support
func RetryWithContext(ctx context.Context, fn func() error, maxRetries int, delay time.Duration) error {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if i < maxRetries {
			backoffDelay := ExponentialBackoff(i) + delay
			select {
			case <-time.After(backoffDelay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return lastErr
}

// RetryFixed retries a function with fixed delay
func RetryFixed(fn func() error, maxRetries int, delay time.Duration) error {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if i < maxRetries {
			time.Sleep(delay)
		}
	}

	return lastErr
}

// RetryWithCondition retries a function with custom retry condition
func RetryWithCondition(fn func() error, shouldRetry func(error) bool, maxRetries int, delay time.Duration) error {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
			if !shouldRetry(err) {
				return err
			}
		}

		if i < maxRetries {
			backoffDelay := ExponentialBackoff(i) + delay
			time.Sleep(backoffDelay)
		}
	}

	return lastErr
}
