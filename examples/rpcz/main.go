package main

import (
	"context"
	"fmt"
	"time"

	"trpc.group/trpc-go/trpc-go/rpcz"
)

func main() {
	// 1. 开启 rpcz (默认可能是关闭的)
	// 在 trpc_go.yaml 中通常通过 global.rpcz 配置
	// 这里我们手动初始化一个 Global RPCZ
	// 这一步其实创建了一个全局的 SpanStore，用来暂存采集到的 Span
	rpcz.GlobalRPCZ = rpcz.NewRPCZ(&rpcz.Config{
		Fraction: 1.0, // 采样率 100% (全采)
		Capacity: 100, // 内存中保留最近 100 条
	})

	fmt.Println("=== Start Trace: Client -> Service A -> Service B ===")

	// 2. Client 发起请求 (Root Span)
	ctx := context.Background()
	// NewSpanContext 会检查 ctx 里有没有父 Span，没有就新建一个 Root Span
	span, end, ctx := rpcz.NewSpanContext(ctx, "Client.Call")

	// 模拟 Client 处理耗时
	time.Sleep(10 * time.Millisecond)
	span.SetAttribute("user_id", "10086") // 打个标签

	fmt.Printf("1. Client Start (SpanID: %d)\n", span.ID())

	// 3. 调用 Service A (传递 ctx)
	callServiceA(ctx)

	// 4. Client 结束
	end.End()
	fmt.Println("4. Client End")
}

func callServiceA(ctx context.Context) {
	// 从 ctx 继承 Span (成为 Client 的子 Span)
	// 这里的 ctx 已经携带了 Client 的 Span 信息
	span, end, ctx := rpcz.NewSpanContext(ctx, "ServiceA.Handle")
	defer end.End()

	fmt.Printf("  -> 2. Service A Start (SpanID: %d)\n", span.ID())
	time.Sleep(20 * time.Millisecond)

	// 调用 Service B
	callServiceB(ctx)
}

func callServiceB(ctx context.Context) {
	// 从 ctx 继承 Span (成为 Service A 的子 Span)
	span, end, _ := rpcz.NewSpanContext(ctx, "ServiceB.Handle")
	defer end.End()

	fmt.Printf("    -> 3. Service B Start (SpanID: %d)\n", span.ID())
	time.Sleep(50 * time.Millisecond)
	span.SetAttribute("db.query", "select * from users")
}
