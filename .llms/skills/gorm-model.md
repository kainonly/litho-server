# GORM 模型生成规范

> 当用户要求创建或修改数据库模型时，使用此规范生成符合项目标准的 GORM 模型代码。

## 触发条件

- 用户提到"创建模型"、"新建表"、"添加实体"
- 用户提供思维导图或表结构定义
- 用户要求修改现有模型字段

## 命名规范

| 元素 | 格式 | 示例 |
|------|------|------|
| 表名 | 小写字母 + 下划线 | `user`、`role_permission` |
| 列名 | 小写字母 + 下划线 | `created_at`、`org_id` |
| Go 字段名 | PascalCase | `CreatedAt`、`OrgID` |
| 关联字段 | `${entity}_id` / `${Entity}ID` | `org_id` / `OrgID` |
| 层级自关联 | `pid` / `PID` | - |

## 基本结构模板

```go
package model

import "time"

// EntityName 实体注释（来自思维导图）
type EntityName struct {
	ID        string    `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
	// 业务字段按思维导图中的顺序排列...
	FieldName string `gorm:"column:field_name;type:text;not null;comment:字段注释"` // 字段注释
}

func (EntityName) TableName() string {
	return "entity_name"
}
```

## 字段定义规则

### 主键
- 类型：`ID string`
- Tag：`gorm:"primaryKey;column:id;type:bigint"`
- 说明：使用 `string` 避免 JavaScript 精度丢失

### 时间字段
```go
CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;index"`
UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
```

### 关联字段（逻辑外键）
```go
OrgID string `gorm:"column:org_id;type:bigint;not null;comment:组织ID"`
```

### 文本字段
```go
Name string `gorm:"column:name;type:text;not null;comment:名称"`
```

### 布尔字段
```go
Active bool `gorm:"column:active;not null;default:true;comment:是否启用"`
```

### 状态字段
```go
State int16 `gorm:"column:state;not null;default:1;comment:状态"` // 从1开始，负值表示异常
```

### 排序字段
```go
Sort int16 `gorm:"column:sort;not null;default:0;comment:排序权重"`
```

### JSONB 字段
```go
Metadata common.M `gorm:"column:metadata;type:jsonb;not null;default:'{}'"`  // 对象
Tags     common.A `gorm:"column:tags;type:jsonb;not null;default:'[]'"`      // 数组
```

## Tag 格式顺序

`gorm:"primaryKey;column:xxx;type:xxx;not null;default:xxx;index/uniqueIndex;comment:xxx"`

## 索引规范

### 单字段索引
```go
Email string `gorm:"column:email;type:text;not null;uniqueIndex;comment:邮箱"`
```

### 联合索引
```go
// 联合唯一索引: (org_id, code)
OrgID string `gorm:"column:org_id;type:bigint;not null;uniqueIndex:idx_resource_org_code,priority:1;comment:组织ID"`
Code  string `gorm:"column:code;type:text;not null;uniqueIndex:idx_resource_org_code,priority:2;comment:编码"`
```

## 字段排列顺序

1. `ID`
2. `CreatedAt`、`UpdatedAt`
3. 业务字段（按思维导图顺序）

## 禁止事项

- ❌ 外键关联（`foreignKey`、`references`）
- ❌ `gorm.Model` 嵌入
- ❌ `serial` 或 `autoIncrement`
- ❌ `varchar(n)` 或 `char(n)`
- ❌ 省略 `column` 标签

## 完整示例

### 基础业务表
```go
// User 系统用户
type User struct {
	ID        string    `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
	Email     string    `gorm:"column:email;type:text;not null;uniqueIndex;comment:邮箱地址"`
	Phone     string    `gorm:"column:phone;type:text;not null;index;comment:手机号码"`
	Name      string    `gorm:"column:name;type:text;not null;comment:用户姓名"`
	Password  string    `gorm:"column:password;type:text;not null;comment:登录密码"`
	Avatar    string    `gorm:"column:avatar;type:text;not null;comment:头像地址"`
	Active    bool      `gorm:"column:active;not null;default:true;comment:是否启用"`
}

func (User) TableName() string {
	return "user"
}
```

### 关联表（多对多中间表）
```go
// RolePermission 角色权限关联
type RolePermission struct {
	RoleID       string `gorm:"primaryKey;column:role_id;type:bigint;comment:角色ID"`
	PermissionID string `gorm:"primaryKey;column:permission_id;type:bigint;comment:权限ID"`
}

func (RolePermission) TableName() string {
	return "role_permission"
}
```

## 检查清单

生成模型后确认：

- [ ] 主键和关联字段使用 `string` + `type:bigint`
- [ ] 包含 `CreatedAt` 和 `UpdatedAt`（`type:timestamptz`）
- [ ] `CreatedAt` 包含 `index` 索引
- [ ] 业务字段按思维导图顺序排列
- [ ] 业务字段包含 `comment:xxx` 标签
- [ ] 文本字段使用 `type:text`
- [ ] 所有字段显式指定 `column`
- [ ] 核心字段包含 `not null`
- [ ] 实现 `TableName()` 方法
- [ ] 无外键关联、无 `gorm.Model` 嵌入

## 后续步骤

模型创建完成后，需要：

1. 在 `cmd/atlas-loader/main.go` 的 `models` 切片中注册模型
2. 在 `tableComments` 中添加表注释
3. 执行数据库迁移（参见 `db-migrate` skill）
