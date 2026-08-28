# GoSingleBoot Development Rules Skill

## 1. Minimal Change Principle

任何需求都应遵循最小修改原则。

优先级：

```text
新增文件
↓
修改业务代码
↓
修改核心代码
```

禁止为了完成普通业务需求而重构框架。

---

## 2. Handler Rules

### Handler 可以直接使用

稳定基础设施可以直接使用，例如：

```text
DB
Redis
Logger
Config
JWT
Kafka
RabbitMQ
NATS
```

### Handler 必须依赖 Interface

以下可替换服务必须通过 Interface 注入：

```text
IEmail
ISMS
IAI
IPayment
IStorage
```

示例：

```go
type UserHandler struct {
    email Interfaces.IEmail
}
```

---

## 3. Handler Forbidden

### 禁止直接依赖 Provider

错误：

```go
type Handler struct {
    email *GoogleEmail.Client
}
```

正确：

```go
type Handler struct {
    email Interfaces.IEmail
}
```

### 禁止 Handler 创建第三方 Client

禁止：

```go
func Login() {
    client := GoogleClient.New()
}
```

原因：

业务层不应该直接感知供应商实现。

---

## 4. Interface Location

所有业务可替换服务的 Interface：

```text
internal/Interfaces/
```

例如：

```text
internal/
└── Interfaces/
    ├── IEmail.go
    ├── ISMS.go
    └── IAI.go
```

---

## 5. Provider Location

第三方实现应放在独立目录。

推荐结构：

```text
internal/
├── Email/
│   ├── GoogleEmail/
│   └── AzureEmail/
├── SMS/
│   ├── GoogleSMS/
│   └── AzureSMS/
├── AI/
│   ├── OpenAI/
│   └── Gemini/
└── Payment/
    └── Stripe/
```

---

## 6. Provider Naming

Provider 命名应体现**供应商 + 服务类型**。

正确：

```text
GoogleEmail
AzureSMS
OpenAI
Stripe
```

禁止使用过于模糊的目录或类型名：

```text
utils
common
helper
thirdparty
service
```

---

## 7. Router Rules

Router 必须使用：

```text
internal/router/
```

不要重新设计 Router。

当前模式：

```text
MainRouter()
   ↓
gin.Engine
   ↓
Middleware
   ↓
Docs
   ↓
API Group
   ↓
RegisterXXXRouter()
```

---

## 8. New Router Module

新增路由模块时，推荐：

```text
internal/router/user_router.go
```

例如：

```go
func RegisterUserRouter(
    rg *gin.RouterGroup,
) {
    userHandler := handler.NewUserHandler()

    rg.GET(
        "/users",
        userHandler.List,
    )
}
```

然后在主路由注册：

```go
MainRouter()

RegisterUserRouter(rg)
```

---

## 9. Core Files Protection

默认禁止修改以下核心文件或模块，除非需求确实要求：

```text
internal/logger/
internal/db/
internal/config/
internal/docs/
internal/jwt/
config.json
cmd/main.go
```

修改前应先判断是否可以通过新增文件或新增模块完成需求。

---

## 10. Infrastructure Addition

新增以下基础设施：

```text
Redis
Kafka
RabbitMQ
NATS
```

必须创建独立 Package，例如：

```text
internal/redis/
internal/mq/
```

允许：

```go
redis.Init()
```

并在：

```text
cmd/main.go
```

中增加初始化调用。

但只能做必要的初始化接入，禁止顺带：

- 修改启动流程
- 修改 Logger
- 修改 Config
- 修改 DB
- 修改 Router
- 重构 DI

---

## 11. Before Coding Checklist

开发前依次检查：

```text
1. 需求是什么？

2. 属于业务还是基础设施？

3. 是否属于第三方供应商？

4. 是否需要 Interface？

5. Interface 放在哪里？

6. Provider 放在哪里？

7. Handler 是否需要注入？

8. Router 如何注册？

9. 是否可以通过新增文件解决？

10. 是否必须修改核心文件？
```

最后必须优先选择**最小影响范围的实现方案**。
