package main

import (
	"time"

	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	// 1. 配置日志
	// 模拟 trpc_go.yaml 中的配置结构
	cfg := log.Config{
		// 输出源 1: 控制台
		{
			Writer:    "console",
			Level:     "debug",
			Formatter: "console",
		},
		// 输出源 2: 文件
		{
			Writer:    "file",
			Level:     "info", // 只记录 Info 及以上
			Formatter: "json", // 使用 JSON 格式
			WriteConfig: log.WriteConfig{
				Filename:   "trpc.log",
				RollType:   "size", // 按大小切割
				MaxSize:    1,      // 1MB
				MaxBackups: 3,      // 保留最近 3 个文件
				Compress:   false,  // 不压缩
			},
		},
	}

	// 2. 初始化并替换全局 Logger
	logger := log.NewZapLog(cfg)
	log.SetLogger(logger)

	// 3. 打印日志
	for i := 0; i < 5; i++ {
		// Debug: 只会在控制台显示 (因为文件的 Level 是 Info)
		log.Debugf("Debug log (console only): %d", i)

		// Info: 会同时显示在控制台和文件
		log.Infof("Info log (console + file): %d", i)

		// Error: 带结构化字段
		log.With(log.Field{Key: "uid", Value: 123}).Error("Error log with field")

		time.Sleep(100 * time.Millisecond)
	}
}
