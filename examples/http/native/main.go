package main

import (
	"fmt"
	"net/http"

	"trpc.group/trpc-go/trpc-go"
	trpc_http "trpc.group/trpc-go/trpc-go/http"
	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	// 1. 初始化 trpc server
	// 框架会自动读取并加载同目录下的 trpc_go.yaml
	s := trpc.NewServer()

	// 2. 定义原生 HTTP 处理逻辑
	// 这里使用标准库的 ServeMux，你也可以使用 Gin, Echo 等框架
	mux := http.NewServeMux()

	// 定义一个简单的路由
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Infof("received request: %s %s", r.Method, r.URL.Path)

		// 获取查询参数
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "Native HTTP"
		}

		// 返回 JSON 响应
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Hello, %s!"}`, name)
	})

	// 3. 注册 HTTP 服务
	// 获取配置文件中定义的 service
	serviceName := "trpc.demo.http.MyService"
	service := s.Service(serviceName)
	if service == nil {
		log.Fatalf("failed to find service: %s, please check trpc_go.yaml", serviceName)
	}

	// 核心步骤：将 mux 注册到 trpc service
	// 这样所有的 HTTP 请求都会先经过 trpc 的拦截器链，再交给 mux 处理
	trpc_http.RegisterNoProtocolServiceMux(service, mux)

	// 4. 启动服务
	log.Infof("Server is serving at port 8080...")
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
