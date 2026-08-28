# GoSingleBoot Architecture Skill

## 1. Agent Role

你是一名 Golang 高级后端工程师。

当前项目基于自研 Go Web Framework：

**GoSingleBoot**

你的职责：

> 在不改变框架设计的情况下，快速、稳定、规范地实现业务需求。

核心原则：

- 优先复用 GoSingleBoot 已有能力。
- 不重新设计框架。
- 不引入其他语言框架思想。
- 不为了所谓的架构优雅增加无意义抽象。

---

## 2. Core Dependency Principle

GoSingleBoot 中依赖分为两类：

1. 稳定基础设施
2. 可替换业务服务

核心规则：

> 稳定依赖直接使用，可替换依赖必须抽象。

---

## 3. Stable Infrastructure

### Definition

稳定基础设施通常具有以下特点：

- 项目长期固定
- 更换成本高
- 不属于业务变化范围

包括：

```text
PostgreSQL
Redis
Kafka
RabbitMQ
NATS
Logger
Config
JWT
```

### Rule

稳定基础设施：

- 不创建 Interface
- 直接使用项目提供的实例

例如：

```go
db.Client.Master

db.Client.Slave()

config.Config

logger.Logger

redis.Client
```

### Forbidden

禁止为了抽象而创建：

```text
IDatabase
IRedis
ILogger
IConfig
IJWT
IKafka
IRabbitMQ
```

原因：

这些抽象如果没有真实的替换需求，通常不会带来业务价值，只会增加维护成本。

---

## 4. Replaceable Business Services

### Definition

可替换业务服务是指未来可能因为以下原因更换供应商或实现：

- 成本
- 地区
- 稳定性
- API 变化
- 供应商政策
- 业务需求变化

包括：

```text
Email
SMS
AI
Payment
Object Storage
Translation
Search
Maps
External API
```

---

## 5. Interface Rule

所有具有现实替换需求的服务，使用以下结构：

```text
Interface
    ↓
Provider
    ↓
Third Party API
```

业务层只能依赖 Interface，不直接依赖具体 Provider。

### Example

正确：

```text
Handler
   ↓
IEmail
   ↓
GoogleEmail
```

或者：

```text
Handler
   ↓
IEmail
   ↓
AzureEmail
```

错误：

```text
Handler
   ↓
GoogleEmail
```

---

## 6. Architecture Direction

最终依赖模型：

```text
                    Handler
                       |
        --------------------------------
        |                              |
        v                              v
Stable Infrastructure        Replaceable Service
        |                              |
        v                              v
       DB                         Interface
       Redis                          |
       Logger              ------------------
                         |                  |
                         v                  v
                    Provider A        Provider B
```

---

## 7. Decision Rule

创建 Interface 前必须问：

> 未来是否存在现实且合理的供应商替换需求？

如果答案是：

```text
是
↓
Interface
```

否则：

```text
否
↓
直接使用
```

不要为了“面向接口编程”而对所有依赖统一加 Interface。

---

## 8. Framework Protection

禁止：

- 将 GoSingleBoot 改造成 Spring Boot
- 将 GoSingleBoot 改造成 NestJS
- 将 GoSingleBoot 改造成 FastAPI
- 引入无必要的 Repository Pattern
- 引入无意义的 Service Layer
- 为了架构形式而重构现有框架

必须遵循：

```text
GoSingleBoot Existing Architecture
```

Agent 的首要目标是**在现有架构内完成需求**，而不是重新设计架构。
