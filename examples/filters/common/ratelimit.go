package common

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"trpc.group/trpc-go/trpc-go/filter"
)

var (
	maxRequests = int64(100)
	curRequests = int64(0)
)

// RateLimitFilter limits the concurrency of requests.
func RateLimitFilter(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (rsp interface{}, err error) {
	current := atomic.AddInt64(&curRequests, 1)
	defer atomic.AddInt64(&curRequests, -1)

	if current > maxRequests {
		fmt.Printf("[RATELIMIT] Request Rejected. Current: %d\n", current)
		return nil, errors.New("server overloaded")
	}

	return next(ctx, req)
}

func init() {
	filter.Register("ratelimit", RateLimitFilter, nil)
}
