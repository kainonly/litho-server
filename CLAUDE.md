# Litho Server 项目规范

> 基于 Go 的后端服务项目，详细规范参见 `.llms/` 目录。

## 技术栈

- **Web 框架**：Hertz (CloudWeGo)
- **ORM**：GORM
- **数据库**：PostgreSQL
- **缓存**：Redis
- **依赖注入**：Wire (github.com/goforj/wire)
- **模型生成**：GORM Gen（builder）

## 项目结构

```
server/
├── api/                 # API 模块
│   ├── api.go           # 路由聚合
│   ├── auth.go          # 认证中间件
│   └── {module}/        # 各业务模块
├── cmd/
│   └── builder/         # GORM Gen 模型生成器
├── common/              # 公共组件
└── model/               # GORM 模型（由 builder 生成）
```

## 核心约定

### GORM 模型

- 模型由 builder 从数据库反向生成，不手动编写
- 主键和关联字段使用 `string` 类型（通过 `gen.FieldType` 覆盖 bigint）
- JSONB 列映射为 `common.Object`
- 禁止手动修改 `./model` 目录下的生成文件

### API 模块

- 依赖注入使用 `wire.NewSet` + `wire.Struct(new(Type), "*")`
- 错误处理统一使用 `c.Error(err)`
- 响应格式：无数据返回 `help.Ok()`，有数据直接返回
- 列表接口设置 `x-total` header

### PostgreSQL

- 表名、字段名全小写 + 下划线
- 使用 `NOT EXISTS` 而非 `NOT IN`
- 时间范围查询使用 `>=` 和 `<` 而非 `BETWEEN`

## 常用命令

```bash
# 模型生成
go run ./cmd/builder              # 从数据库生成 model

# Go
go mod tidy                   # 整理依赖
go generate ./...             # Wire 生成
go run .                      # 运行
go test ./...                 # 测试
```

## 新增功能流程

1. **创建表** → 在数据库中创建表结构
2. **注册表** → 在 `cmd/builder/main.go` 中添加表映射
3. **生成模型** → 运行 `go run ./cmd/builder`
4. **创建 API** → 参考 `.llms/skills/api-module.md`
5. **注册路由** → 在 `api/api.go`

## Skills 索引

| Skill | 文件 | 触发场景 |
|-------|------|----------|
| GORM 模型规范 | `.llms/skills/gorm-model.md` | 了解模型字段规范、命名约定 |
| Builder 模型生成 | `.llms/skills/builder.md` | 从数据库生成/更新 model 代码 |
| API 模块 | `.llms/skills/api-module.md` | 创建 API 接口、新增模块、路由 |
