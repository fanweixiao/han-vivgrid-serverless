# han-vivgrid-serverless

## 1. vivgrid控制台

https://dashboard.vivgrid.com

修改System Prompt：
你是一个网络设备专家，用户可能会向你询问有关网关、路由器等相关问题，请用中文回答。
回答中不要涉及具体品牌的产品，仅向用户回答一般的网络相关知识。

## 2. 启动mock服务

```sh
pip install uvicorn fastapi pydantic

uvicorn mock_api:app --port 9999
```

## 3. 创建访客SFN

```sh
cd guest_account

export YOMO_SFN_NAME=guest_account
export YOMO_SFN_ZIPPER="zipper.vivgrid.com:9000"
export YOMO_SFN_CREDENTIAL="app-key-secret:****.****"

yomo run app.go
```

```sh
curl https://api.vivgrid.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ****.****" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role": "user", "content": "请帮我创建一个新的访客账号，用户名是“李明”"}]
  }'
```

```txt
已为您创建了一个新的访客账号，用户名为“李明”。用户ID是“7be36e7b-0842-4020-b1b1-5ff225f37fc0”，初始密码是“By_emXo-”。请尽快登录并修改密码以确保安全。
```

## 4. 设备RAG SFN

```sh
cd device_info

zip app.zip app.go *.txt

cp yc.yml.example yc.yml
# 编辑yc.yml，填入app-key和app-secret

yc deploy app.zip
```

```sh
curl https://api.vivgrid.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ****.****" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role": "user", "content": "HAN设备「AP211」可以支持5G频段吗？"}]
  }'
```

```txt
是的，HAN设备「AP211」支持5G频段。根据设备文档，AP211是一款室内802.11ac MU-MIMO AP，能够同时工作在2.4GHz和5GHz双频段。5GHz频段的最大无线速率可达867Mbps。
```
