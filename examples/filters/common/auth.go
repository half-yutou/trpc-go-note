package common

import (
	"context"
	"errors"
	"fmt"

	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/filter"
)

// AuthFilter verifies if the request carries a valid token.
func AuthFilter(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (rsp interface{}, err error) {
	msg := codec.Message(ctx)
	md := msg.ServerMetaData()

	token := string(md["authorization"])
	if token != "secret-token-123" {
		fmt.Printf("[AUTH] Invalid Token: %s\n", token)
		return nil, errors.New("unauthorized: invalid token")
	}

	fmt.Println("[AUTH] Token Verified")
	return next(ctx, req)
}

func init() {
	filter.Register("auth", AuthFilter, nil)
}
