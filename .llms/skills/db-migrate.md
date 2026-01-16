# 数据库迁移规范

> 当用户要求执行数据库迁移、查看迁移状态或处理迁移问题时，使用此规范。

## 触发条件

- 用户创建或修改了 GORM 模型后
- 用户要求"迁移数据库"、"同步表结构"
- 用户遇到迁移相关问题

## 环境准备

### 前置依赖

- **Go 1.24+**
- **Atlas CLI**：`curl -sSf https://atlasgo.sh | sh`
- **Docker**

### 配置文件

迁移工具从 `config/values.yml` 读取数据库配置：

```yaml
database:
  dsn: host=localhost user=postgres password=your-password dbname=litho port=5432 TimeZone=Asia/Shanghai sslmode=disable
```

## 核心命令

迁移工具位于 `cmd/migrate/main.go`，需要先编译：

```bash
# 编译迁移工具
go build -o migrate ./cmd/migrate

# 或直接运行
go run ./cmd/migrate <command>
```

| 命令 | 用途 | 示例 |
|------|------|------|
| `diff` | 对比模型变更，生成迁移文件 | `./migrate diff` |
| `apply` | 应用迁移到数据库 | `./migrate apply` |
| `status` | 查看迁移状态 | `./migrate status` |
| `inspect` | 查看数据库当前 schema | `./migrate inspect` |
| `push` | 直接推送 schema（仅开发环境） | `./migrate push` |
| `hash` | 重新生成迁移文件哈希 | `./migrate hash` |

### 使用不同配置文件

```bash
# 默认使用 config/values.yml
./migrate diff

# 使用指定配置文件
./migrate -config config/prod.yml apply
```

### 快捷方式（使用 go run）

```bash
go run ./cmd/migrate diff
go run ./cmd/migrate apply
go run ./cmd/migrate -config config/prod.yml apply
```

## 标准工作流

### 新增模型后

```bash
# 1. 在 cmd/atlas-loader/main.go 注册模型（必须先做）
# models 切片和 tableComments 映射

# 2. 生成迁移文件
go run ./cmd/migrate diff

# 3. 检查生成的 SQL（人工审核）
cat migrations/*.sql | tail -50

# 4. 应用迁移
go run ./cmd/migrate apply

# 5. 验证状态
go run ./cmd/migrate status
```

### 修改模型字段后

```bash
# 1. 生成迁移
go run ./cmd/migrate diff

# 2. 审核 SQL（特别关注破坏性变更）
cat migrations/*.sql | tail -50

# 3. 应用迁移
go run ./cmd/migrate apply
```

## Atlas Loader 注册

新增模型必须在 `cmd/atlas-loader/main.go` 中注册：

```go
models := []any{
    &model.User{},
    &model.Org{},
    // 新增模型在此处添加
    &model.NewEntity{},
}

var tableComments = map[string]string{
    "user": "用户表",
    "org":  "组织表",
    // 新增表注释
    "new_entity": "新实体表",
}
```

**重要**：未注册的模型不会被 Atlas 识别，迁移时将被忽略。

## 危险操作警告

以下操作需要用户明确确认：

| 操作 | 风险 |
|------|------|
| 使用生产配置 | 生产环境操作 |
| 删除列或表 | 数据丢失 |
| 修改列类型 | 数据转换失败 |
| 删除索引 | 查询性能下降 |

执行前必须：
1. 明确告知用户操作风险
2. 等待用户确认
3. 建议先在本地环境验证

## 常见问题

### 迁移文件哈希不匹配

```bash
go run ./cmd/migrate hash
```

### Atlas 命令未找到

```bash
curl -sSf https://atlasgo.sh | sh
```

### 迁移状态不一致

```bash
# 查看当前状态
go run ./cmd/migrate status

# 查看数据库实际 schema
go run ./cmd/migrate inspect
```

## Go 项目常用命令

### 依赖管理

```bash
go mod download    # 下载依赖
go mod tidy        # 整理依赖
go get -u pkg      # 更新依赖
```

### 代码生成

```bash
go generate ./...  # Wire 依赖注入生成
```

### 运行与测试

```bash
go run .           # 运行应用
go test ./...      # 运行测试
go test -cover ./...  # 带覆盖率测试
```

## 自动化检查清单

### 新增模型后

- [ ] 在 `models` 切片中注册模型
- [ ] 在 `tableComments` 中添加表注释
- [ ] 执行 `go run ./cmd/migrate diff`
- [ ] 人工审核 SQL 文件
- [ ] 执行 `go run ./cmd/migrate apply`

### 修改模型后

- [ ] 执行 `go run ./cmd/migrate diff`
- [ ] 审核 SQL（关注破坏性变更）
- [ ] 执行 `go run ./cmd/migrate apply`

### 添加新依赖后

- [ ] 执行 `go mod tidy`
