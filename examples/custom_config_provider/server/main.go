package main

import (
	"fmt"
	"time"

	"trpc.group/trpc-go/trpc-go/config"
	"trpc.group/trpc-go/trpc-go/log"
)

func init() {
	// 注册我们的自定义插件
	config.RegisterProvider(NewMockProvider())
}

func main() {
	// 1. 初始化远程配置中心的初始值
	UpdateRemoteConfig("app.yaml", `
server:
  timeout: 1000
  msg: "initial version"
`)

	// 2. 加载配置
	// 指定使用我们刚才注册的 "mock-remote" provider
	cfg, err := config.Load("app.yaml",
		config.WithProvider("mock-remote"),
		config.WithCodec("yaml"),
		config.WithWatch(),
	)
	if err != nil {
		panic(err)
	}

	// 3. 启动一个协程，模拟远程配置每 3 秒变一次
	go func() {
		version := 1
		for {
			time.Sleep(3 * time.Second)
			version++
			newValue := fmt.Sprintf(`
server:
  timeout: %d
  msg: "version %d"
`, 1000+version, version)

			// 模拟在远程控制台点击了“发布”
			UpdateRemoteConfig("app.yaml", newValue)
		}
	}()

	// 4. 主循环：观察配置是否自动变了
	for {
		timeout := cfg.GetInt("server.timeout", 0)
		msg := cfg.GetString("server.msg", "")

		log.Infof("[Main] Current Config -> timeout: %d, msg: %s", timeout, msg)
		time.Sleep(1 * time.Second)
	}
}
