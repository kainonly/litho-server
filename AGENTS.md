# Litho Server 项目规范

> 基于 Go 的后端服务项目，详细规范参见 `.llms/` 目录。

## 技术栈

- Web 框架：Hertz (CloudWeGo)
- ORM：GORM
- 数据库：PostgreSQL
- 缓存：Redis
- 依赖注入：Wire (github.com/goforj/wire)
- 数据库迁移：Atlas

## 项目结构

```
server/
├── api/                 # API 模块
│   ├── api.go           # 路由聚合
│   ├── auth.go          # 认证中间件
│   └── {module}/        # 各业务模块（common.go + 方法文件）
├── cmd/
│   ├── migrate/         # 迁移命令
│   └── atlas-loader/    # Atlas 迁移加载器
├── common/              # 公共组件
├── model/               # GORM 模型
├── migrations/          # 迁移文件
└── config/              # 配置文件
```

## 核心编码规范

### GORM 模型

```go
type Entity struct {
    ID        string    `gorm:"primaryKey;column:id;type:bigint"`
    CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;index"`
    UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
    // 业务字段...
    Name string `gorm:"column:name;type:text;not null;comment:名称"`
}

func (Entity) TableName() string {
    return "entity"
}
```

**规则**：
- 主键/关联字段：`string` + `type:bigint`
- 时间字段：`type:timestamptz`
- 文本字段：`type:text`
- JSONB：对象用 `common.M`，数组用 `common.A`
- 禁止：`gorm.Model`、外键关联、`autoIncrement`

### API 模块

每个模块包含 `common.go` 和各方法文件：

```go
// common.go
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

**错误处理**：`c.Error(err)`
**成功响应**：无数据用 `help.Ok()`，有数据直接返回
**列表接口**：设置 `c.Header("x-total", ...)`

### 路由约定

| 路由 | Method | 说明 |
|------|--------|------|
| `/{module}/:id` | GET | 根据 ID 获取 |
| `/{module}` | GET | 获取列表 |
| `/{module}/create` | POST | 创建 |
| `/{module}/update` | POST | 更新 |
| `/{module}/delete` | POST | 删除 |

## 常用命令

```bash
go run ./cmd/migrate diff     # 生成迁移
go run ./cmd/migrate apply    # 应用迁移
go mod tidy                   # 整理依赖
go generate ./...             # Wire 生成
go run .                      # 运行
```

## 新增功能流程

1. 创建模型 → `.llms/skills/gorm-model.md`
2. 注册模型 → `cmd/atlas-loader/main.go`
3. 执行迁移 → `.llms/skills/db-migrate.md`
4. 创建 API → `.llms/skills/api-module.md`
5. 注册路由 → `api/api.go`

## Skills 参考

- GORM 模型：`.llms/skills/gorm-model.md`
- 数据库迁移：`.llms/skills/db-migrate.md`
- API 模块：`.llms/skills/api-module.md`
