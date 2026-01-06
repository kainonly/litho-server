# Uber FX 依赖注入指南

## 目录

1. [概述](#1-概述)
2. [核心概念](#2-核心概念)
3. [项目结构](#3-项目结构)
4. [代码示例](#4-代码示例)
5. [高级用法](#5-高级用法)
6. [注意事项](#6-注意事项)

---

## 1. 概述

### 1.1 为什么选择 FX

- **生命周期管理**：内置优雅启动/关闭支持（OnStart/OnStop）
- **无代码生成**：减少构建步骤，简化 CI/CD
- **内置日志和调试**：便于排查依赖问题
- **运行时灵活性**：支持条件依赖和动态配置

> **注**：Google Wire 项目已于 2025 年 8 月 25 日[归档](https://github.com/google/wire)，由于 `golang.org/x/tools` 依赖版本过旧，在 Go 1.25 中[无法正常使用](https://github.com/google/wire/issues/431)。

---

## 2. 核心概念

### 2.1 Provider

Provider 是创建依赖的函数，通过 `fx.Provide` 注册：

```go
func NewGorm(v *common.Values) (*gorm.DB, error) {
    return gorm.Open(postgres.New(postgres.Config{
        DSN: v.Database.Url,
    }), &gorm.Config{})
}

// 注册
fx.Provide(NewGorm)
```

### 2.2 结构体注入 (fx.In)

使用 `fx.In` 结构体收集依赖，避免冗长的参数列表：

```go
type ServiceParams struct {
    fx.In

    Inject    *common.Inject
    Passport  *passport.Passport
    SessionsX *sessions.Service
}

func NewService(p ServiceParams) *Service {
    return &Service{
        Inject:    p.Inject,
        Passport:  p.Passport,
        SessionsX: p.SessionsX,
    }
}
```

### 2.3 生命周期钩子

使用 `fx.Lifecycle` 管理资源的启动和关闭：

```go
func NewGorm(lc fx.Lifecycle, v *common.Values) (*gorm.DB, error) {
    orm, err := gorm.Open(...)
    if err != nil {
        return nil, err
    }

    db, _ := orm.DB()

    lc.Append(fx.Hook{
        OnStop: func(ctx context.Context) error {
            return db.Close()
        },
    })

    return orm, nil
}
```

---

## 3. 项目结构

```
bootstrap/
└── bootstrap.go        # Provider 函数实现

api/
├── api.go              # API 主结构
├── auth.go             # 认证中间件
└── index/
    ├── common.go       # Service/Controller 定义
    └── controller.go   # 业务逻辑
```

---

## 4. 代码示例

### 4.1 Bootstrap Providers

```go
package bootstrap

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
    "go.uber.org/fx"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "server/common"
)

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

    lc.Append(fx.Hook{
        OnStop: func(ctx context.Context) error {
            return db.Close()
        },
    })

    return orm, nil
}

func NewRedis(lc fx.Lifecycle, v *common.Values) (*redis.Client, error) {
    opts, err := redis.ParseURL(v.Database.Redis)
    if err != nil {
        return nil, err
    }

    client := redis.NewClient(opts)
    if err = client.Ping(context.Background()).Err(); err != nil {
        return nil, err
    }

    lc.Append(fx.Hook{
        OnStop: func(ctx context.Context) error {
            return client.Close()
        },
    })

    return client, nil
}

type InjectParams struct {
    fx.In

    V       *common.Values
    Db      *gorm.DB
    RDb     *redis.Client
    Captcha *captcha.Captcha
    Locker  *locker.Locker
}

func NewInject(p InjectParams) *common.Inject {
    return &common.Inject{
        V:       p.V,
        Db:      p.Db,
        RDb:     p.RDb,
        Captcha: p.Captcha,
        Locker:  p.Locker,
    }
}
```

### 4.2 API Providers (以 Index 为例)

```go
package index

import (
    "go.uber.org/fx"

    "server/api/sessions"
    "server/common"
    "github.com/kainonly/go/passport"
)

type ServiceParams struct {
    fx.In

    Inject    *common.Inject
    Passport  *passport.Passport
    SessionsX *sessions.Service
}

func NewService(p ServiceParams) *Service {
    return &Service{
        Inject:    p.Inject,
        Passport:  p.Passport,
        SessionsX: p.SessionsX,
    }
}

type ControllerParams struct {
    fx.In

    V      *common.Values
    IndexX *Service
}

func NewController(p ControllerParams) *Controller {
    return &Controller{
        V:      p.V,
        IndexX: p.IndexX,
    }
}
```

### 4.3 API 主结构

```go
package api

import (
    "go.uber.org/fx"

    "server/api/index"
    "server/common"
)

type APIParams struct {
    fx.In

    Inject *common.Inject
    Hertz  *server.Hertz
    Index  *index.Controller
    IndexX *index.Service
}

func NewAPI(p APIParams) *API {
    return &API{
        Inject: p.Inject,
        Hertz:  p.Hertz,
        Index:  p.Index,
        IndexX: p.IndexX,
    }
}
```

### 4.4 main.go

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
    "server/api/index"
    "server/bootstrap"
    "server/common"
)

func main() {
    fx.New(
        fx.WithLogger(func() fxevent.Logger {
            return &fxevent.ConsoleLogger{W: os.Stdout}
        }),

        fx.Provide(
            // 配置
            func() (*common.Values, error) {
                return bootstrap.LoadStaticValues("./config/values.yml")
            },

            // 基础设施
            bootstrap.NewGorm,
            bootstrap.NewRedis,
            bootstrap.NewPassport,
            bootstrap.NewLocker,
            bootstrap.NewCaptcha,
            bootstrap.NewHertz,
            bootstrap.NewInject,

            // API
            index.NewService,
            index.NewController,
            api.NewAPI,
        ),

        fx.Invoke(registerHooks),
    ).Run()
}

func registerHooks(lc fx.Lifecycle, apiInstance *api.API) {
    var h *server.Hertz

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            var err error
            h, err = apiInstance.Initialize(ctx)
            if err != nil {
                return fmt.Errorf("failed to initialize API: %w", err)
            }
            go h.Spin()
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

---

## 5. 高级用法

### 5.1 使用 fx.Out 提供多个值

```go
type BootstrapResult struct {
    fx.Out

    DB      *gorm.DB
    Redis   *redis.Client
    Locker  *locker.Locker
}

func NewInfrastructure(v *common.Values) (BootstrapResult, error) {
    // ...
    return BootstrapResult{
        DB:     db,
        Redis:  rdb,
        Locker: lock,
    }, nil
}
```

### 5.2 命名依赖

当同一类型有多个实例时：

```go
// 提供方
fx.Provide(
    fx.Annotate(NewPrimaryDB, fx.ResultTags(`name:"primary"`)),
    fx.Annotate(NewReplicaDB, fx.ResultTags(`name:"replica"`)),
)

// 使用方
type Params struct {
    fx.In

    PrimaryDB *gorm.DB `name:"primary"`
    ReplicaDB *gorm.DB `name:"replica"`
}
```

### 5.3 可选依赖

```go
type Params struct {
    fx.In

    Cache *redis.Client `optional:"true"`
}
```

### 5.4 调试依赖图

```go
fx.New(
    // ... providers
    fx.Invoke(func(g fx.DotGraph) {
        fmt.Println(g)  // 输出 Graphviz DOT 格式
    }),
)
```

---

## 6. 注意事项

### 6.1 依赖循环

FX 会在启动时检测循环依赖并报错。解决方法：

1. 重新设计依赖关系，拆分服务
2. 使用接口解耦
3. 使用 `fx.Invoke` 延迟初始化

### 6.2 测试

使用 `fxtest` 包简化测试：

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

### 6.3 性能

FX 使用反射，但只在启动时执行一次，对于大多数应用启动开销可忽略（通常 < 100ms）。

---

## 参考资料

- [Uber FX 官方文档](https://uber-go.github.io/fx/)
- [FX GitHub 仓库](https://github.com/uber-go/fx)
