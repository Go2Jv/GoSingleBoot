---
name: gosingleboot
description: GoSingleBoot 项目的项目级开发规范。用于本项目中新增或修改业务代码、Handler、Router、项目已有封装、第三方服务和基础设施。
---

# GoSingleBoot

在现有 GoSingleBoot 架构内实现需求。优先理解、搜索并复用项目已有能力；不要为了新需求重新设计框架。

## 1. 开发原则

按以下优先级工作：

```text
现有项目能力
    ↓
已有模块/模式
    ↓
新增文件或模块
    ↓
修改业务代码
    ↓
只有确实必要时才修改核心代码
```

实现需求前先搜索代码，确认是否已经存在可以复用的错误处理、JWT、DTO、Response、DB、Redis、Logger、Config、Middleware 或其他项目封装。

已存在相同职责的能力必须优先复用。不要重新实现相同逻辑，也不要创建职责重复的 Utility / Helper。

## 2. 依赖边界

### 稳定基础设施

稳定、长期固定、没有现实供应商替换需求的依赖直接使用项目已有实例：

```text
PostgreSQL / Redis / Kafka / RabbitMQ / NATS
Logger / Config / JWT
```

例如：

```go
db.Client.Master
redis.Client
logger.Logger
config.Config
```

不要为了统一 DI 而创建：

```text
IDatabase / IRedis / ILogger / IConfig / IJWT / IKafka / IRabbitMQ
```

### 可替换服务

存在现实供应商替换需求时，通过 Interface 隔离：

```text
Email / SMS / AI / Payment / Object Storage
Translation / Search / Maps / External APIs
```

依赖方向：

```text
Handler → Interface → Provider → 第三方 API / SDK
```

Interface：

```text
internal/Interfaces/
```

Provider：

```text
internal/<Service>/<Provider>/
```

例如：

```text
internal/
├── Interfaces/
│   ├── IEmail.go
│   └── ISMS.go
├── Email/
│   ├── GoogleEmail/
│   └── AzureEmail/
└── SMS/
    ├── GoogleSMS/
    └── AzureSMS/
```

同一替换边界已经存在 Interface 时直接复用。

## 3. bizErr 与错误控制流

业务代码优先使用项目已有 `bizErr`，不要重复展开相同的 `if err != nil` 逻辑。

约定：

```text
检查类方法返回 true → 没有错误，可以继续
检查类方法返回 false → 错误已经处理，当前 Handler 必须 return
```

优先写：

```go
if !bizErr.Validation(c, "非法传参", err) {
    return
}

if !bizErr.SQLNotFound(c, "用户不存在", err) {
    return
}

if !bizErr.DBError(c, err) {
    return
}
```

项目已有错误检查封装时，不要重新写等价的：

```go
if err != nil {
    ...
}
```

当代码已经明确知道这里就是业务错误，使用项目的直接错误提交方法，然后显式 `return`：

```go
bizErr.Abort(c, 400, "账号密码错误", nil)
return
```

正常业务错误应通过 `*BizErr` 进入 `c.Error()`，由全局错误中间件统一响应和记录日志。不要在 Handler 中对已经交给全局错误中间件的错误重复打日志。

## 4. 业务模块结构

每个业务模块拥有自己的 Handler 和 Router。

以 User 为例：

```text
internal/
├── handler/
│   └── userHandler.go
└── router/
    └── userRouter.go
```

实现新模块时遵循：

```text
检查已有能力
    ↓
创建 <name>Handler.go
    ↓
定义 <Name>Handler
    ↓
放入实际需要的可替换服务 Interface
    ↓
需要时创建 New<Name>Handler(...)
    ↓
实现 Info / Orders / ... Endpoint 方法
    ↓
创建 <name>Router.go
    ↓
Register<Name>Router(rg) 中准备具体依赖/Provider
    ↓
New<Name>Handler(...) 完成注入
    ↓
在 rg 上注册 Endpoint
    ↓
Main Router 调用 Register<Name>Router(rg)
```

## 5. Handler

Handler 只放**实际需要注入**的可替换服务。

```go
type UserHandler struct {
    email Interfaces.IEmail
    sms   Interfaces.ISMS
}

func NewUserHandler(
    email Interfaces.IEmail,
    sms Interfaces.ISMS,
) *UserHandler {
    return &UserHandler{email: email, sms: sms}
}

func (h *UserHandler) Info(c *gin.Context) {
    // ...
}

func (h *UserHandler) Orders(c *gin.Context) {
    // ...
}
```

规则：

- 稳定基础设施按项目现有方式直接使用，不为了统一 DI 强行塞进 struct。
- 可替换服务使用 Interface，不使用具体 Provider。
- Handler 不创建第三方 Client。
- 一个业务模块通常共用一个 Handler，不为每个 Endpoint 单独创建 Handler。
- Router 与 Handler 不同 package 时，Endpoint 方法必须导出，例如 `Info`、`Orders`。
- Endpoint 内优先使用项目已有封装，如 `bizErr`、JWT、DTO 等。

## 6. Router

业务模块必须有独立 Router 文件：

```text
internal/router/userRouter.go
```

入口：

```go
func RegisterUserRouter(rg *gin.RouterGroup) {
    // 准备依赖
    // 创建 Handler
    // 注册 Endpoint
}
```

`Register<Name>Router()` 是模块的**依赖组装点**。具体 Provider 在这里选择，再以 Interface 传给 Handler：

```go
func RegisterUserRouter(rg *gin.RouterGroup) {
    email := /* GoogleEmail / AzureEmail */
    sms := /* GoogleSMS / AzureSMS */

    userHandler := handler.NewUserHandler(email, sms)

    rg.GET("/info", userHandler.Info)
    rg.GET("/orders", userHandler.Orders)
}
```

不要把业务 Endpoint 直接写到 `mainRouter.go`。

## 7. Main Router

`mainRouter.go` 只负责组合业务 Router，不实现业务 Endpoint：

```go
RegisterUserRouter(rg)
RegisterOrderRouter(rg)
```

保持现有结构：

```text
MainRouter
    ↓
gin.Engine
    ↓
Middleware / Docs / Groups
    ↓
Register<Name>Router(...)
```

## 8. 核心代码与基础设施

默认谨慎修改：

```text
internal/logger/
internal/db/
internal/config/
internal/docs/
internal/jwt/
config.json
cmd/main.go
```

只有需求确实要求时才修改。

新增 Redis、Kafka、RabbitMQ、NATS 等基础设施时，优先创建独立 Package，只做必要的初始化/启动接入，不要顺手重构 logger、config、DB、router、middleware 或 DI。

## 9. 禁止

```text
重新设计 GoSingleBoot 的整体架构
模仿 Spring Boot / NestJS / FastAPI 等框架重写现有结构
为稳定基础设施创建没有真实替换需求的 Interface
Handler 直接依赖具体 Provider / 第三方 SDK
Handler 内创建第三方 Client
把业务 Endpoint 集中到 mainRouter.go
重复实现项目已有的 bizErr / JWT / Response / Utility 能力
为了“架构完整”增加无实际需求的 Repository / Service / Manager / Factory
创建 utils/、common/、helper/、service/、thirdparty/ 等万能目录
实现小功能时顺手重构无关代码
```

## 10. 完成前检查

```text
[ ] 已先搜索并复用项目已有能力
[ ] 没有重复实现已有封装
[ ] 依赖已区分稳定基础设施与可替换服务
[ ] 可替换服务使用 Interface
[ ] Handler / Router 位于正确目录
[ ] Handler 只包含实际需要的注入依赖
[ ] Provider 在 Register<Name>Router() 中组装
[ ] 业务 Endpoint 由模块 Router 注册
[ ] mainRouter.go 只组合模块 Router
[ ] bizErr 检查函数遵循 true=继续、false=return
[ ] 已交给全局错误处理的错误没有重复日志
[ ] 没有无关核心重构
```
