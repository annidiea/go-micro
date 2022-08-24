package client

import "context"

// note that returning either false or a non-nil error will result in the call not being retried
type RetryFunc func(ctx context.Context, req Request, retryCount int, err error) (bool, error)

// RetryAlways always retry on error
func RetryAlways(ctx context.Context, req Request, retryCount int, err error) (bool, error) {
	return true, nil
}
