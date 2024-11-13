package main

import (
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/yomorun/yomo/serverless"
)

var (
	//go:embed AP211.txt
	docAP211 string

	//go:embed AP271.txt
	docAP271 string
)

type Parameter struct {
	DeviceId string `json:"username" jsonschema:"description=网络设备型号ID，例如：AP211、AP271等"`
}

func Description() string {
	return `这个函数将根据用户提出的有关 HAN（傲天) 网络设备问题，返回具体设备型号的知识文档。
以下是一些用户可能提出的问题：
1.请概要介绍一下 AP211 的使用场景？
2.AP271 是否支持5G频段？`
}

func InputSchema() any {
	return &Parameter{}
}

func Handler(ctx serverless.Context) {
	slog.Info("[sfn] receive", "ctx.data", string(ctx.Data()))

	var msg Parameter
	err := ctx.ReadLLMArguments(&msg)
	if err != nil {
		slog.Error("[sfn] json.Marshal error", "err", err)
		return
	}

	var doc string
	switch msg.DeviceId {
	case "AP211":
		doc = docAP211
	case "AP271":
		doc = docAP271
	default:
		ctx.WriteLLMResult(fmt.Sprintf("抱歉，我没有找到关于「%s」产品的文档", msg.DeviceId))
	}

	ctx.WriteLLMResult(fmt.Sprintf("以下内容是 HAN 设备「%s」的文档，请根据文档内容回答用户问题：\n%s", msg.DeviceId, doc))
}

func DataTags() []uint32 {
	return []uint32{0x11}
}
