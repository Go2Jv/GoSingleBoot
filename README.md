# GoSingleBoot

一个基于 **Go + Gin + Bun** 的单体后端快速开发脚手架。开箱即用地提供配置管理、日志、数据库主从读写分离、JWT 认证、统一错误处理、CORS、Panic 恢复、优雅关闭、API 文档自动生成等能力,让开发者只需专注于编写业务代码。

> 🔨 **goBootCLI 脚手架工具**:本项目配套了官方的交互式 CLI 工具 [goBootCLI](./gobootCLI),一条命令即可从本仓库模板创建并配置好新项目,无需手动改 `go.mod`、import 路径和 `config.json`。详见 [goBootCLI 使用指南](#gobootcli-脚手架工具)。

> 📌 **数据库支持说明**:ORM 层选用 [Bun](https://bun.uptrace.dev/),当前仅集成 **PostgreSQL**(`pgdialect` + `pgdriver`),**暂不支持 MySQL**。

> ❤️ **长期维护承诺**:本项目持续维护中,欢迎通过 [Issues](https://github.com/Go2Jv/GoSingleBoot/issues) 反馈问题、提出建议或贡献代码。

## 解决什么问题

在从零搭建 Go 后端服务时,每次都要重复处理大量"基础设施"代码:

| 痛点 | 本项目解决方案 |
| --- | --- |
| 每个 handler 都要手写错误判断、打日志、返回 JSON,代码冗长且容易漏打日志或重复打日志 | 统一的 `BizErr` 错误体系 + 全局错误中间件:**业务代码只管抛错,日志与响应由中间件统一处理** |
| 日志内容和返回给用户的内容混在一起,容易把内部细节(如 SQL 语句、DSN)泄露给客户端 | `BizErr` 同时携带 `Msg`(给用户看)与 `Log`(记日志用),**日志信息与用户展示信息分离** |
| 重复配置 CORS / Panic 恢复 / 404 / JWT / 配置读取 | 全部预置为中间件和工具包,通过 `config.json` 一键配置 |
| 进程被 `kill` 时无法优雅处理在途请求 | 监听 `SIGINT` / `SIGTERM`,调用 `server.Shutdown` 优雅关闭(15s 超时) |
| 接口文档与代码脱节 | 开发环境启动时自动执行 `respec` 扫描代码生成 OpenAPI 文档,并通过 Scalar 提供可视化页面 |
| 数据库读写压力集中 | 内置 Master/Slave 主从连接管理:`Slave()` 无从库时回退主库、单从库直用、多从库随机负载均衡 |
| 配置项与代码耦合、启动时才发现配置错误 | 集中式 `config.json`,启动阶段即完成 CORS 配置合法性校验(如 `AllowCredentials=true` 时禁止空 Origin 或 `*`) |

## 技术栈

| 类别 | 选型 |
| --- | --- |
| 语言 | Go 1.26 |
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [Bun](https://bun.uptrace.dev/)(PostgreSQL 专用:`pgdialect` / `pgdriver`) |
| 认证 | [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)(HS256) |
| 日志 | [Uber Zap](https://github.com/uber-go/zap)(生产/开发双模式) |
| CORS | [gin-contrib/cors](https://github.com/gin-contrib/cors) |
| API 文档 | [go-respec](https://github.com/Zachacious/go-respec)(开发环境自动生成 OpenAPI)+ [Scalar](https://github.com/MarceloPetrucio/go-scalar-api-reference)(可视化页面) |
| 脚手架工具 | [goBootCLI](./gobootCLI)(纯 Go 标准库,零第三方依赖) |
| 许可证 | MIT |

## 核心功能

### 1. 集中式配置(`config.json`)

- **Application**:应用名、监听端口、运行模式(`Text` 字段控制生产/开发模式:开发模式开启 Gin Debug 日志与 API 文档生成)
- **Database**:`Master` 主库 DSN + `Slaves` 从库 DSN 列表(支持主从读写分离)
- **Jwt**:密钥、过期时间(小时)、签发者
- **Cors**:允许的来源列表、是否允许携带凭证
- 启动时自动校验 CORS 配置合法性,非法配置直接 `panic` 快速失败

### 2. 业务错误体系(`bizErr`)+ 全局错误中间件

业务代码中无需手写错误响应与日志,只需一行抛出:

```go
// 通用业务错误:自动包装并交给全局错误中间件处理
bizErr.Throw(c, 400, "账号密码错误", nil)

// 参数绑定校验:一行替代原来的一大堆 if err != nil
if ok := bizErr.Validation(c, err); !ok {
    return
}

// 单行查询未命中:自动识别 sql.ErrNoRows 返回 404
if ok := bizErr.SQLNotFound(c, err); !ok {
    return
}
```

全局错误中间件统一兜底:`BizErr` → 返回 `code`/`msg` 给用户,`Log` 写入日志;**未知错误** → 返回统一的"服务器繁忙"提示,不泄露内部细节。

### 3. Panic 恢复中间件

任何 handler 中的 `panic` 都会被捕获,记录 Error 级日志(含堆栈)并返回 500 友好响应,进程不崩溃。

### 4. JWT 封装

`GenerateToken` / `ParseToken` / `ValidateToken` 三个方法完成 HS256 签名令牌的生成、解析与校验,并强制校验签名算法,防止算法混淆攻击。

### 5. 优雅关闭

监听 `SIGINT` / `SIGTERM` 信号,收到后停止接收新请求,15 秒内等待在途请求处理完毕再退出。

### 6. API 文档自动生成(仅开发模式)

开发模式下启动时自动执行 `respec` 扫描代码注释生成 `openapi.yaml`,访问 `/docs` 即可在浏览器中查看 Scalar 渲染的交互式 API 文档。

### 7. 数据库主从读写分离

- `db.Client.Master`:写操作
- `db.Client.Slave()`:读操作,自动选择——无从库 → 主库;单个从库 → 直接使用;多个从库 → 随机负载均衡
- 连接失败在启动阶段即 `panic`,快速暴露配置问题

## goBootCLI 脚手架工具

`gobootCLI` 是 GoSingleBoot 官方配套的交互式脚手架工具:它把本仓库当作项目模板,**一条命令完成** clone 模板 → 替换 module 与 import 路径 → 更新配置 → `go mod tidy` → 清理模板文件 → 可选初始化 Git,生成的项目即可直接编译运行。工具本身仅使用 Go 标准库,零第三方依赖。

### 下载安装(Release 二进制)

前往 [Releases](https://github.com/Go2Jv/GoSingleBoot/releases) 页面下载对应平台的 `goboot` 二进制(当前版本:`v1.0.1`):

| 平台 | 文件 |
| --- | --- |
| Windows x64 / x86 | `goboot-windows-amd64.exe` / `goboot-windows-386.exe` |
| macOS Universal(Apple Silicon + Intel 通用) | `goboot-darwin-universal` |

macOS 首次运行若提示无法验证开发者,可执行 `xattr -d com.apple.quarantine goboot-darwin-universal` 解除隔离(Windows 下载 `.exe` 后如被 SmartScreen 拦截,点击"仍要运行"即可,均为未签名的开源二进制)。

也可以随时从源码运行:

```bash
go run ./gobootCLI
```

### 使用流程

运行后按提示依次输入:

```text
Project name: my-api                                 ← 必填,不能为空/含路径分隔符/以 . 开头
Project location (default: current directory):       ← 新项目创建位置,支持 ~ 展开,回车为当前目录
Cloning template into my-api...                      ← git clone --depth 1 本仓库模板
Removed template files: .git, README.md, LICENSE, gobootCLI   ← 立即清理模板残留
Updated Application.Name in my-api/config.json       ← 自动把应用名改为项目名
Go module path (default: my-api):                    ← 回车用项目名,或输入自定义 module path
Setting go.mod module to my-api...                   ← go mod edit -module
Updated imports in ...                               ← 仅 .go 文件 import 中的旧路径 → 新路径
Running go mod tidy in my-api...                     ← 在新项目目录执行,输出直通终端
Initialize Git repository? (Y/n):                    ← 默认 Yes;n/no 则跳过 git init

Project created successfully!
...
```

### 关键行为

- **安全性**:目标目录已存在则直接报错退出,绝不覆盖;clone 失败自动清理半成品目录
- **精确替换**:只改写 `go.mod` 的 module 行与 `.go` 文件 import 中 `"GoSingleBoot/` 开头的路径,不做全项目无脑字符串替换;跳过 `.git`
- **模板清理**:clone 完成后立即删除 `.git`、`README.md`、`LICENSE` 以及 `gobootCLI` 自身(避免新项目残留嵌套 module 导致 `go mod tidy` 失败)
- **配置联动**:自动把新项目 `config.json` 中的 `Application.Name` 更新为项目名(仅最小化修改,保留原格式与注释)
- **环境要求**:本机需安装 `git` 与 `go`(clone 与 `go mod edit` / `go mod tidy` 依赖)

### 常用命令

```bash
# 源码运行
go run ./gobootCLI

# 编译到当前目录
go build -o goboot ./gobootCLI
./goboot
```

## 项目结构

```
.
├── cmd/main.go                        # 程序入口:初始化配置/日志/数据库,启动 HTTP 服务并优雅关闭
├── config.json                        # 集中式配置文件
├── gobootCLI/                         # 脚手架工具(独立 module,纯 Go 标准库)
│   ├── main.go                        # CLI 入口
│   └── cli/cli.go                     # 交互式创建项目全流程
├── openapi.yaml                       # 自动生成的 OpenAPI 文档(开发模式)
├── internal/
│   ├── bizErr/                        # 业务错误体系:Wrap/Throw/Validation/SQLNotFound
│   ├── config/                        # 配置读取与校验
│   ├── db/                            # Bun + PostgreSQL 客户端(主从管理)
│   ├── docs/                          # OpenAPI 自动生成 + Scalar 文档页面
│   ├── dto/
│   │   ├── req/                       # 请求体结构(含 binding 校验)
│   │   └── resp/                      # 统一响应结构 CodeMsgAndData / CodeAndMsg
│   ├── handler/                       # 业务处理器(LoginHandler 演示)
│   ├── Interfaces/                    # 接口定义(如 IEmail)
│   ├── jwt/                           # JWT 生成/解析/校验封装
│   ├── logger/                        # Zap 日志初始化(生产/开发双模式)
│   ├── middleware/                    # CORS / Panic 恢复 / 全局错误处理
│   ├── model/                         # Bun ORM 模型(User 演示)
│   ├── router/                        # 主路由(注册中间件)与业务路由分组
│   └── Test/                          # 单元测试
```

## 版本与维护

- **最新 Release**:[v1.0.1](https://github.com/Go2Jv/GoSingleBoot/releases)(提供 Windows x64/x86 与 macOS Universal 的 goBootCLI 二进制)
- **Go 版本要求**:
  - goBootCLI 工具本身:**Go 1.22 及以上**(纯标准库)
  - 由模板生成的项目:当前模板依赖(gin、cors、x/net 等)要求 **Go 1.25 及以上**;CLI 会自动把新项目 go.mod 的 go 版本号同步为本机版本,低于依赖下限时 `go mod tidy` 会给出依赖自身的明确提示
- **维护承诺**:本项目长期维护(LTS),持续修复问题、更新依赖,欢迎通过 [Issues](https://github.com/Go2Jv/GoSingleBoot/issues) 反馈

## 快速开始

### 环境要求

- Go 1.26+
- PostgreSQL(已存在 `users` 表:含 `id`、`username`、`password` 字段)

### 步骤

1. 修改 `config.json` 中的数据库 DSN 与端口等配置
2. 运行服务:

```bash
go run cmd/main.go
```

3. 开发模式下访问 API 文档:<http://localhost:8080/docs>

### 示例接口:登录

```
POST /api/login
Content-Type: application/json

{
  "username": "admin",
  "password": "123456"
}
```

成功返回统一响应结构:

```json
{
  "code": 200,
  "msg": "登陆成功",
  "data": "<JWT Token>"
}
```

## 约定与设计说明

- **错误处理分工**:handler 只负责"抛出"错误(`bizErr.Throw` / `Validation` / `SQLNotFound`),日志打印与响应渲染统一由 `GlobalErrorMiddleware` 完成,避免重复打日志
- **日志与展示分离**:`BizErr.Msg` 面向用户,`BizErr.Log` 面向日志,两者内容可以不同
- **统一响应结构**:所有接口返回 `code` + `msg`(+ 可选 `data`)的 JSON 结构
- **404 兜底**:未注册路由统一返回 `{"code": 404, "msg": "不存在该Api"}`

## 许可证

[MIT](./LICENSE)
