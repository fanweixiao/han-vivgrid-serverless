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

var (
	HanUrl = "http://127.0.0.1:9999"
)

type Parameter struct {
	Username string `json:"username" jsonschema:"description=the guest account username"`
}

func Description() string {
	return `this function is used to create a new guest account`
}

func InputSchema() any {
	return &Parameter{}
}

func Init() error {
	if v, ok := os.LookupEnv("HAN_URL"); ok {
		HanUrl = v
	}
	return nil
}

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

func DataTags() []uint32 {
	return []uint32{0x10}
}

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
