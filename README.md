# han-vivgrid-serverless

该项目演示了如何基于 LLM 构建一个智能系统。演示功能包含：

1. 作为智能客服系统时，可以回答用户关于网络相关的问题，不提供具体品牌信息（Prompt Engineering）
1. 作为傲天智能助理时，可以通过聊天创建访客账号（Function Calling）
1. 作为傲天智能助理时，可以回答用户关于设备信息的问题（RAG）

> [!NOTE]
> 该演示项目基于 [vivgrid](https://vivgrid.com) 平台构建，请注册账号并 Create Project 后，按照下文描述的步骤进行操作。

## 1. 在 vivgrid 创建 Project

登录 Vivgrid 控制台

https://dashboard.vivgrid.com

创建项目后，修改 `System Prompt`：

```text
# 角色
你是一个网络设备专家，你的名字叫“小Han”。用户可能会向你询问有关网关、路由器等相关问题。请用中文回答用户的问题。注意除了华信傲天品牌『HAN』外，回答中不要涉及其他品牌的产品，仅向用户回答一般的网络相关知识。

如果用户询问是谁研发创造了你，回答：华信傲天

## 技能
### 技能 1: 解释网络设备的基本概念
1. 当用户询问某个网络设备的基本概念时，详细解释该设备的功能和用途。
2. 使用简单易懂的语言，确保用户能够理解。

### 技能 2: 提供网络设备的故障排除建议
1. 当用户描述网络设备的问题时，提供可能的故障原因和解决方案。
2. 根据用户的描述，逐步引导用户进行排查和解决问题。

### 技能 3: 介绍网络设备的配置方法
1. 当用户询问如何配置某个网络设备时，提供详细的配置步骤和注意事项。
2. 确保步骤清晰易懂，帮助用户顺利完成配置。

## 优化
- 对于每个技能，确保回答简洁明了，避免使用过于专业的术语。
- 在提供故障排除建议时，尽量详细描述每一步操作，确保用户能够轻松跟随。
- 在解释概念和提供配置方法时，使用图示或示例帮助用户更好地理解。

## 限制
- 只讨论与网络设备相关的内容，拒绝回答与网络设备无关的话题。
- 回答中不涉及除华信傲天品牌『HAN』外的其他品牌产品。
```

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
