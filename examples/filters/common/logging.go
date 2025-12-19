package common

import (
	"context"
	"fmt"
	"time"

	"trpc.group/trpc-go/trpc-go/filter"
)

// LoggingFilter returns a server filter that logs request and response.
func LoggingFilter(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (rsp interface{}, err error) {
	start := time.Now()
	fmt.Printf("[LOG] Recv Request: %+v\n", req)

	rsp, err = next(ctx, req)

	cost := time.Since(start)
	if err != nil {
		fmt.Printf("[LOG] Handle Error: %v, Cost: %v\n", err, cost)
	} else {
		fmt.Printf("[LOG] Send Response: %+v, Cost: %v\n", rsp, cost)
	}
	return rsp, err
}

// Register filters
func init() {
	filter.Register("logging", LoggingFilter, nil)
}
