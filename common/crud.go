package common

import (
	"fmt"
	"gorm.io/gorm/clause"
	"regexp"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/help"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type Controller interface {
	Create(ctx context.Context, c *app.RequestContext)
	Find(ctx context.Context, c *app.RequestContext)
	FindById(ctx context.Context, c *app.RequestContext)
	Update(ctx context.Context, c *app.RequestContext)
	Delete(ctx context.Context, c *app.RequestContext)
}

func SetPipe(ctx context.Context, i any) context.Context {
	return context.WithValue(ctx, "pipe", i)
}

var ToOrderBy = map[string]string{
	"1":  "",
	"-1": "desc",
}

type ExistsDto struct {
	// 需要查重的字段
	// @example ?keys=name 则等于 where name=$q
	Key string `query:"key,omitempty"`
	Q   string `query:"q,omitempty"` // 关键词查询
}

type ExistsPipe struct {
	fields map[string]bool
}

func NewExistsPipe(keys ...string) *ExistsPipe {
	fields := make(map[string]bool)
	for _, key := range keys {
		fields[key] = true
	}
	return &ExistsPipe{
		fields: fields,
	}
}

func (x *ExistsDto) Get(ctx context.Context) *ExistsPipe {
	return ctx.Value("pipe").(*ExistsPipe)
}

type ExistsResult struct {
	Exists bool `json:"exists"`
}

func (x *ExistsDto) Exists(ctx context.Context, do *gorm.DB) (result ExistsResult, err error) {
	p := x.Get(ctx)
	if !p.fields[x.Key] {
		err = help.E(0, fmt.Sprintf(`[%s] 字段不允许查重`, x.Key))
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

type FindDto struct {
	PageSize int64    `header:"x-pagesize" vd:"omitempty,min=0,max=1000"`
	Page     int64    `header:"x-page" vd:"omitempty,min=0"`
	Q        string   `query:"q,omitempty"`                             // 关键词查询
	Sort     []string `query:"sort,omitempty" vd:"omitempty,dive,sort"` // 排序
}

func (x *FindDto) GetPageSize() int {
	if x.PageSize == 0 {
		x.PageSize = 1000
	}
	return int(x.PageSize)
}

func (x *FindDto) GetOffset() int {
	return int(x.Page) * int(x.PageSize)
}

func (x *FindDto) IsCode() bool {
	regex, _ := regexp.Compile(`^N-`)
	return regex.MatchString(x.Q)
}

func (x *FindDto) IsNo() bool {
	regex, _ := regexp.Compile(`^[NBM]-`)
	return regex.MatchString(x.Q)
}

func (x *FindDto) GetCode() string {
	return x.Q[2:]
}

func (x *FindDto) GetKeyword() string {
	return fmt.Sprintf(`%%%s%%`, x.Q)
}

type FindPipe struct {
	ts   bool     // 存在默认时间，例如排除：create_time,update_time，使用 create_time 倒序
	sort bool     // 排序处理
	page bool     // 分页处理
	keys []string // 指定返回字段，不为空时 omit 失效
	omit []string // 排除返回字段
}

func (x *FindDto) Get(ctx context.Context) *FindPipe {
	return ctx.Value("pipe").(*FindPipe)
}

func NewFindPipe() *FindPipe {
	return &FindPipe{
		ts:   true,
		sort: true,
		page: true,
	}
}

func (x *FindPipe) SkipTs() *FindPipe {
	x.ts = false
	return x
}

func (x *FindPipe) SkipSort() *FindPipe {
	x.sort = false
	return x
}

func (x *FindPipe) SkipPage() *FindPipe {
	x.page = false
	return x
}

func (x *FindPipe) Select(keys ...string) *FindPipe {
	x.keys = keys
	return x
}

func (x *FindPipe) Omit(keys ...string) *FindPipe {
	x.omit = keys
	return x
}

func (x *FindDto) Factory(ctx context.Context, do *gorm.DB) *gorm.DB {
	p := x.Get(ctx)
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
			do = do.Order(fmt.Sprintf(`%s %s`, rule[0], ToOrderBy[rule[1]]))
		}
	}

	if p.page {
		do = do.Limit(x.GetPageSize()).Offset(x.GetOffset())
	}
	return do
}

func (x *FindDto) Find(ctx context.Context, do *gorm.DB, i interface{}) (err error) {
	if err = x.Factory(ctx, do).Find(i).Error; err != nil {
		return
	}
	return
}

type FindByIdDto struct {
	ID   string `path:"id"`
	Full int    `query:"full,omitempty"` // 全字段处理（即关联返回）
}

func (x *FindByIdDto) IsFull() bool {
	return x.Full == 1
}

type FindByIdPipe struct {
	ts    bool     // 存在默认时间，例如排除：create_time,update_time，使用 create_time 倒序
	keys  []string // 指定返回字段，不为空时 omit 失效
	omit  []string // 排除返回字段
	fKeys []string // 指定完整模式返回字段，不为空时 fOmit 失效
	fOmit []string // 指定完整模式排除字段
}

func (x *FindByIdDto) Get(ctx context.Context) *FindByIdPipe {
	return ctx.Value("pipe").(*FindByIdPipe)
}

func NewFindByIdPipe() *FindByIdPipe {
	return &FindByIdPipe{
		ts: true,
	}
}

func (x *FindByIdPipe) SkipTs() *FindByIdPipe {
	x.ts = false
	return x
}

func (x *FindByIdPipe) Select(keys ...string) *FindByIdPipe {
	x.keys = keys
	return x
}

func (x *FindByIdPipe) Omit(keys ...string) *FindByIdPipe {
	x.omit = keys
	return x
}

func (x *FindByIdPipe) FullSelect(keys ...string) *FindByIdPipe {
	x.fKeys = keys
	return x
}

func (x *FindByIdPipe) FullOmit(keys ...string) *FindByIdPipe {
	x.fOmit = keys
	return x
}

func (x *FindByIdDto) Take(ctx context.Context, do *gorm.DB, i interface{}) (err error) {
	p := x.Get(ctx)
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

	return do.Where(`id = ?`, x.ID).Take(&i).Error
}

type SearchDto struct {
	M   string `query:"m,omitempty"`   // 使用查询模式，1：简化字段，用于异步返回
	Q   string `query:"q,omitempty"`   // 关键词查询
	IDs string `query:"ids,omitempty"` // 已存在 IDs
}

func (x *SearchDto) IsCode() bool {
	regex, _ := regexp.Compile(`^N-`)
	return regex.MatchString(x.Q)
}

func (x *SearchDto) IsNo() bool {
	regex, _ := regexp.Compile(`^[NBM]-`)
	return regex.MatchString(x.Q)
}

func (x *SearchDto) GetCode() string {
	return x.Q[2:]
}

func (x *SearchDto) GetKeyword() string {
	return fmt.Sprintf(`%%%s%%`, x.Q)
}

type SearchPipe struct {
	keys  []string // 指定字段，默认：id,name
	async bool     // 前端异步
}

func (x *SearchPipe) SkipAsync() *SearchPipe {
	x.async = false
	return x
}

func NewSearchPipe(keys ...string) *SearchPipe {
	search := &SearchPipe{
		keys:  []string{},
		async: true,
	}
	for _, key := range keys {
		search.keys = append(search.keys, key)
	}
	if len(keys) == 0 {
		search.keys = []string{"id", "name"}
	}
	return search
}

func (x *SearchDto) Get(ctx context.Context) *SearchPipe {
	return ctx.Value("pipe").(*SearchPipe)
}

func (x *SearchDto) Factory(ctx context.Context, do *gorm.DB) *gorm.DB {
	p := x.Get(ctx)
	if p.async {
		do = do.Limit(50)
	}
	return do.Select(p.keys)
}

func (x *SearchDto) Find(ctx context.Context, do *gorm.DB, i any) (err error) {
	p := x.Get(ctx)
	if x.IDs != "" {
		ids := strings.Split(x.IDs, ",")
		return do.Raw(`(?) union all (?)`,
			do.WithContext(ctx).Select(p.keys).Where(`id in (?)`, ids),
			x.Factory(ctx, do.WithContext(ctx)).Where(`id not in (?)`, ids),
		).Find(i).Error
	} else {
		return x.Factory(ctx, do).Find(i).Error
	}
}

type SearchResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DeleteDto struct {
	IDs []string `json:"ids"`
}
