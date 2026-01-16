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
- **Python 3.10+**
- **Docker**

### 配置文件

1. 复制 `.env.example` 为 `.env`
2. 填入数据库连接信息
3. `.env` 不提交到版本控制

## 核心命令

所有命令通过 `scripts/migrate.py` 执行：

| 命令 | 用途 | 示例 |
|------|------|------|
| `diff` | 对比模型变更，生成迁移文件 | `scripts/migrate.py diff` |
| `apply` | 应用迁移到数据库 | `scripts/migrate.py apply` |
| `status` | 查看迁移状态 | `scripts/migrate.py status` |
| `inspect` | 查看数据库当前 schema | `scripts/migrate.py inspect` |
| `push` | 直接推送 schema（仅开发环境） | `scripts/migrate.py push` |
| `hash` | 重新生成迁移文件哈希 | `scripts/migrate.py hash` |

### 环境切换

```bash
# 本地环境（默认）
scripts/migrate.py diff

# 生产环境
scripts/migrate.py --env prod apply
```

## 标准工作流

### 新增模型后

```bash
# 1. 在 cmd/atlas-loader/main.go 注册模型（必须先做）
# models 切片和 tableComments 映射

# 2. 生成迁移文件
scripts/migrate.py diff

# 3. 检查生成的 SQL（人工审核）
cat migrations/*.sql | tail -50

# 4. 应用迁移
scripts/migrate.py apply

# 5. 验证状态
scripts/migrate.py status
```

### 修改模型字段后

```bash
# 1. 生成迁移
scripts/migrate.py diff

# 2. 审核 SQL（特别关注破坏性变更）
cat migrations/*.sql | tail -50

# 3. 应用迁移
scripts/migrate.py apply
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
| `--env prod` | 生产环境操作 |
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
scripts/migrate.py hash
```

### Atlas 命令未找到

```bash
curl -sSf https://atlasgo.sh | sh
```

### 迁移状态不一致

```bash
# 查看当前状态
scripts/migrate.py status

# 查看数据库实际 schema
scripts/migrate.py inspect
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
- [ ] 执行 `scripts/migrate.py diff`
- [ ] 人工审核 SQL 文件
- [ ] 执行 `scripts/migrate.py apply`

### 修改模型后

- [ ] 执行 `scripts/migrate.py diff`
- [ ] 审核 SQL（关注破坏性变更）
- [ ] 执行 `scripts/migrate.py apply`

### 添加新依赖后

- [ ] 执行 `go mod tidy`
