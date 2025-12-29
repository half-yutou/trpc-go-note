package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trpc.group/trpc-go/trpc-go"
	trpc_http "trpc.group/trpc-go/trpc-go/http"
	"trpc.group/trpc-go/trpc-go/log"
)

// HelloRequest 定义请求参数结构体
// Gin 使用 `form` tag 绑定 Query 参数，`json` tag 绑定 JSON Body
type HelloRequest struct {
	Name string `form:"name" json:"name" binding:"required"` // binding:"required" 增加非空校验
	Age  int    `form:"age"  json:"age"`
}

func main() {
	// 0. 注册拦截器
	// 必须在 trpc.NewServer() 之前注册，否则配置加载时会找不到 filter
	InitFilters()

	// 1. 初始化 trpc server
	s := trpc.NewServer()

	// 2. 初始化 Gin Engine
	g := gin.Default()

	// 3. 定义路由和处理函数
	g.GET("/hello", func(c *gin.Context) {
		var req HelloRequest
		// 自动绑定 Query 或 JSON
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
			return
		}

		log.Infof("Handling request for user: %s, age: %d", req.Name, req.Age)

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"greeting": "Hello, " + req.Name + " from Gin!",
				"is_adult": req.Age >= 18,
			},
		})
	})

	// 演示 POST 请求
	g.POST("/user", func(c *gin.Context) {
		var req HelloRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "created", "user": req})
	})

	// 4. 获取 trpc Service 并注册 Gin
	serviceName := "trpc.demo.http.MyService"
	service := s.Service(serviceName)
	if service == nil {
		log.Fatalf("failed to find service: %s", serviceName)
	}

	// 核心魔法：将 Gin Engine 注册为 trpc 的 Handler
	trpc_http.RegisterNoProtocolServiceMux(service, g)

	// 5. 启动
	log.Infof("Server is serving at port 8080...")
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
