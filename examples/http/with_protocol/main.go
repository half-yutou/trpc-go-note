package main

import (
	"context"
	"fmt"

	pb "trpc-go-note/examples/helloworld/pb"

	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

// GreeterServerImpl 实现 GreeterService 接口
type GreeterServerImpl struct{}

// Hello 具体的业务逻辑
// 这个方法既会被 RPC 调用触发，也会被 HTTP 调用触发
func (s *GreeterServerImpl) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Infof("Received request: %s", req.Msg)

	// 无论是 RPC 还是 HTTP，这里处理逻辑是一样的
	// HTTP 请求的 JSON Body 会被框架自动反序列化为 req *pb.HelloRequest
	rsp := &pb.HelloReply{
		Msg: "Hello " + req.Msg + " (from RPC-Style HTTP)",
	}

	return rsp, nil
}

func main() {
	s := trpc.NewServer()
	impl := &GreeterServerImpl{}

	// 注册到 RPC Service
	pb.RegisterGreeterService(s.Service("trpc.helloworld.GreeterRPC"), impl)

	// 注册到 HTTP Service
	// 注意：这里必须显式获取 GreeterHTTP 这个 service name
	pb.RegisterGreeterService(s.Service("trpc.helloworld.Greeter"), impl)

	log.Infof("Server is serving...")
	if err := s.Serve(); err != nil {
		fmt.Println(err)
	}
}
