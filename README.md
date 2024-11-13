# han-vivgrid-serverless

该项目演示了如何基于 LLM 构建一个智能系统。演示功能包含：

1. 作为智能客服系统时，可以回答用户关于网络相关的问题，不提供具体品牌信息（Prompt Engineering）
1. 作为傲天智能助理时，可以通过聊天创建访客账号（Function Calling）
1. 作为傲天智能助理时，可以回答用户关于设备信息的问题（RAG）

基于开源的 NextChat 在 Vercel 上搭建了一套测试用的界面，点击访问：https://han-vivgrid.vercel.app/

（若不可访问，请挂🪜）

> [!TIP]
> 该演示项目基于 [vivgrid](https://vivgrid.com) 平台构建，请注册账号并 Create Project 后，按照下文描述的步骤进行操作。

## 0. 在 vivgrid 创建 Project

登录 Vivgrid 控制台

https://dashboard.vivgrid.com

## 1. 基于通用知识的智能客服助手：“AP 和 WiFi 的区别是什么？”

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

在 Vivgrid 的界面上操作方法：

![image](https://github.com/user-attachments/assets/c3675fd3-5bbd-4b56-860a-79ed30742e39)

1. 🔴 红色箭头处选择 `Overwrite`，然后点击 `Save`。该设置会忽略所有的 API Request 中的 System Prompt 部分
2. 🔵 蓝色箭头指向 Prompt Evaluation 功能，方便您快速验证 Prompt 的修改是否有效。Prompt Engineering 就是反复调试知道最终可以向 LLM 清晰明了的表达准确意图。

该功能演示了 Prompt Engineering 技术，可以在 [OpenAI 官方的 Tutorial](https://platform.openai.com/docs/guides/prompt-engineering) 里查看更多细节。

## 2. 傲天私有 API 与 LLM 整合：“创建一个新的访客账号，用户名是：李明”

> [!NOTE]
> 演示的 API 是 `https://192.168.40.183/api/v1.0/`，但为了演示效果，我们搭建了一个 mock 服务。之后，可以随时将该部分代码在内网环境中启动。只要内网环境可以 udp 连接到 `zipper.vivgrid.com:9000`，即可使用。

[mock_api.py](./mock_api.py) 是 Mock API Server，首先本地启动它：

```sh
pip install uvicorn fastapi pydantic

uvicorn mock_api:app --port 9999
```

然后，以 Serverless 的方式编写 LLM Function Calling，让 LLM 知道当用户的要求是创建访客账号时，需要调用该程序（SFN）完成任务。
代码在 [sfn_create_guest_account](./sfn_create_guest_account) 目录下。

> [!TIP]
> 在私有环境中 Hosting 该 Function Calling 需要使用开源的 [YoMo](https://github.com/yomorun/yomo) 项目。

先安装 YoMo CLI:

```bash
curl -fsSL https://get.yomo.run | sh
```

在 Vivgrid Dashboard 上获取该项目的 Token：

![image](https://github.com/user-attachments/assets/173eb46d-966a-4b4b-bd0f-2fca688b3544)

之后，在机器上启动该 Function Calling：

```sh
cd sfn_create_guest_account

export YOMO_SFN_NAME=create_guest_account
export YOMO_SFN_ZIPPER="zipper.vivgrid.com:9000"
export YOMO_SFN_CREDENTIAL="app-key-secret:****.****"

yomo run app.go
```

启动之后，即可回到 Vivgrid Dashboard，在 AI Bridge 页面的 Prompt Evaluation 进行测试：

![image](https://github.com/user-attachments/assets/11618ff6-d3b6-4493-be91-b9c7a3d4288e)

在这里，我们看到 LLM 已经调用了 API，并完成了账号创建工作。

当然，也可以使用 API 访问，您完全可以使用 OpenAI API 接口兼容的 SDK 进行对 Vivgrid 的调用。

注：`Authorization` 中的 Token 的获得方式在上文已有描述。

```sh
curl https://api.vivgrid.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ****.****" \
  -d '{
    "model": "gpt-4o",
    "messages": [{"role": "user", "content": "请帮我创建一个新的访客账号，用户名是“李明”"}]
  }'
```

API 响应结果为：

```txt
已为您创建了一个新的访客账号，用户名为“李明”。用户ID是“7be36e7b-0842-4020-b1b1-5ff225f37fc0”，初始密码是“By_emXo-”。请尽快登录并修改密码以确保安全。
```

该功能演示了如何基于 Vivgrid 和开源的 YoMo 工具方便快速的构建一个 LLM Function Calling 服务，并已 Serverless 的方式编写，这也使得之后的部署和运维工作变得异常容易。
如果想了解更多关于 Function Calling 的知识，可以访问 [OpenAI 官方文档 - Function Calling](https://platform.openai.com/docs/guides/function-calling)

> [!TIP]
> 主流 LLM 均提供了 Function Calling 功能，但不同的 LLM 定义的 Function Calling 规范不同，代码完全不可复用。可以对照 Google Gemini API 的官方文档查看与 OpenAI API 中对 Function Calling 的规范差异。
> 
> 我们的 Vivgrid 提供了对不同 LLM 的支持，这意味着当您编写的 Function Calling Serverless，在更换 LLM 时，无需任何改动！😃

借助 LLM Function Calling 功能，可以将各种 API 分别包装成服务，以拓展 LLM 与现场业务系统的结合。

## 3. 基于私有知识的智能助手："「AP211」可以支持5G频段吗？"

该功能演示如何将私有知识应用于回答用户问题。

`RAG` 技术在正式使用场景中，因为其技术缺陷，往往导致无法精准理解用户的问题，我们的 F500 客户使用了下面演示的方法精准的理解和回答用户问题，如下图：

![image](https://github.com/user-attachments/assets/ef9eb17d-c192-4aff-92e3-f6a79312ae2d)

该类内容因无需访问私有 API，因此我们可以将其部署至 Vivgrid，这将大幅降低您的运维管理成本：

```sh
cd sfn_device_info

zip app.zip app.go *.txt

cp yc.yml.example yc.yml

# 编辑 yc.yml，填入 app-key 和 app-secret

yc deploy app.zip
```

其中，`app-key` 和 `app-secret` 可以在 Vivgrid Dashboard - Configuration 中找到：

![image](https://github.com/user-attachments/assets/bde783bc-bd4a-400b-bccf-1907dd0a1564)

当然，该功能也支持使用 API 访问：

```sh
curl https://api.vivgrid.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ****.****" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role": "user", "content": "HAN设备「AP211」可以支持5G频段吗？"}]
  }'
```

响应：

```txt
是的，HAN设备「AP211」支持5G频段。根据设备文档，AP211是一款室内802.11ac MU-MIMO AP，能够同时工作在2.4GHz和5GHz双频段。5GHz频段的最大无线速率可达867Mbps。
```
