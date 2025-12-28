package main

import (
	"fmt"
	"sync"

	"trpc.group/trpc-go/trpc-go/config"
	"trpc.group/trpc-go/trpc-go/log"
)

// MockRemoteServer 模拟远程配置中心服务端
var MockRemoteServer = struct {
	data  map[string]string
	mu    sync.RWMutex
	chans []chan map[string]string // 模拟推送通道
}{
	data:  make(map[string]string),
	chans: make([]chan map[string]string, 0),
}

// UpdateRemoteConfig 模拟在远程配置中心修改配置
func UpdateRemoteConfig(key, value string) {
	MockRemoteServer.mu.Lock()
	defer MockRemoteServer.mu.Unlock()

	MockRemoteServer.data[key] = value
	fmt.Printf("[RemoteServer] Config updated: %s = %s\n", key, value)

	// 推送给所有连接的客户端
	for _, ch := range MockRemoteServer.chans {
		// 非阻塞推送，防止卡死
		select {
		case ch <- MockRemoteServer.data:
		default:
		}
	}
}

// ------------------------------------------------------------------

// MockProvider 实现 config.DataProvider 接口
type MockProvider struct {
	name string
}

func NewMockProvider() *MockProvider {
	return &MockProvider{name: "mock-remote"}
}

func (p *MockProvider) Name() string {
	return p.name
}

// Read 第一次加载时调用
func (p *MockProvider) Read(path string) ([]byte, error) {
	MockRemoteServer.mu.RLock()
	defer MockRemoteServer.mu.RUnlock()

	if val, ok := MockRemoteServer.data[path]; ok {
		return []byte(val), nil
	}
	return nil, fmt.Errorf("config not found: %s", path)
}

// Watch 监听变更
func (p *MockProvider) Watch(cb config.ProviderCallback) {
	// 1. 建立连接（这里用 channel 模拟长连接）
	updateCh := make(chan map[string]string, 10)

	MockRemoteServer.mu.Lock()
	MockRemoteServer.chans = append(MockRemoteServer.chans, updateCh)
	MockRemoteServer.mu.Unlock()

	// 2. 启动协程监听
	go func() {
		for data := range updateCh {
			// 假设我们只关心 "app.yaml" 这个配置
			// 在真实场景中，这里会根据 path 过滤
			if val, ok := data["app.yaml"]; ok {
				log.Infof("[MockProvider] Received update for app.yaml, triggering callback...")
				// 3. 核心：调用框架回调，更新内存
				cb("app.yaml", []byte(val))
			}
		}
	}()
}
