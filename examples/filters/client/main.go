package main

import (
	"context"
	"fmt"
	"log"

	pb "trpc-go-note/examples/helloworld/pb"

	"trpc.group/trpc-go/trpc-go/client"
)

func main() {
	// Create client proxy
	// Note: We inject metadata for AuthFilter
	proxy := pb.NewGreeterClientProxy(
		client.WithTarget("ip://127.0.0.1:8000"),
		client.WithMetaData("authorization", []byte("secret-token-123")),
	)

	// 1. Normal Request
	fmt.Println("--- Test 1: Normal Request ---")
	rsp, err := proxy.Hello(context.Background(), &pb.HelloRequest{Msg: "World"})
	if err != nil {
		log.Println(err)
	} else {
		log.Println(rsp.Msg)
	}

	// 2. Panic Request (Test Recovery)
	fmt.Println("\n--- Test 2: Panic Request ---")
	_, err = proxy.Hello(context.Background(), &pb.HelloRequest{Msg: "panic"})
	if err != nil {
		log.Printf("Received Expected Error: %v\n", err)
	}

	// 3. Auth Fail Request
	fmt.Println("\n--- Test 3: Auth Fail ---")
	badProxy := pb.NewGreeterClientProxy(
		client.WithTarget("ip://127.0.0.1:8000"),
		client.WithMetaData("authorization", []byte("wrong-token")),
	)
	_, err = badProxy.Hello(context.Background(), &pb.HelloRequest{Msg: "Hacker"})
	if err != nil {
		log.Printf("Received Expected Error: %v\n", err)
	}
}
