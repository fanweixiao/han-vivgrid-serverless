# han-vivgrid-serverless

## Usage

```sh
uvicorn mock_api:app --port 9999
```

```sh
export YOMO_SFN_NAME=han
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
