# 泛HTTP服务与流式服务

- 泛HTTP服务：`http`模块
- 流式服务：`stream`模块

# HTTP 服务分类

tRPC-Go 框架支持搭建与 HTTP 相关的三种服务，它们的核心区别在于**对 Protobuf 的依赖程度**以及**请求处理方式**。

| 服务类型 | 核心特点 | 适用场景 | 配置文件 (protocol) | 桩代码/IDL |
| :--- | :--- | :--- | :--- | :--- |
| **1. 泛 HTTP 标准服务**<br>(No Protocol) | 类似原生 `net/http`，直接操作 `http.Request/ResponseWriter`。<br>框架**不自动**读取/解析 Body。 | 文件上传下载、Webhook 回调、迁移现有 Web 项目 (Gin/Echo)。 | `http_no_protocol` | 不需要 |
| **2. 泛 HTTP RPC 服务**<br>(Standard RPC) | 基于 RPC 逻辑，HTTP 请求被自动映射为 RPC 请求。<br>框架**自动**读取 Body 并反序列化为 PB 结构体。 | 微服务内部通信、简单的 HTTP 接口暴露。 | `http` | 需要 (复用 RPC) |
| **3. 泛 HTTP RESTful 服务**<br>(RESTful API) | 基于 `google.api.http` 注解，支持复杂的 URL 路径映射 (如 `/users/{id}`)。 | 对外开放的 REST API、需要精细化 URL 设计的场景。 | `http` (配合 RESTful 插件) | 需要 (PB + 注解) |

> **注意**：`protocol: http` 与 `protocol: http_no_protocol` 的关键区别在于底层 Codec 是否启用 `AutoReadBody`。前者自动读取，后者需手动流式读取。

---

示例代码:`../examples/http`

---


# 泛HTTP服务搭建 (Native HTTP)

本节主要介绍上述第 1 类：**泛 HTTP 标准服务**。它的核心思想是将标准的 `http.Handler` 挂载到 trpc 的 Service 上，从而复用 trpc 的治理能力（日志、监控、服务发现等）。

### 1. 核心 API

```go
// 将任意实现了 http.Handler 接口的对象注册到 trpc Service
import (
    thttp "trpc.group/trpc-go/trpc-go/http"
)

thttp.RegisterNoProtocolServiceMux(service, mux)
```

### 2. 路由注册方式演进

#### 方式 A：标准库 `net/http` (推荐 Go 1.22+)
最轻量，无第三方依赖。Go 1.22+ 已增强了路由匹配能力。

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello"))
})
trpc_http.RegisterNoProtocolServiceMux(service, mux)
```

#### 方式 B：Gorilla Mux (官方示例)
曾是 Go 社区事实标准，功能强大但目前维护状态一般。官方文档示例中大量使用。

```go
r := mux.NewRouter()
r.HandleFunc("/user/{id}", handler).Methods("GET")
trpc_http.RegisterNoProtocolServiceMux(service, r)
```

#### 方式 C：集成 Gin 框架 (生产推荐)
利用 Gin 强大的参数绑定、校验和 JSON 渲染能力，提升开发效率。

```go
// 1. 创建 Gin Engine
g := gin.Default()

// 2. 定义路由与业务逻辑
g.GET("/users", func(c *gin.Context) {
    var req UserReq
    // 自动绑定 Query 或 JSON
    if err := c.ShouldBind(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"data": "success"})
})

// 3. 注册到 trpc
trpc_http.RegisterNoProtocolServiceMux(service, g)
```

### 3. 配置注意
在 `trpc_go.yaml` 中，务必将协议设置为 `http_no_protocol`，以避免框架尝试自动解析 Body 导致冲突。

```yaml
server:
  service:
    - name: trpc.demo.http.MyService
      protocol: http_no_protocol  # <--- 关键配置
      port: 8080
```

---

# 泛HTTP RPC 服务(with_protocol)

本节介绍上述第 2 类：**泛 HTTP RPC 服务** (Standard RPC over HTTP)。
它的核心理念是 **"Write RPC, Serve HTTP"**：你只需要编写标准的 RPC 服务代码，框架会自动将 HTTP 请求映射为 RPC 调用，并处理 JSON 到 PB 结构体的序列化。

### 1. 核心流程
1.  **定义 Proto**：编写标准的 `.proto` 文件。
2.  **生成代码**：使用 `trpc` 工具生成桩代码。
3.  **配置双协议**：在 `trpc_go.yaml` 中配置两个 Service 入口，一个跑 `trpc` 协议，一个跑 `http` 协议。
4.  **复用实现**：将同一个 `ServerImpl` 注册到这两个 Service 上。

### 2. 配置文件关键点 (`trpc_go.yaml`)

```yaml
server:
  service:
    # 1. 标准 RPC 服务 (供微服务内部调用)
    - name: trpc.helloworld.GreeterRPC
      protocol: trpc
      port: 8000

    # 2. HTTP 服务 (供外部调用)
    - name: trpc.helloworld.Greeter
      protocol: http     # <--- 关键：启用 RPC-Style 模式 (自动 JSON 转换)
      port: 8080
```

> **特别注意**：为了确保 HTTP URL 路由能正确匹配，建议 **HTTP Service 的 name 必须与 Proto Package 定义的 Service Name 完全一致**。如果名字不一致，可能导致框架无法正确解析 URL 路由。

### 3. 代码实现 (`main.go`)

```go
func main() {
    s := trpc.NewServer()
    impl := &GreeterImpl{}

    // 注册到 RPC Service
    pb.RegisterGreeterService(s.Service("trpc.helloworld.GreeterRPC"), impl)

    // 注册到 HTTP Service
    // 注意：这里使用的是与 proto package 一致的名字 "trpc.helloworld.Greeter"
    pb.RegisterGreeterService(s.Service("trpc.helloworld.Greeter"), impl)

    s.Serve()
}
```

### 4. 请求路由规则
默认情况下，URL 路径必须严格遵循以下格式：
```
http://<host>:<port>/<Package>.<Service>/<Method>
```

例如：
- Proto 定义：`package trpc.helloworld; service Greeter { rpc Hello... }`
- 请求 URL：`POST http://127.0.0.1:8080/trpc.helloworld.Greeter/Hello`
- Body：`{"msg": "world"}` (自动映射到 `HelloRequest`)

> **提示**：如果希望自定义 URL (如 `/hello`)，可以使用 `trpc.alias` 选项（需引入 `trpc.proto`），或者进阶使用 RESTful 模式（见下一节）。

---

# 泛HTTP RPC Restful 服务(with_protocol_restful)

