# 项目规范

请遵循 `.llms/llms.txt` 中的 PostgreSQL 和 GORM 模型规范。

## 关键要点

### GORM 模型
- 主键和关联字段使用 `string` 类型 + `type:bigint`（避免 JS 精度丢失）
- 文本字段使用 `type:text`
- JSONB 对象使用 `common.M`，数组使用 `common.A`
- 必须实现 `TableName()` 方法
- 禁止使用 `gorm.Model`、外键关联、`autoIncrement`

### PostgreSQL
- 表名、字段名全小写 + 下划线
- 时间字段使用 `timestamptz`
- 使用 `NOT EXISTS` 而非 `NOT IN`
- 时间范围查询使用 `>=` 和 `<` 而非 `BETWEEN`

详细规范请参考 `.llms/llms.txt`。
