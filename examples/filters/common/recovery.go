package common

import (
	"context"
	"fmt"
	"runtime/debug"

	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/filter"
)

// RecoveryFilter catches panic and converts it to error.
func RecoveryFilter(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (rsp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errs.New(errs.RetServerSystemErr, fmt.Sprintf("panic: %v", r))
			fmt.Printf("[RECOVERY] Captured Panic: %v\nStack: %s\n", r, debug.Stack())
		}
	}()

	return next(ctx, req)
}

func init() {
	filter.Register("recovery", RecoveryFilter, nil)
}
