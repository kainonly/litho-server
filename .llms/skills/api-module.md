# API 模块规范

> 当用户要求创建 API 接口、新增模块或修改路由时，使用此规范生成符合项目标准的代码。

## 触发条件

- 用户要求"创建 API"、"新增接口"、"添加路由"
- 用户要求为某个模型创建 CRUD 接口
- 用户要求修改现有 API 模块

## 目录结构

```
api/
├── api.go              # 路由聚合与生命周期管理
├── auth.go             # 认证中间件
├── index/              # 根路由模块
│   ├── common.go       # Provides、Controller、Service
│   └── ping.go         # 具体 API 方法
├── users/              # /users 路由模块
│   ├── common.go
│   ├── find.go
│   ├── find_by_id.go
│   ├── create.go
│   ├── update.go
│   └── delete.go
└── {module}/           # 其他模块
```

## common.go 模板

每个模块必须包含 `common.go`：

```go
package {module}

import (
	"server/common"

	"github.com/goforj/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	{Module}X *Service
}

type Service struct {
	*common.Inject
}
```

### 带额外依赖的模块

```go
package users

import (
	"server/api/sessions"
	"server/common"

	"github.com/goforj/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	UsersX *Service
}

type Service struct {
	*common.Inject

	SessionsX *sessions.Service
}
```

## 命名规则

| 元素 | 格式 | 示例 |
|------|------|------|
| 包名 | 小写模块名 | `users`、`orgs` |
| Controller 服务字段 | `{Module}X` | `UsersX`、`OrgsX` |
| Service 嵌入字段 | `*common.Inject` | 固定 |

## 标准路由约定

| 路由 | Method | 方法 | 说明 |
|------|--------|------|------|
| `/{module}/:id` | GET | `FindById` | 根据 ID 获取 |
| `/{module}` | GET | `Find` | 获取列表 |
| `/{module}/create` | POST | `Create` | 创建 |
| `/{module}/bulk_create` | POST | `BulkCreate` | 批量创建（可选） |
| `/{module}/update` | POST | `Update` | 更新 |
| `/{module}/delete` | POST | `Delete` | 删除 |

### 可选路由

| 路由 | Method | 方法 | 说明 |
|------|--------|------|------|
| `/{module}/_search` | GET | `Search` | 异步搜索 |
| `/{module}/_exists` | GET | `Exists` | 异步验证 |

### 操作类路由

| 路由格式 | Method | 方法格式 |
|----------|--------|----------|
| `/{module}/set_{noun}` | POST | `Set{Noun}` |
| `/{module}/{verb}` | POST | `{Verb}` |

示例：`/users/set_active` → `SetActive`，`/users/sort` → `Sort`

### 数据获取路由

| 路由格式 | Method | 方法格式 |
|----------|--------|----------|
| `/{module}/fetch_{data}` | POST | `Fetch{Data}` |

示例：`/users/fetch_statistics` → `FetchStatistics`

## api.go 结构

```go
package api

import (
	"context"
	"server/api/index"
	"server/api/orgs"
	"server/api/users"
	"server/common"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/goforj/wire"
	"github.com/kainonly/go/csrf"
)

var Provides = wire.NewSet(
	index.Provides,
	orgs.Provides,
	users.Provides,
)

type API struct {
	*common.Inject

	Hertz  *server.Hertz
	Csrf   *csrf.Csrf
	Index  *index.Controller
	IndexX *index.Service
	Orgs   *orgs.Controller
	Users  *users.Controller
	UsersX *users.Service
}

func (x *API) Initialize(ctx context.Context) (_ *server.Hertz, err error) {
	_csrf := x.Csrf.VerifyToken(!x.V.IsRelease())
	_auth := x.Auth()

	// 无需认证的路由
	x.Hertz.GET("", x.Index.Ping)
	x.Hertz.POST("login", _csrf, x.Index.Login)

	// 需要认证的路由组
	m := x.Hertz.Group(``, _csrf, _auth)
	{
		// orgs 模块
		m.GET("/orgs/:id", x.Orgs.FindById)
		m.GET("/orgs", x.Orgs.Find)
		m.POST("/orgs/create", x.Orgs.Create)
		m.POST("/orgs/update", x.Orgs.Update)
		m.POST("/orgs/delete", x.Orgs.Delete)
	}

	return x.Hertz, nil
}
```

### 注册新模块步骤

1. 在 `import` 中添加模块路径
2. 在 `Provides` 的 `wire.NewSet()` 中添加 `{module}.Provides`
3. 在 `API` 结构体中添加 `{Module} *{module}.Controller`
4. 如需在 Auth 使用，添加 `{Module}X *{module}.Service`
5. 在 `Initialize` 中注册路由

## API 方法示例

### find.go

```go
package orgs

import (
	"context"
	"server/common"
	"server/model"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

type FindDto struct {
	common.FindDto
}

func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto FindDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	total, data, err := x.OrgsX.Find(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-total", strconv.Itoa(int(total)))
	c.JSON(200, data)
}

type FindResult struct {
	ID     string `json:"id"`
	Active bool   `json:"active"`
	Name   string `json:"name"`
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (total int64, results []*FindResult, err error) {
	do := x.Db.Model(&model.Org{}).WithContext(ctx)
	if dto.Q != "" {
		do = do.Where(`name like ?`, dto.GetKeyword())
	}

	if err = do.Count(&total).Error; err != nil {
		return
	}

	results = make([]*FindResult, 0)
	ctx = common.SetPipe(ctx, common.NewFindPipe())
	if err = dto.Find(ctx, do, &results); err != nil {
		return
	}
	return
}
```

### find_by_id.go

```go
package orgs

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) FindById(ctx context.Context, c *app.RequestContext) {
	var dto common.FindByIdDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	data, err := x.OrgsX.FindById(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type FindByIdResult struct {
	ID     string `json:"id"`
	Active bool   `json:"active"`
	Name   string `json:"name"`
}

func (x *Service) FindById(ctx context.Context, user *common.IAMUser, dto common.FindByIdDto) (result FindByIdResult, err error) {
	do := x.Db.Model(model.Org{}).WithContext(ctx)
	ctx = common.SetPipe(ctx, common.NewFindByIdPipe())
	if err = dto.Take(ctx, do, &result); err != nil {
		return
	}
	return
}
```

### create.go

```go
package orgs

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID     string `json:"-"`
	Name   string `json:"name" vd:"required"`
	Active *bool  `json:"active" vd:"required"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	dto.ID = help.SID()
	if err := x.OrgsX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := model.Org{
		ID:     dto.ID,
		Active: *dto.Active,
		Name:   dto.Name,
	}
	if err = x.Db.WithContext(ctx).Create(&data).Error; err != nil {
		return
	}
	return
}
```

### update.go

```go
package orgs

import (
	"context"
	"time"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/kainonly/go/help"
)

type UpdateDto struct {
	ID     string `json:"id" vd:"required"`
	Name   string `json:"name" vd:"required"`
	Active *bool  `json:"active" vd:"required"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.OrgsX.Update(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	if err = x.Db.Model(model.Org{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(utils.H{
			`updated_at`: time.Now(),
			`active`:     *dto.Active,
			`name`:       dto.Name,
		}).Error; err != nil {
		return
	}
	return
}
```

### delete.go

```go
package orgs

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

func (x *Controller) Delete(ctx context.Context, c *app.RequestContext) {
	var dto common.DeleteDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.OrgsX.Delete(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Delete(ctx context.Context, user *common.IAMUser, dto common.DeleteDto) (err error) {
	return x.Db.WithContext(ctx).Delete(model.Org{}, dto.IDs).Error
}
```

## 通用 DTO 类型

| DTO | 说明 |
|-----|------|
| `common.FindDto` | 列表查询，提供 `Q` 和 `GetKeyword()` |
| `common.FindByIdDto` | 单条查询，提供 `Take()` |
| `common.SearchDto` | 异步搜索 |
| `common.ExistsDto` | 异步验证 |
| `common.DeleteDto` | 批量删除（`IDs []string`） |

## 错误处理与响应规范

| 场景 | 处理方式 |
|------|----------|
| 错误 | `c.Error(err)` |
| 无数据成功 | `c.JSON(200, help.Ok())` |
| 有数据成功 | `c.JSON(200, data)` |
| 列表总数 | `c.Header("x-total", ...)` |
| 获取用户 | `common.GetIAM(c)` |
| ID 生成 | `help.SID()` |

## Service 可用依赖

| 字段 | 类型 | 说明 |
|------|------|------|
| `x.Db` | `*gorm.DB` | 数据库连接 |
| `x.RDb` | `*redis.Client` | Redis 连接 |
| `x.V` | `*Values` | 应用配置 |
| `x.Captcha` | `*captcha.Captcha` | 验证码服务 |
| `x.Locker` | `*locker.Locker` | 分布式锁服务 |

## common.IAMUser 结构

```go
type IAMUser struct {
	ID     string `json:"id"`
	OrgID  string `json:"org_id"`
	RoleID string `json:"role_id"`
	Active bool   `json:"active"`
	Ip     string `json:"-"`
}
```

## 新增模块检查清单

### 模块文件创建

- [ ] 创建 `api/{module}/` 目录
- [ ] 创建 `common.go`（Provides、Controller、Service）
- [ ] 创建 `find.go`（FindDto、FindResult、Find）
- [ ] 创建 `find_by_id.go`（FindByIdResult、FindById）
- [ ] 创建 `create.go`（CreateDto、Create）
- [ ] 创建 `update.go`（UpdateDto、Update）
- [ ] 创建 `delete.go`（Delete）

### api.go 注册

- [ ] 导入模块：`"server/api/{module}"`
- [ ] 添加 Provides：`{module}.Provides`
- [ ] 添加 Controller：`{Module} *{module}.Controller`
- [ ] 注册路由到认证组

### 代码规范

- [ ] 使用 `c.Error(err)` 处理错误
- [ ] 使用 `common.GetIAM(c)` 获取用户
- [ ] 无数据返回 `help.Ok()`
- [ ] 列表设置 `x-total` header
