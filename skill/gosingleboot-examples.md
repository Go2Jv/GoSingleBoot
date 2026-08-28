# GoSingleBoot Coding Examples Skill

## 1. Third-Party Dependency

### Bad

```go
type LoginHandler struct {
    email *GoogleEmail.Client
}
```

问题：

- Handler 绑定供应商
- 无法低成本切换 Provider
- 业务层与第三方 SDK 耦合

### Good

```go
type LoginHandler struct {
    email Interfaces.IEmail
}
```

依赖关系：

```text
Handler
   ↓
IEmail
   ↓
GoogleEmail
```

---

## 2. Creating Client Inside Handler

### Bad

```go
func Login() {
    email := GoogleEmail.NewClient()
    email.Send()
}
```

问题：

业务层直接感知供应商。

### Good

```go
func Login() {
    h.email.Send()
}
```

业务层只依赖抽象接口。

---

## 3. Fake Abstraction

### Bad

无真实替换需求，却创建：

```text
IDatabase
IRedis
ILogger
```

问题：

- 增加类型数量
- 增加 DI / Mock / 维护成本
- 没有真实的业务收益

### Good

稳定基础设施直接使用：

```go
db.Client.Master
```

---

## 4. Wrong Third-Party Directory

### Bad

```text
internal/
└── service/
    ├── email.go
    └── sms.go
```

问题：

不同供应商和不同职责容易混杂，扩展时边界不清晰。

### Good

```text
internal/
├── Interfaces/
│   └── IEmail.go
└── Email/
    ├── GoogleEmail/
    └── AzureEmail/
```

---

## 5. Router Redesign

### Bad

为了增加一个 API，新建完全不同的路由系统：

```text
routes/
router_manager/
module_router/
```

原因：

这会破坏 GoSingleBoot 现有 Router 设计。

### Good

保持：

```text
internal/router/
├── user_router.go
└── login_router.go
```

按照已有模式注册新的 Router。

---

## 6. Framework Rewrite

### Bad

为了新增 Redis，同时修改：

```text
logger
config
db
middleware
router
main.go
```

问题：

需求范围很小，却产生大范围修改，容易引入回归问题。

### Good

新增独立模块：

```text
internal/redis/
└── redis.go
```

然后仅在启动入口增加必要初始化：

```go
redis.Init()
```

---

## 7. Code Review Questions

### Dependency

- Handler 是否依赖具体 Provider？
- 第三方服务是否通过 Interface？
- 是否创建了无意义的 Interface？

### Architecture

- 是否修改了框架核心设计？
- 是否引入了其他框架的思想？
- 是否违反了 Router 规范？

### Scope

- 是否可以通过新增文件解决？
- 是否修改了无关文件？
- 是否影响核心模块？

---

## Final Rule

永远保持：

```text
Stable Infrastructure
        ↓
   Direct Use
```

以及：

```text
Replaceable Service
        ↓
    Interface
        ↓
    Provider
        ↓
 Third Party API
```

禁止：

```text
Handler
   ↓
GoogleEmail
```

禁止：

```text
Handler
   ↓
AzureSMS
```

核心目标：**让业务依赖稳定，让可替换部分真正可替换，同时保持修改范围最小。**
