package main

import (
	"context"

	"trpc-go-note/examples/helloworld/pb"

	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	c := pb.NewGreeterClientProxy(client.WithTarget("ip://127.0.0.1:8000"))
	rsp, err := c.Hello(context.Background(), &pb.HelloRequest{Msg: "world"})
	if err != nil {
		log.Error(err)
	}
	log.Info(rsp.Msg)
}
