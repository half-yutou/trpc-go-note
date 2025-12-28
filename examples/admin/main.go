package main

import (
	"encoding/json"
	"net/http"
	"time"

	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/admin"
	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	// 1. 创建 Server (会自动启动 admin)
	// 默认 Admin 监听在 :9028
	s := trpc.NewServer()

	// 2. 注册自定义 Admin 命令
	// 访问: http://localhost:9028/cmds/my_status
	admin.HandleFunc("/cmds/my_status", func(w http.ResponseWriter, r *http.Request) {
		status := map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"app_state": "running",
			"uptime":    "1h 20m",
		}
		_ = json.NewEncoder(w).Encode(status)
	})

	log.Info("Admin server is running on :9028")
	log.Info("Try: curl http://localhost:9028/cmds")
	log.Info("Try: curl -X PUT -d 'value=debug' http://localhost:9028/cmds/loglevel?logger=default")

	// 3. 模拟业务运行
	// s.Serve() 会阻塞并处理信号 (Ctrl+C)
	// 如果没有任何 Service 注册，它会 panic，但 trpc.NewServer() 默认会自动加载 admin service
	if err := s.Serve(); err != nil {
		log.Error(err)
	}
}
