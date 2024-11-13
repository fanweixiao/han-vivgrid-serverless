package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/yomorun/yomo/serverless"
)

// 这里定义的是 API 地址
var (
	HanUrl = "http://127.0.0.1:9999"
	// HanUrl = "https://192.168.40.183"
)

// 定义 LLM 调用该 Function 时的参数签名
type Parameter struct {
	Username string `json:"username" jsonschema:"description=the guest account username"`
}

// 必要方法。准确的描述该 Function 的功能，有助于 LLM 匹配用户的问题
func Description() string {
	return `this function is used to create a new guest account`
}

// 必要方法
func InputSchema() any {
	return &Parameter{}
}

// 非必要方法，用于初始化工作
func Init() error {
	if v, ok := os.LookupEnv("HAN_URL"); ok {
		HanUrl = v
	}
	return nil
}

// 必要方法。每次 LLM 调用该 Function Calling 时，该函数会被唤醒运行一次
func Handler(ctx serverless.Context) {
	slog.Info("[sfn] receive", "ctx.data", string(ctx.Data()))

	var msg Parameter
	err := ctx.ReadLLMArguments(&msg)
	if err != nil {
		slog.Error("[sfn] json.Marshal error", "err", err)
		return
	}

	id, password, err := createGuestAccount(msg.Username)
	if err != nil {
		slog.Error("[sfn] >> createGuestAccount error", "err", err)
		return
	}

	result := fmt.Sprintf("已为【%s】创建新访客账号，用户ID为【%s】，初始密码为【%s】", msg.Username, id, password)

	err = ctx.WriteLLMResult(result)
	if err != nil {
		slog.Error("[sfn] >> write error", "err", err)
		return
	}
}

// 必要方法
func DataTags() []uint32 {
	return []uint32{0x10}
}

//////////// 以下均为调用 API 时的业务代码，与 Function Calling 本身无关 ////////////

type GuestAccountRequest struct {
	Username              string `json:"username"`
	Password              string `json:"password"`
	ServiceLevel          string `json:"serviceLevel"`
	AccountValidityPeriod int64  `json:"accountValidityPeriod"`
}

type Response[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type GuestAccountResponse struct {
	Id string `json:"id"`
}

func createGuestAccount(username string) (string, string, error) {
	url := HanUrl + "/api/v1.0/am/accounts/guest/accounts"

	password, err := gonanoid.New(8) // 生成一个初始密码
	if err != nil {
		return "", "", err
	}

	reqBody := &GuestAccountRequest{
		Username:              username,
		Password:              password,
		ServiceLevel:          "Default Service Level",
		AccountValidityPeriod: time.Now().AddDate(1, 0, 0).UnixMilli(), // 默认一年有效期
	}

	buf, err := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", "", errors.New("bad status: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var res Response[GuestAccountResponse]
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", "", err
	}

	return res.Data.Id, password, nil
}
