package main

import (
	"context"
	"errors"
	"time"

	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/http"
	"trpc.group/trpc-go/trpc-go/log"
)

// InitFilters 注册所有拦截器
func InitFilters() {
	filter.Register("demo_logger", DemoLoggerFilter, nil)
	filter.Register("demo_auth", DemoAuthFilter, nil)
}

// DemoLoggerFilter 通用日志拦截器
// 这个拦截器对 RPC 和 HTTP 都适用，因为它只关心 context 和耗时
func DemoLoggerFilter(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (interface{}, error) {
	start := time.Now()

	// 继续执行后续链路
	rsp, err := next(ctx, req)

	cost := time.Since(start)
	log.Infof("[Filter] Request processed in %v", cost)

	return rsp, err
}

// DemoAuthFilter HTTP 专用鉴权拦截器
// 演示如何从 HTTP Header 中获取信息
func DemoAuthFilter(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (interface{}, error) {
	// 关键点：从 context 中提取 HTTP 头部信息
	// trpc-go 将 http.Request 封装在了 Head 中
	head := http.Head(ctx)
	if head == nil || head.Request == nil {
		// 如果不是 HTTP 协议，可能直接放行或报错，视业务而定
		return next(ctx, req)
	}

	// 获取 Authorization Header
	token := head.Request.Header.Get("Authorization")

	// 简单的鉴权逻辑：Token 必须是 "secret-token"
	if token != "secret-token" {
		log.Warnf("[Filter] Unauthorized access attempt, token: %s", token)

		// 尝试更加友好的错误返回 (401 Unauthorized)
		// 注意：直接在这里写 ResponseWriter 是不推荐的，因为可能破坏洋葱模型。
		// 标准做法是返回特定错误，由 ErrorHandler 处理。
		// 这里简单返回 error，框架默认会转为 500，但在日志中我们可以清晰看到拦截原因。
		return nil, errors.New("unauthorized: invalid token")
	}

	log.Infof("[Filter] Auth passed for token: %s", token)
	return next(ctx, req)
}
