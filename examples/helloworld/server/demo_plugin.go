package main

import (
	"context"
	
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/plugin"
)

const (
	pluginType = "demo"
	pluginName = "logger"
)

func init() {
	plugin.Register(pluginName, &DemoPluginFactory{})
}

type DemoPluginFactory struct{}

func (f *DemoPluginFactory) Type() string {
	return pluginType
}

func (f *DemoPluginFactory) Setup(name string, dec plugin.Decoder) error {
	var cfg struct {
		Prefix string `yaml:"prefix"`
	}
	if err := dec.Decode(&cfg); err != nil {
		return err
	}

	filter.Register(pluginName, func(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (rsp interface{}, err error) {
		log.Infof("%s Request received!", cfg.Prefix)
		return next(ctx, req)
	}, nil)

	log.Infof("Demo plugin setup success: %s", name)
	return nil
}
