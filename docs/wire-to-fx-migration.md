# Wire 到 Uber FX 迁移指南

## 目录

1. [概述](#1-概述)
2. [核心概念对比](#2-核心概念对比)
3. [迁移步骤](#3-迁移步骤)
4. [代码示例](#4-代码示例)
5. [注意事项](#5-注意事项)
6. [迁移检查清单](#6-迁移检查清单)

---

## 1. 概述

### 1.1 当前项目架构

项目当前使用 **Google Wire** 进行依赖注入，主要涉及以下文件：

| 文件 | 作用 |
|------|------|
| `bootstrap/wire.go` | Wire 注入定义（编译时处理） |
| `bootstrap/wire_gen.go` | Wire 自动生成的注入代码 |
| `bootstrap/bootstrap.go` | Provider 函数实现 |
| `api/api.go` | API 结构体和 Provider Sets |
| `api/*/common.go` | 各模块的 Provider Sets |

### 1.2 Wire vs FX 对比

| 特性 | Google Wire | Uber FX |
|------|-------------|---------|
| 注入时机 | 编译时（代码生成） | 运行时（反射） |
| 性能 | 零运行时开销 | 极小的启动开销 |
| 生命周期管理 | 无 | 内置（OnStart/OnStop） |
| 热重载 | 不支持 | 支持 |
| 代码生成 | 需要 `go generate` | 不需要 |
| 错误发现 | 编译时 | 启动时 |
| 学习曲线 | 较陡 | 相对平缓 |

### 1.3 为什么选择 FX

- **生命周期管理**：内置优雅启动/关闭支持
- **无代码生成**：减少构建步骤，简化 CI/CD
- **更灵活的模块组织**：`fx.Module` 提供更好的封装
- **内置日志和调试**：便于排查依赖问题
- **运行时灵活性**：支持条件依赖和动态配置

---

## 2. 核心概念对比

### 2.1 Provider 定义

**Wire 方式：**
```go
// bootstrap/bootstrap.go
func UseGorm(v *common.Values) (*gorm.DB, error) {
    return gorm.Open(postgres.New(postgres.Config{
        DSN: v.Database.Url,
    }), &gorm.Config{})
}
```

**FX 方式：**
```go
// bootstrap/providers.go
func NewGorm(v *common.Values) (*gorm.DB, error) {
    return gorm.Open(postgres.New(postgres.Config{
        DSN: v.Database.Url,
    }), &gorm.Config{})
}

// 注册方式
fx.Provide(NewGorm)
```

### 2.2 Provider Sets / Modules

**Wire 方式：**
```go
// api/users/common.go
var Provides = wire.NewSet(
    wire.Struct(new(Controller), "*"),
    wire.Struct(new(Service), "*"),
)
```

**FX 方式：**
```go
// api/users/module.go
var Module = fx.Module("users",
    fx.Provide(
        NewService,
        NewController,
    ),
)

func NewService(inject *common.Inject, sessions *sessions.Service) *Service {
    return &Service{
        Inject:    inject,
        SessionsX: sessions,
    }
}

func NewController(v *common.Values, service *Service) *Controller {
    return &Controller{
        V:      v,
        UsersX: service,
    }
}
```

### 2.3 结构体自动组装

**Wire 方式：**
```go
wire.Struct(new(Controller), "*")  // 自动填充所有公开字段
```

**FX 方式：**
```go
// 方式 1: 显式构造函数（推荐）
func NewController(v *common.Values, service *Service) *Controller {
    return &Controller{
        V:      v,
        UsersX: service,
    }
}

// 方式 2: 使用 fx.In 参数结构体
type ControllerParams struct {
    fx.In
    V      *common.Values
    UsersX *Service
}

func NewController(p ControllerParams) *Controller {
    return &Controller{
        V:      p.V,
        UsersX: p.UsersX,
    }
}
```

### 2.4 依赖注入入口

**Wire 方式：**
```go
// bootstrap/wire.go
//go:build wireinject

func NewAPI(values *common.Values) (*api.API, error) {
    wire.Build(
        wire.Struct(new(api.API), "*"),
        wire.Struct(new(common.Inject), "*"),
        UseGorm,
        UseRedis,
        UsePassport,
        UseLocker,
        UseCaptcha,
        UseHertz,
        api.Provides,
    )
    return &api.API{}, nil
}
```

**FX 方式：**
```go
// main.go
func main() {
    fx.New(
        // 配置
        fx.Provide(func() (*common.Values, error) {
            return bootstrap.LoadStaticValues("./config/values.yml")
        }),

        // 基础设施
        bootstrap.Module,

        // API 模块
        api.Module,

        // 启动 HTTP 服务器
        fx.Invoke(startServer),
    ).Run()
}

func startServer(lc fx.Lifecycle, api *api.API) {
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            h, err := api.Initialize(ctx)
            if err != nil {
                return err
            }
            go h.Spin()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            // 优雅关闭
            return nil
        },
    })
}
```

---

## 3. 迁移步骤

### 步骤 1：添加 FX 依赖

```bash
go get go.uber.org/fx@latest
```

更新 `go.mod`：
```go
require (
    go.uber.org/fx v1.23.0  // 或最新版本
    // 移除 wire（迁移完成后）
    // github.com/google/wire v0.7.0
)
```

### 步骤 2：创建 Bootstrap Module

创建 `bootstrap/module.go`：

```go
package bootstrap

import (
    "context"
    "database/sql"
    "regexp"
    "time"

    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/common/config"
    "github.com/go-playground/validator/v10"
    "github.com/hertz-contrib/cors"
    go_playground "github.com/hertz-contrib/validator/go-playground"
    "github.com/kainonly/go/captcha"
    "github.com/kainonly/go/help"
    "github.com/kainonly/go/locker"
    "github.com/kainonly/go/passport"
    "github.com/redis/go-redis/v9"
    "go.uber.org/fx"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "server/common"
)

// Module 导出所有基础设施依赖
var Module = fx.Module("bootstrap",
    fx.Provide(
        NewGorm,
        NewRedis,
        NewPassport,
        NewLocker,
        NewCaptcha,
        NewHertz,
        NewInject,
    ),
)

// NewGorm 创建数据库连接
func NewGorm(lc fx.Lifecycle, v *common.Values) (*gorm.DB, error) {
    orm, err := gorm.Open(
        postgres.New(postgres.Config{
            DSN:                  v.Database.Url,
            PreferSimpleProtocol: true,
        }),
        &gorm.Config{
            SkipDefaultTransaction: true,
            PrepareStmt:            true,
        },
    )
    if err != nil {
        return nil, err
    }

    db, err := orm.DB()
    if err != nil {
        return nil, err
    }

    db.SetMaxIdleConns(10)
    db.SetMaxOpenConns(100)
    db.SetConnMaxLifetime(time.Hour)

    // 生命周期管理
    lc.Append(fx.Hook{
        OnStop: func(ctx context.Context) error {
            return db.Close()
        },
    })

    return orm, nil
}

// NewRedis 创建 Redis 连接
func NewRedis(lc fx.Lifecycle, v *common.Values) (*redis.Client, error) {
    opts, err := redis.ParseURL(v.Database.Redis)
    if err != nil {
        return nil, err
    }

    client := redis.NewClient(opts)
    if err = client.Ping(context.Background()).Err(); err != nil {
        return nil, err
    }

    // 生命周期管理
    lc.Append(fx.Hook{
        OnStop: func(ctx context.Context) error {
            return client.Close()
        },
    })

    return client, nil
}

// NewPassport 创建认证服务
func NewPassport(v *common.Values) *passport.Passport {
    return passport.New(
        passport.SetKey(v.Key),
        passport.SetIssuer(v.Domain),
    )
}

// NewLocker 创建分布式锁服务
func NewLocker(client *redis.Client) *locker.Locker {
    return locker.New(client)
}

// NewCaptcha 创建验证码服务
func NewCaptcha(client *redis.Client) *captcha.Captcha {
    return captcha.New(client)
}

// NewHertz 创建 HTTP 服务器
func NewHertz(v *common.Values) (*server.Hertz, error) {
    vd := go_playground.NewValidator()
    vd.SetValidateTag("vd")
    vdx := vd.Engine().(*validator.Validate)

    vdx.RegisterValidation("snake", func(fl validator.FieldLevel) bool {
        matched, err := regexp.MatchString("^[a-z_]+$", fl.Field().Interface().(string))
        return err == nil && matched
    })

    vdx.RegisterValidation("sort", func(fl validator.FieldLevel) bool {
        matched, err := regexp.MatchString("^[a-z_.]+:(-1|1)$", fl.Field().Interface().(string))
        return err == nil && matched
    })

    opts := []config.Option{
        server.WithHostPorts(v.Address),
        server.WithCustomValidator(vd),
    }

    h := server.Default(opts...)
    h.Use(
        help.ErrorHandler(),
        cors.New(cors.Config{
            AllowOrigins:     v.Cors.Origins,
            AllowMethods:     []string{"GET", "POST"},
            AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Page", "X-Pagesize", "X-Requested-With"},
            ExposeHeaders:    []string{"X-Total"},
            AllowCredentials: true,
            MaxAge:           12 * time.Hour,
        }),
    )

    return h, nil
}

// NewInject 创建通用注入结构体
func NewInject(
    v *common.Values,
    db *gorm.DB,
    rdb *redis.Client,
    captcha *captcha.Captcha,
    locker *locker.Locker,
) *common.Inject {
    return &common.Inject{
        V:       v,
        Db:      db,
        RDb:     rdb,
        Captcha: captcha,
        Locker:  locker,
    }
}
```

### 步骤 3：迁移各 API 模块

#### 3.1 Sessions 模块

创建 `api/sessions/module.go`：

```go
package sessions

import (
    "go.uber.org/fx"
    "server/common"
)

var Module = fx.Module("sessions",
    fx.Provide(
        NewService,
        NewController,
    ),
)

func NewService(inject *common.Inject) *Service {
    return &Service{
        Inject: inject,
    }
}

func NewController(v *common.Values, service *Service) *Controller {
    return &Controller{
        V:         v,
        SessionsX: service,
    }
}
```

#### 3.2 Index 模块

创建 `api/index/module.go`：

```go
package index

import (
    "go.uber.org/fx"
    "server/api/sessions"
    "server/common"
    "github.com/kainonly/go/passport"
)

var Module = fx.Module("index",
    fx.Provide(
        NewService,
        NewController,
    ),
)

func NewService(
    inject *common.Inject,
    passport *passport.Passport,
    sessionsX *sessions.Service,
) *Service {
    return &Service{
        Inject:    inject,
        Passport:  passport,
        SessionsX: sessionsX,
    }
}

func NewController(v *common.Values, service *Service) *Controller {
    return &Controller{
        V:      v,
        IndexX: service,
    }
}
```

#### 3.3 Users 模块

创建 `api/users/module.go`：

```go
package users

import (
    "go.uber.org/fx"
    "server/api/sessions"
    "server/common"
)

var Module = fx.Module("users",
    fx.Provide(
        NewService,
        NewController,
    ),
)

func NewService(inject *common.Inject, sessionsX *sessions.Service) *Service {
    return &Service{
        Inject:    inject,
        SessionsX: sessionsX,
    }
}

func NewController(v *common.Values, service *Service) *Controller {
    return &Controller{
        V:      v,
        UsersX: service,
    }
}
```

#### 3.4 其他模块（Jobs, Schedulers, Teams）

按照相同模式创建各自的 `module.go` 文件。

### 步骤 4：创建 API 主模块

修改 `api/api.go`，添加 FX Module：

```go
package api

import (
    "context"

    "github.com/cloudwego/hertz/pkg/app/server"
    "go.uber.org/fx"

    "server/api/index"
    "server/api/jobs"
    "server/api/schedulers"
    "server/api/sessions"
    "server/api/teams"
    "server/api/users"
    "server/common"
)

// Module 聚合所有 API 子模块
var Module = fx.Module("api",
    // 子模块
    sessions.Module,
    index.Module,
    users.Module,
    jobs.Module,
    schedulers.Module,
    teams.Module,

    // 主 API 结构体
    fx.Provide(NewAPI),
)

// APIParams 定义 API 依赖
type APIParams struct {
    fx.In

    Inject     *common.Inject
    Hertz      *server.Hertz
    Index      *index.Controller
    IndexX     *index.Service
    Jobs       *jobs.Controller
    Schedulers *schedulers.Controller
    Sessions   *sessions.Controller
    Teams      *teams.Controller
    Users      *users.Controller
    UsersX     *users.Service
}

// NewAPI 创建 API 实例
func NewAPI(p APIParams) *API {
    return &API{
        Inject:     p.Inject,
        Hertz:      p.Hertz,
        Index:      p.Index,
        IndexX:     p.IndexX,
        Jobs:       p.Jobs,
        Schedulers: p.Schedulers,
        Sessions:   p.Sessions,
        Teams:      p.Teams,
        Users:      p.Users,
        UsersX:     p.UsersX,
    }
}

type API struct {
    *common.Inject
    Hertz      *server.Hertz
    Index      *index.Controller
    IndexX     *index.Service
    Jobs       *jobs.Controller
    Schedulers *schedulers.Controller
    Sessions   *sessions.Controller
    Teams      *teams.Controller
    Users      *users.Controller
    UsersX     *users.Service
}

func (x *API) Initialize(ctx context.Context) (*server.Hertz, error) {
    authx := x.Auth()

    // 公开路由
    x.Hertz.GET("", x.Index.Ping)
    x.Hertz.POST("login", x.Index.Login)
    x.Hertz.GET("verify", x.Index.Verify)
    x.Hertz.POST("logout", authx, x.Index.Logout)

    // 受保护的路由
    r := x.Hertz.Group("", authx)
    r.GET("user", x.Index.GetUser)
    r.POST("user/set_password", x.Index.SetUserPassword)

    // Users CRUD
    r.GET("users", x.Users.Find)
    r.GET("users/:id", x.Users.FindById)
    r.POST("users/create", x.Users.Create)
    r.POST("users/update", x.Users.Update)
    r.POST("users/delete", x.Users.Delete)
    r.GET("users/_exists", x.Users.Exists)
    r.GET("users/_search", x.Users.Search)
    r.POST("users/set_statuses", x.Users.SetStatuses)

    // 其他路由...

    return x.Hertz, nil
}
```

### 步骤 5：重写 main.go

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/cloudwego/hertz/pkg/app/server"
    "go.uber.org/fx"
    "go.uber.org/fx/fxevent"

    "server/api"
    "server/bootstrap"
    "server/common"
)

func main() {
    app := fx.New(
        // 可选：配置日志
        fx.WithLogger(func() fxevent.Logger {
            return &fxevent.ConsoleLogger{W: os.Stdout}
        }),

        // 提供配置
        fx.Provide(provideValues),

        // 基础设施模块
        bootstrap.Module,

        // API 模块
        api.Module,

        // 启动服务器
        fx.Invoke(registerHooks),
    )

    app.Run()
}

// provideValues 加载配置
func provideValues() (*common.Values, error) {
    return bootstrap.LoadStaticValues("./config/values.yml")
}

// registerHooks 注册生命周期钩子
func registerHooks(lc fx.Lifecycle, apiInstance *api.API) {
    var h *server.Hertz

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            var err error
            h, err = apiInstance.Initialize(ctx)
            if err != nil {
                return fmt.Errorf("failed to initialize API: %w", err)
            }

            // 在 goroutine 中启动服务器（非阻塞）
            go func() {
                h.Spin()
            }()

            return nil
        },
        OnStop: func(ctx context.Context) error {
            if h != nil {
                return h.Shutdown(ctx)
            }
            return nil
        },
    })
}
```

### 步骤 6：清理 Wire 相关文件

迁移完成并测试通过后：

1. 删除 `bootstrap/wire.go`
2. 删除 `bootstrap/wire_gen.go`
3. 删除各模块 `common.go` 中的 `wire.NewSet` 定义
4. 从 `go.mod` 移除 Wire 依赖：
   ```bash
   go mod tidy
   ```

---

## 4. 代码示例

### 4.1 完整的模块结构示例

```
api/
├── api.go              # API 主结构 + fx.Module
├── auth.go             # 认证中间件（无变化）
├── index/
│   ├── module.go       # 新增：fx.Module 定义
│   ├── common.go       # 移除 wire.NewSet，保留结构体定义
│   ├── login.go        # 无变化
│   └── ...
├── users/
│   ├── module.go       # 新增：fx.Module 定义
│   ├── common.go       # 移除 wire.NewSet，保留结构体定义
│   ├── find.go         # 无变化
│   └── ...
└── ...

bootstrap/
├── module.go           # 新增：fx.Module 定义
├── bootstrap.go        # 重命名函数（UseXxx → NewXxx）
└── (删除 wire.go 和 wire_gen.go)
```

### 4.2 使用 fx.In 处理多依赖

当构造函数参数过多时，使用 `fx.In`：

```go
type ServiceParams struct {
    fx.In

    Inject    *common.Inject
    Passport  *passport.Passport
    SessionsX *sessions.Service
    UsersX    *users.Service
    // 可以继续添加更多依赖
}

func NewService(p ServiceParams) *Service {
    return &Service{
        Inject:    p.Inject,
        Passport:  p.Passport,
        SessionsX: p.SessionsX,
        UsersX:    p.UsersX,
    }
}
```

### 4.3 使用 fx.Out 提供多个值

```go
type BootstrapResult struct {
    fx.Out

    DB      *gorm.DB
    Redis   *redis.Client
    Locker  *locker.Locker
    Captcha *captcha.Captcha
}

func NewInfrastructure(v *common.Values) (BootstrapResult, error) {
    // 初始化所有基础设施...
    return BootstrapResult{
        DB:      db,
        Redis:   rdb,
        Locker:  lock,
        Captcha: cap,
    }, nil
}
```

### 4.4 条件依赖（FX 独有优势）

```go
// 根据配置决定是否启用某功能
func NewOptionalFeature(v *common.Values) (*Feature, error) {
    if !v.FeatureEnabled {
        return nil, nil  // FX 会忽略 nil 返回
    }
    return &Feature{}, nil
}

// 或使用 fx.Supply 提供可选值
fx.Provide(fx.Annotate(
    NewOptionalFeature,
    fx.ResultTags(`optional:"true"`),
))
```

### 4.5 命名依赖

当同一类型有多个实例时：

```go
// 提供方
fx.Provide(
    fx.Annotate(
        NewPrimaryDB,
        fx.ResultTags(`name:"primary"`),
    ),
    fx.Annotate(
        NewReplicaDB,
        fx.ResultTags(`name:"replica"`),
    ),
)

// 使用方
type Params struct {
    fx.In

    PrimaryDB *gorm.DB `name:"primary"`
    ReplicaDB *gorm.DB `name:"replica"`
}
```

---

## 5. 注意事项

### 5.1 依赖循环问题

FX 会在启动时检测循环依赖并报错。如果遇到此问题：

1. **重新设计依赖关系**：拆分服务，提取公共部分
2. **使用接口解耦**：通过接口打破直接依赖
3. **延迟初始化**：使用 `fx.Invoke` 延迟某些操作

### 5.2 启动顺序

FX 会自动处理依赖顺序。如需显式控制：

```go
lc.Append(fx.Hook{
    OnStart: fx.Hook{
        OnStart: func(ctx context.Context) error {
            // 这里的代码会在所有依赖准备好后执行
            return nil
        },
    },
})
```

### 5.3 测试支持

FX 提供 `fxtest` 包简化测试：

```go
func TestService(t *testing.T) {
    app := fxtest.New(t,
        fx.Provide(NewMockDB),
        fx.Provide(NewService),
        fx.Invoke(func(s *Service) {
            // 测试逻辑
        }),
    )
    app.RequireStart()
    defer app.RequireStop()
}
```

### 5.4 调试依赖图

FX 提供内置的依赖可视化：

```go
// 打印依赖图
fx.New(
    // ... modules
    fx.Options(fx.NopLogger),  // 禁用默认日志
    fx.Invoke(func(g fx.DotGraph) {
        fmt.Println(g)  // 输出 Graphviz DOT 格式
    }),
)
```

### 5.5 性能考虑

- Wire 在编译时生成代码，运行时无反射开销
- FX 使用反射，但只在启动时执行一次
- 对于大多数应用，FX 的启动开销可以忽略不计（通常 < 100ms）

---

## 6. 迁移检查清单

### 准备阶段

- [ ] 阅读 [FX 官方文档](https://uber-go.github.io/fx/)
- [ ] 添加 FX 依赖：`go get go.uber.org/fx`
- [ ] 理解项目现有的 Wire 依赖图

### 迁移阶段

- [ ] 创建 `bootstrap/module.go`
  - [ ] 迁移 `UseGorm` → `NewGorm`
  - [ ] 迁移 `UseRedis` → `NewRedis`
  - [ ] 迁移 `UsePassport` → `NewPassport`
  - [ ] 迁移 `UseLocker` → `NewLocker`
  - [ ] 迁移 `UseCaptcha` → `NewCaptcha`
  - [ ] 迁移 `UseHertz` → `NewHertz`
  - [ ] 创建 `NewInject` 构造函数
  - [ ] 添加生命周期钩子（数据库/Redis 连接关闭）

- [ ] 迁移 API 模块
  - [ ] `api/sessions/module.go`
  - [ ] `api/index/module.go`
  - [ ] `api/users/module.go`
  - [ ] `api/jobs/module.go`
  - [ ] `api/schedulers/module.go`
  - [ ] `api/teams/module.go`

- [ ] 创建 `api/api.go` 中的 fx.Module
- [ ] 重写 `main.go`

### 测试阶段

- [ ] 运行应用，检查启动日志
- [ ] 测试所有 API 端点
- [ ] 测试优雅关闭（Ctrl+C）
- [ ] 检查数据库/Redis 连接是否正确关闭

### 清理阶段

- [ ] 删除 `bootstrap/wire.go`
- [ ] 删除 `bootstrap/wire_gen.go`
- [ ] 清理 `common.go` 中的 `wire.NewSet`
- [ ] 运行 `go mod tidy`
- [ ] 更新项目文档

---

## 附录：关键文件变更摘要

| 原文件 | 操作 | 新文件/变更 |
|--------|------|-------------|
| `bootstrap/wire.go` | 删除 | - |
| `bootstrap/wire_gen.go` | 删除 | - |
| `bootstrap/bootstrap.go` | 修改 | 函数重命名，添加 `NewInject` |
| - | 新增 | `bootstrap/module.go` |
| `api/api.go` | 修改 | 添加 `fx.Module`，移除 `wire.NewSet` |
| `api/*/common.go` | 修改 | 移除 `wire.NewSet`，添加构造函数 |
| - | 新增 | `api/*/module.go` |
| `main.go` | 重写 | 使用 `fx.New().Run()` |
| `go.mod` | 修改 | 添加 `go.uber.org/fx`，移除 Wire |

---

## 参考资料

- [Uber FX 官方文档](https://uber-go.github.io/fx/)
- [FX GitHub 仓库](https://github.com/uber-go/fx)
- [Google Wire 文档](https://github.com/google/wire)
- [FX vs Wire 对比文章](https://blog.drewolson.org/dependency-injection-in-go)
