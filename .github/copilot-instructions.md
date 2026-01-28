# Litho Server 项目规范

Go 后端服务，使用 Hertz + GORM + PostgreSQL + Redis。

## GORM 模型规范

- 主键：`ID string` + `gorm:"primaryKey;column:id;type:bigint"`
- 时间：`type:timestamptz`（非 `timestamp`）
- 文本：`type:text`（非 `varchar`）
- JSONB：对象用 `common.M`，数组用 `common.A`
- 必须实现 `TableName()` 方法
- 禁止：`gorm.Model`、外键关联、`autoIncrement`、`serial`

```go
type Entity struct {
    ID        string    `gorm:"primaryKey;column:id;type:bigint"`
    CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;index"`
    UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
    Name      string    `gorm:"column:name;type:text;not null;comment:名称"`
}

func (Entity) TableName() string {
    return "entity"
}
```

## API 模块规范

- 依赖注入：`wire.NewSet` + `wire.Struct(new(Type), "*")`
- 错误处理：`c.Error(err)`
- 成功响应：无数据用 `help.Ok()`，有数据直接返回
- 列表接口：`c.Header("x-total", strconv.Itoa(int(total)))`

```go
// Controller 方法
func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
    var dto CreateDto
    if err := c.BindAndValidate(&dto); err != nil {
        c.Error(err)
        return
    }
    if err := x.ModuleX.Create(ctx, dto); err != nil {
        c.Error(err)
        return
    }
    c.JSON(200, help.Ok())
}
```

## 路由约定

| 路由 | Method | 说明 |
|------|--------|------|
| `/{module}/:id` | GET | 根据 ID 获取 |
| `/{module}` | GET | 获取列表 |
| `/{module}/create` | POST | 创建 |
| `/{module}/update` | POST | 更新 |
| `/{module}/delete` | POST | 删除 |

## 命令

```bash
go run ./cmd/migrate diff   # 生成迁移
go run ./cmd/migrate apply  # 应用迁移
go generate ./...           # Wire 生成
```

详细规范参见 `.llms/` 目录。
