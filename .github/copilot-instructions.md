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

### 模块常量（common.go）

每个模块在 `common.go` 中声明两个常量：

```go
const (
    Key   = "roles"  // 模块标识，用于缓存 key 等
    Label = "权限组" // 模块中文名，用于日志、审计
)
```

### 操作常量

每个操作文件声明对应的操作常量（紧跟 Dto 类型之后）：

```go
const ICreate = "新增"   // create.go
const IUpdate = "更新"   // update.go
const IDelete = "删除"   // delete.go
const ISort   = "排序"   // sort.go
// 其他操作类文件同理，如 IRegroup = "重新分组"
```

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
| `/{module}/_search` | GET | 异步搜索 |
| `/{module}/_exists` | GET | 异步验证 |
| `/{module}/create` | POST | 创建 |
| `/{module}/update` | POST | 更新 |
| `/{module}/delete` | POST | 删除 |
| `/{module}/get_{noun}` | GET | 读取单个字段（如 `get_strategy`） |
| `/{module}/set_{noun}` | POST | 更新单个字段（如 `set_strategy`） |
| `/{module}/{verb}` | POST | 其他操作（如 `sort`、`regroup`） |

## 模型生成

模型由 builder 从数据库反向生成（GORM Gen），禁止手动修改 `./model` 目录下的生成文件。

```bash
go run ./cmd/builder            # 从数据库生成 model
```

新增表需在 `cmd/builder/main.go` 的 `g.ApplyBasic()` 中注册表映射。

## 命令

```bash
go run ./cmd/builder            # 从数据库生成 model
go generate ./...           # Wire 生成
go mod tidy                 # 整理依赖
```

详细规范参见 `.llms/` 目录。
