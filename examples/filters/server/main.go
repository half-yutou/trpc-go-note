package main

import (
	"context"
	"fmt"

	_ "trpc-go-note/examples/filters/common" // Import filters
	pb "trpc-go-note/examples/helloworld/pb"

	trpc "trpc.group/trpc-go/trpc-go"
)

type greeterImpl struct{}

func (s *greeterImpl) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	// Simulate panic for recovery test
	if req.Msg == "panic" {
		panic("something went wrong")
	}
	return &pb.HelloReply{Msg: "Hello " + req.Msg}, nil
}

func main() {
	s := trpc.NewServer()

	pb.RegisterGreeterService(s, &greeterImpl{})
	if err := s.Serve(); err != nil {
		fmt.Println(err)
	}
}
