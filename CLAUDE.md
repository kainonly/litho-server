# Litho Server 项目规范

> 基于 Go 的后端服务项目，详细规范参见 `.llms/` 目录。

## 技术栈

- **Web 框架**：Hertz (CloudWeGo)
- **ORM**：GORM
- **数据库**：PostgreSQL
- **缓存**：Redis
- **依赖注入**：Wire (github.com/goforj/wire)
- **数据库迁移**：Atlas

## 项目结构

```
server/
├── api/                 # API 模块
│   ├── api.go           # 路由聚合
│   ├── auth.go          # 认证中间件
│   └── {module}/        # 各业务模块
├── cmd/
│   └── atlas-loader/    # Atlas 迁移加载器
├── common/              # 公共组件
├── model/               # GORM 模型
└── migrations/          # 迁移文件
```

## 核心约定

### GORM 模型

- 主键和关联字段使用 `string` 类型 + `type:bigint`（避免 JS 精度丢失）
- 时间字段必须使用 `type:timestamptz`
- 文本字段使用 `type:text`
- JSONB 对象使用 `common.M`，数组使用 `common.A`
- 必须实现 `TableName()` 方法
- 禁止使用 `gorm.Model`、外键关联、`autoIncrement`

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
# 迁移
go run ./cmd/migrate diff     # 生成迁移
go run ./cmd/migrate apply    # 应用迁移
go run ./cmd/migrate status   # 查看状态

# Go
go mod tidy                   # 整理依赖
go generate ./...             # Wire 生成
go run .                      # 运行
go test ./...                 # 测试
```

## 新增功能流程

1. **创建模型** → 参考 `.llms/skills/gorm-model.md`
2. **注册模型** → 在 `cmd/atlas-loader/main.go`
3. **执行迁移** → 参考 `.llms/skills/db-migrate.md`
4. **创建 API** → 参考 `.llms/skills/api-module.md`
5. **注册路由** → 在 `api/api.go`

## Skills 索引

| Skill | 文件 | 触发场景 |
|-------|------|----------|
| GORM 模型生成 | `.llms/skills/gorm-model.md` | 创建/修改数据库模型、表结构 |
| 数据库迁移 | `.llms/skills/db-migrate.md` | 执行迁移、同步表结构、迁移问题 |
| API 模块 | `.llms/skills/api-module.md` | 创建 API 接口、新增模块、路由 |
