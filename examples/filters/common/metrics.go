package common

import (
	"context"
	"fmt"
	"time"

	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/filter"
)

// MetricsFilter collects execution time metrics.
func MetricsFilter(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (rsp interface{}, err error) {
	start := time.Now()

	rsp, err = next(ctx, req)

	cost := time.Since(start)
	msg := codec.Message(ctx)
	rpcName := msg.ServerRPCName()

	// In real world, you should report to Prometheus/InfluxDB here.
	fmt.Printf("[METRICS] RPC: %s, Cost: %dms\n", rpcName, cost.Milliseconds())

	return rsp, err
}

func init() {
	filter.Register("metrics", MetricsFilter, nil)
}
