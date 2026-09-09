// Package common 提供通用 CRUD 操作与 REST API 构建工具。
//
// 采用 "Pipe" 模式：配置对象经 SetPipe 存入 context，由 DTO 方法取出以定制查询行为。
//
//   - ExistsDto + ExistsPipe：按字段检查记录是否存在
//   - FindDto + FindPipe：分页列表查询
//   - FindByIdDto + FindByIdPipe：按 ID 获取单条记录
//   - SearchDto + SearchPipe：轻量搜索/自动补全
//
// 安全性：列名经白名单或格式校验，ID 经格式校验，全部使用参数化查询防注入。
package common

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// pipeKey 用作 context 键的自定义类型，避免键冲突。
type pipeKey struct{}

// validColumnName 列名校验，防 SQL 注入。
var validColumnName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// validSnowflakeID 雪花 ID 校验（纯数字，通常 18-19 位）。
var validSnowflakeID = regexp.MustCompile(`^[0-9]+$`)

// IDValidator 自定义 ID 校验函数。
type IDValidator func(id string) bool

// DefaultIDValidator 默认按雪花 ID（纯数字）校验。
var DefaultIDValidator IDValidator = func(id string) bool {
	return validSnowflakeID.MatchString(id)
}

// Controller 标准 CRUD 接口。
type Controller interface {
	Create(ctx context.Context, c *app.RequestContext)
	Find(ctx context.Context, c *app.RequestContext)
	FindById(ctx context.Context, c *app.RequestContext)
	Update(ctx context.Context, c *app.RequestContext)
	Delete(ctx context.Context, c *app.RequestContext)
}

// SetPipe 将管道配置存入 context。
func SetPipe(ctx context.Context, i any) context.Context {
	return context.WithValue(ctx, pipeKey{}, i)
}

// getPipe 从 context 取出指定类型的管道。
func getPipe[T any](ctx context.Context) (T, bool) {
	v, ok := ctx.Value(pipeKey{}).(T)
	return v, ok
}

// ToOrderBy 排序方向映射："1" 升序（空串），"-1" 降序。
var ToOrderBy = map[string]string{
	"1":  "",
	"-1": "desc",
}

// ExistsDto 存在性检查 DTO。
// 查询参数：key 列名，q 值。
type ExistsDto struct {
	Key string `query:"key,omitempty"`
	Q   string `query:"q,omitempty"`
}

// ExistsPipe 存在性检查的字段白名单。
type ExistsPipe struct {
	fields map[string]bool
}

// NewExistsPipe 创建 ExistsPipe，仅允许检查指定的字段。
func NewExistsPipe(keys ...string) *ExistsPipe {
	fields := make(map[string]bool)
	for _, key := range keys {
		fields[key] = true
	}
	return &ExistsPipe{
		fields: fields,
	}
}

// Get 从 context 取出 ExistsPipe。
func (x *ExistsDto) Get(ctx context.Context) (*ExistsPipe, error) {
	p, ok := getPipe[*ExistsPipe](ctx)
	if !ok {
		return nil, help.E(0, "上下文中未找到 ExistsPipe")
	}
	return p, nil
}

// ExistsResult 存在性检查响应。
type ExistsResult struct {
	Exists bool `json:"exists"`
}

// Exists 检查指定列是否存在匹配值的记录，列名须经白名单校验。
func (x *ExistsDto) Exists(ctx context.Context, do *gorm.DB) (result ExistsResult, err error) {
	p, err := x.Get(ctx)
	if err != nil {
		return
	}
	// 白名单校验
	if !p.fields[x.Key] {
		err = help.E(0, fmt.Sprintf(`字段 [%s] 不允许进行存在性检查`, x.Key))
		return
	}
	var count int64
	if err = do.
		Where(`? = ?`, clause.Column{Name: x.Key}, x.Q).
		Count(&count).Error; err != nil {
		return
	}
	result = ExistsResult{
		Exists: count != 0,
	}
	return
}

// FindDto 分页列表查询 DTO。
//
// 请求头：x-pagesize 每页记录数（默认 1000，最大 1000），x-page 页码（从 0 起）。
// 查询参数：q 搜索关键字，sort 排序规则 "column:direction"（如 "name:1"、"create_time:-1"）。
type FindDto struct {
	PageSize int64    `header:"x-pagesize" vd:"omitempty,min=0,max=1000"`
	Page     int64    `header:"x-page" vd:"omitempty,min=0"`
	Q        string   `query:"q,omitempty"`
	Sort     []string `query:"sort,omitempty" vd:"omitempty,dive,sort"`
}

// GetPageSize 返回每页记录数，默认 1000。
func (x *FindDto) GetPageSize() int {
	if x.PageSize == 0 {
		x.PageSize = 1000
	}
	return int(x.PageSize)
}

// GetOffset 计算分页偏移量。
func (x *FindDto) GetOffset() int {
	return int(x.Page) * int(x.PageSize)
}

// GetKeyword 返回包裹 LIKE 通配符的关键字，如 "test" -> "%test%"。
func (x *FindDto) GetKeyword() string {
	return fmt.Sprintf(`%%%s%%`, x.Q)
}

// FindPipe Find 查询配置。
type FindPipe struct {
	ts       bool            // 自动省略时间戳列（create_time、update_time）
	sort     bool            // 应用排序
	page     bool            // 应用分页
	keys     []string        // 指定返回的列
	omit     []string        // 排除的列
	sortable map[string]bool // 可排序列白名单
}

// Get 从 context 取出 FindPipe。
func (x *FindDto) Get(ctx context.Context) (*FindPipe, error) {
	p, ok := getPipe[*FindPipe](ctx)
	if !ok {
		return nil, help.E(0, "上下文中未找到 FindPipe")
	}
	return p, nil
}

// NewFindPipe 创建默认配置的 FindPipe：省略时间戳列、
// 未指定排序时按 create_time 倒序、启用分页。
func NewFindPipe() *FindPipe {
	return &FindPipe{
		ts:   true,
		sort: true,
		page: true,
	}
}

// SkipTs 禁用时间戳列自动省略。
func (x *FindPipe) SkipTs() *FindPipe {
	x.ts = false
	return x
}

// SkipSort 禁用排序。
func (x *FindPipe) SkipSort() *FindPipe {
	x.sort = false
	return x
}

// SkipPage 禁用分页，返回全部记录。
func (x *FindPipe) SkipPage() *FindPipe {
	x.page = false
	return x
}

// Select 指定返回的列。
func (x *FindPipe) Select(keys ...string) *FindPipe {
	x.keys = keys
	return x
}

// Omit 指定排除的列，Select 优先。
func (x *FindPipe) Omit(keys ...string) *FindPipe {
	x.omit = keys
	return x
}

// Sortable 设置可排序列白名单。未设置时仅做格式校验，建议显式指定。
func (x *FindPipe) Sortable(keys ...string) *FindPipe {
	x.sortable = make(map[string]bool)
	for _, key := range keys {
		x.sortable[key] = true
	}
	return x
}

// Factory 根据管道配置应用列选择/排除、排序与分页。
// 排序列名经白名单（若配置）或格式校验，防 SQL 注入。
func (x *FindDto) Factory(ctx context.Context, do *gorm.DB) (*gorm.DB, error) {
	p, err := x.Get(ctx)
	if err != nil {
		return nil, err
	}
	if len(p.keys) != 0 {
		do = do.Select(p.keys)
	} else {
		if len(p.omit) == 0 && p.ts {
			do = do.Omit(`create_time`, `update_time`)
		}
		if len(p.omit) != 0 {
			do = do.Omit(p.omit...)
		}
	}

	if p.sort {
		if len(x.Sort) == 0 && p.ts {
			do = do.Order("create_time desc")
		}
		for _, v := range x.Sort {
			rule := strings.Split(v, ":")
			if len(rule) != 2 {
				return nil, help.E(0, fmt.Sprintf(`排序格式无效: %s`, v))
			}
			columnName := rule[0]
			// 白名单优先，否则格式校验
			if len(p.sortable) > 0 {
				if !p.sortable[columnName] {
					return nil, help.E(0, fmt.Sprintf(`列 [%s] 不可排序`, columnName))
				}
			} else {
				if !validColumnName.MatchString(columnName) {
					return nil, help.E(0, fmt.Sprintf(`排序中的列名无效: %s`, columnName))
				}
			}
			order, ok := ToOrderBy[rule[1]]
			if !ok {
				return nil, help.E(0, fmt.Sprintf(`排序方向无效: %s`, rule[1]))
			}
			do = do.Order(fmt.Sprintf(`%s %s`, columnName, order))
		}
	}

	if p.page {
		do = do.Limit(x.GetPageSize()).Offset(x.GetOffset())
	}
	return do, nil
}

// Find 执行分页查询并扫描到指定切片。
func (x *FindDto) Find(ctx context.Context, do *gorm.DB, i any) (err error) {
	db, err := x.Factory(ctx, do)
	if err != nil {
		return
	}
	return db.Find(i).Error
}

// FindByIdDto 按 ID 获取单条记录。
// 路径参数 id；查询参数 full=1 时返回全部字段。
type FindByIdDto struct {
	ID   string `path:"id"`
	Full int    `query:"full,omitempty"`
}

// IsFull 是否为完整模式。
func (x *FindByIdDto) IsFull() bool {
	return x.Full == 1
}

// FindByIdPipe FindById 查询配置，普通/完整模式可分别设置列。
type FindByIdPipe struct {
	ts          bool        // 自动省略时间戳列
	keys        []string    // 普通模式返回的列
	omit        []string    // 普通模式排除的列
	fKeys       []string    // 完整模式返回的列
	fOmit       []string    // 完整模式排除的列
	idValidator IDValidator // ID 校验函数
}

// Get 从 context 取出 FindByIdPipe。
func (x *FindByIdDto) Get(ctx context.Context) (*FindByIdPipe, error) {
	p, ok := getPipe[*FindByIdPipe](ctx)
	if !ok {
		return nil, help.E(0, "上下文中未找到 FindByIdPipe")
	}
	return p, nil
}

// NewFindByIdPipe 创建默认配置的 FindByIdPipe（雪花 ID 校验）。
func NewFindByIdPipe() *FindByIdPipe {
	return &FindByIdPipe{
		ts:          true,
		idValidator: DefaultIDValidator,
	}
}

// SkipTs 禁用时间戳列自动省略。
func (x *FindByIdPipe) SkipTs() *FindByIdPipe {
	x.ts = false
	return x
}

// Select 指定普通模式返回的列。
func (x *FindByIdPipe) Select(keys ...string) *FindByIdPipe {
	x.keys = keys
	return x
}

// Omit 指定普通模式排除的列。
func (x *FindByIdPipe) Omit(keys ...string) *FindByIdPipe {
	x.omit = keys
	return x
}

// FullSelect 指定完整模式返回的列。
func (x *FindByIdPipe) FullSelect(keys ...string) *FindByIdPipe {
	x.fKeys = keys
	return x
}

// FullOmit 指定完整模式排除的列。
func (x *FindByIdPipe) FullOmit(keys ...string) *FindByIdPipe {
	x.fOmit = keys
	return x
}

// SetIDValidator 设置自定义 ID 校验函数。
func (x *FindByIdPipe) SetIDValidator(v IDValidator) *FindByIdPipe {
	x.idValidator = v
	return x
}

// SkipIDValidation 禁用 ID 校验，仅在信任输入来源时使用。
func (x *FindByIdPipe) SkipIDValidation() *FindByIdPipe {
	x.idValidator = nil
	return x
}

// Take 按 ID 获取记录，full 参数决定列配置，ID 先经格式校验。
func (x *FindByIdDto) Take(ctx context.Context, do *gorm.DB, i any) (err error) {
	p, err := x.Get(ctx)
	if err != nil {
		return
	}
	// 校验 ID 格式
	if p.idValidator != nil && !p.idValidator(x.ID) {
		return help.E(0, fmt.Sprintf(`ID 格式无效: %s`, x.ID))
	}
	if !x.IsFull() {
		if len(p.keys) != 0 {
			do = do.Select(p.keys)
		} else {
			if len(p.omit) == 0 && p.ts {
				do = do.Omit(`create_time`, `update_time`)
			}
			if len(p.omit) != 0 {
				do = do.Omit(p.omit...)
			}
		}
	} else {
		if len(p.fKeys) != 0 {
			do = do.Select(p.fKeys)
		} else {
			if len(p.fOmit) != 0 {
				do = do.Omit(p.fOmit...)
			}
		}
	}

	return do.Where(`id = ?`, x.ID).Take(i).Error
}

// SearchDto 轻量搜索 DTO，用于自动补全。
// 查询参数：m 模式，q 关键字，ids 逗号分隔的优先 ID 列表。
type SearchDto struct {
	M   string `query:"m,omitempty"`
	Q   string `query:"q,omitempty"`
	IDs string `query:"ids,omitempty"`
}

// GetKeyword 返回包裹 LIKE 通配符的关键字。
func (x *SearchDto) GetKeyword() string {
	return fmt.Sprintf(`%%%s%%`, x.Q)
}

// SearchPipe Search 查询配置。
type SearchPipe struct {
	keys        []string    // 返回的列（默认 id、name）
	async       bool        // 限制结果数量（异步/自动补全场景）
	idValidator IDValidator // IDs 参数校验函数
}

// SkipAsync 禁用结果数量限制（默认最多 50 条）。
func (x *SearchPipe) SkipAsync() *SearchPipe {
	x.async = false
	return x
}

// NewSearchPipe 创建 SearchPipe，默认返回 id/name 列、限 50 条、校验 IDs。
func NewSearchPipe(keys ...string) *SearchPipe {
	search := &SearchPipe{
		async:       true,
		idValidator: DefaultIDValidator,
	}
	if len(keys) == 0 {
		search.keys = []string{"id", "name"}
	} else {
		search.keys = keys
	}
	return search
}

// SetIDValidator 为 IDs 参数设置自定义校验函数。
func (x *SearchPipe) SetIDValidator(v IDValidator) *SearchPipe {
	x.idValidator = v
	return x
}

// SkipIDValidation 禁用 IDs 参数校验。
func (x *SearchPipe) SkipIDValidation() *SearchPipe {
	x.idValidator = nil
	return x
}

// Get 从 context 取出 SearchPipe。
func (x *SearchDto) Get(ctx context.Context) (*SearchPipe, error) {
	p, ok := getPipe[*SearchPipe](ctx)
	if !ok {
		return nil, help.E(0, "上下文中未找到 SearchPipe")
	}
	return p, nil
}

// Factory 根据管道配置应用列选择与结果数量限制。
func (x *SearchDto) Factory(ctx context.Context, do *gorm.DB) (*gorm.DB, error) {
	p, err := x.Get(ctx)
	if err != nil {
		return nil, err
	}
	if p.async {
		do = do.Limit(50)
	}
	return do.Select(p.keys), nil
}

// Find 执行搜索，提供 ids 时优先返回（UNION ALL），IDs 先经校验。
func (x *SearchDto) Find(ctx context.Context, do *gorm.DB, i any) (err error) {
	p, err := x.Get(ctx)
	if err != nil {
		return
	}
	if x.IDs != "" {
		ids := strings.Split(x.IDs, ",")
		// 校验每个 ID
		if p.idValidator != nil {
			for _, id := range ids {
				id = strings.TrimSpace(id)
				if id != "" && !p.idValidator(id) {
					return help.E(0, fmt.Sprintf(`IDs 中的 ID 格式无效: %s`, id))
				}
			}
		}
		factory, err := x.Factory(ctx, do.WithContext(ctx))
		if err != nil {
			return err
		}
		return do.Raw(`(?) union all (?)`,
			do.WithContext(ctx).Select(p.keys).Where(`id in (?)`, ids),
			factory.Where(`id not in (?)`, ids),
		).Find(i).Error
	}
	factory, err := x.Factory(ctx, do)
	if err != nil {
		return
	}
	return factory.Find(i).Error
}

// SearchResult 搜索响应结构。
type SearchResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DeleteDto 批量删除 DTO。
type DeleteDto struct {
	IDs []string `json:"ids"`
}
