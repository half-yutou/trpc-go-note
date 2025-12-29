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
	// 1. 初始化 trpc server
	s := trpc.NewServer()

	// 2. 初始化 Gin Engine
	// gin.Default() 默认带有 Logger 和 Recovery 中间件
	g := gin.Default()

	// 3. 定义路由和处理函数
	// Gin 提供了更强大的路由功能，比如分组、中间件等
	g.GET("/hello", func(c *gin.Context) {
		// 绑定参数
		var req HelloRequest

		// ShouldBind 会自动根据 Content-Type 选择绑定方式 (Query 或 JSON)
		if err := c.ShouldBind(&req); err != nil {
			// 如果绑定失败（例如缺少 required 字段），返回 400 错误
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  err.Error(),
			})
			return
		}

		// 业务逻辑
		log.Infof("Handling request for user: %s, age: %d", req.Name, req.Age)

		// 统一返回 JSON
		// gin.H 是 map[string]interface{} 的快捷方式
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
		// 对于 POST，ShouldBind 默认解析 JSON Body
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
