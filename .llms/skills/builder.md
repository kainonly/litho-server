# Builder 模型生成规范

> 当用户要求从数据库生成或更新 GORM 模型时，使用此规范。

## 触发条件

- 用户创建或修改了数据库表结构后
- 用户要求"生成模型"、"同步模型"、"更新 model"
- 用户新增了数据库表

## 工作原理

Builder 使用 GORM Gen 连接 PostgreSQL 数据库，反向读取表结构并生成 Go model 文件到 `./model` 目录。每次运行会**删除并重新生成**整个 `./model` 目录。

### 配置文件

Builder 从 `config/values.yml` 读取数据库配置：

```yaml
database:
  dsn: host=localhost user=postgres password=your-password dbname=litho port=5432 TimeZone=Asia/Shanghai sslmode=disable
```

## 核心命令

```bash
# 运行 builder 生成 model
go run ./cmd/builder
```

## 代码结构

### cmd/builder/common.go

通用生成函数，负责：
- 读取数据库配置
- 连接 PostgreSQL
- 配置 GORM Gen（nullable、coverable、signable）
- 执行代码生成
- 清理临时 query 目录

### cmd/builder/main.go

主入口，定义：
- **表名到结构体的映射**：通过 `g.GenerateModelAs("table_name", "StructName")` 注册
- **全局字段类型覆盖**：通过 `gen.FieldType("column_name", "go_type")` 指定
- **JSONB 类型映射**：jsonb 列映射为 `common.Object`

## 标准工作流

### 新增表后

```bash
# 1. 在数据库中创建表（手动或通过 SQL）

# 2. 在 cmd/builder/main.go 中注册表映射
#    在 g.ApplyBasic() 中添加:
#    g.GenerateModelAs("new_table", "NewTable"),

# 3. 如有新的 bigint 关联字段，在 gen.FieldType() 中添加类型覆盖
#    gen.FieldType("new_table_id", "string"),

# 4. 运行 builder
go run ./cmd/builder

# 5. 检查生成的 model 文件
cat ./model/new_table.gen.go
```

### 修改表结构后

```bash
# 1. 在数据库中修改表结构

# 2. 重新运行 builder（会重新生成所有 model）
go run ./cmd/builder
```

## 注册新表示例

在 `cmd/builder/main.go` 的 `g.ApplyBasic()` 中添加：

```go
g.ApplyBasic(
    g.GenerateModelAs("cap", "Cap"),
    g.GenerateModelAs("org", "Org"),
    g.GenerateModelAs("resource", "Resource"),
    g.GenerateModelAs("role", "Role"),
    g.GenerateModelAs("route", "Route"),
    g.GenerateModelAs("user", "User"),
    // 新增表在此处添加
    g.GenerateModelAs("new_table", "NewTable"),
)
```

如果新表包含 bigint 类型的关联字段，需要添加类型覆盖：

```go
gen.FieldType("id", "string"),
gen.FieldType("pid", "string"),
gen.FieldType("org_id", "string"),
gen.FieldType("role_id", "string"),
gen.FieldType("user_id", "string"),
// 新增关联字段
gen.FieldType("new_table_id", "string"),
```

## 重要说明

- Builder 每次运行会**删除**整个 `./model` 目录后重新生成，不要手动修改生成的文件
- JSONB 列自动映射为 `common.Object` 类型
- bigint 类型的 id/关联字段需要通过 `gen.FieldType()` 覆盖为 `string`（避免 JS 精度丢失）
- 生成的文件位于 `./model/*.gen.go`

## 自动化检查清单

### 新增表后

- [ ] 在数据库中创建表
- [ ] 在 `cmd/builder/main.go` 的 `g.ApplyBasic()` 中注册表映射
- [ ] 添加必要的 `gen.FieldType()` 类型覆盖
- [ ] 执行 `go run ./cmd/builder`
- [ ] 检查生成的 model 文件

### 修改表结构后

- [ ] 在数据库中修改表结构
- [ ] 执行 `go run ./cmd/builder`
- [ ] 检查生成的 model 文件

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
