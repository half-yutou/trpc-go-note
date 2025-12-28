package main

import (
	"math/rand"
	"time"

	"trpc.group/trpc-go/trpc-go/metrics"
)

func main() {
	// 1. 注册 Sink (数据接收端)
	// 我们用 ConsoleSink，它会把数据直接打印到屏幕上
	// 在生产环境中，这里通常会注册 PrometheusSink 或其他监控系统的 Sink
	metrics.RegisterMetricsSink(metrics.NewConsoleSink())

	// 模拟奶茶店营业
	go func() {
		for {
			// --- 场景 1: Counter (卖出一杯) ---
			// 每次调用 Incr()，计数器 +1
			metrics.Counter("tea.sold.total").Incr()

			// --- 场景 2: Gauge (当前排队人数) ---
			// Set() 设置当前值，可以忽高忽低
			currentQueue := float64(rand.Intn(10))
			metrics.Gauge("tea.queue.size").Set(currentQueue)

			// --- 场景 3: Timer (制作耗时) ---
			// Record() 记录一次耗时
			// 模拟制作耗时 100ms - 500ms
			cost := time.Duration(100+rand.Intn(400)) * time.Millisecond
			metrics.Timer("tea.make.cost").RecordDuration(cost)

			time.Sleep(1 * time.Second)
		}
	}()

	// 阻塞主进程
	select {}
}
