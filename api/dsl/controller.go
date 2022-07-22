package dsl

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/route"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Controller struct {
	DslService *Service
}

func (x *Controller) In(r *route.RouterGroup) {
	r.POST("", x.Create)
	r.POST("bulk-create", x.BulkCreate)
	r.GET("_size", x.Size)
	r.GET("_exists", x.Exists)
	r.GET("", x.Find)
	r.GET("_pages", x.FindPages)
	r.GET("_one", x.FindOne)
	r.GET(":id", x.FindById)
	r.PATCH("", x.Update)
	r.PATCH(":id", x.UpdateById)
	r.PUT(":id", x.Replace)
	r.DELETE(":id", x.Delete)
	r.POST("bulk-delete", x.BulkDelete)
	r.POST("sort", x.Sort)
}

// Create 创建文档
// @router /dsl/:model [POST]
func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 文档数据
		Data M `json:"data,required" vd:"len($)>0;msg:'文档不能为空数据'"`
		// 文档字段格式转换
		Format []string `json:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.Create(ctx, dto.Model, dto.Data, dto.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, any)
}

// BulkCreate 批量创建文档
// @router /dsl/:model/bulk-create [POST]
func (x *Controller) BulkCreate(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model,required" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 批量文档数据
		Data []M `json:"data,required" vd:"len($)>0 && range($,len(#v)>0);msg:'批量文档不能存在空数据'"`
		// 文档字段格式转换
		Format []string `json:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.BulkCreate(ctx, dto.Model, dto.Data, dto.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, any)
}

// Size 获取文档总数
// @router /dsl/:model/_size [GET]
func (x *Controller) Size(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model,required" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 筛选条件
		Filter M `query:"filter"`
		// 筛选条件格式转换
		Format []string `query:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	size, err := x.DslService.Size(ctx, dto.Model, dto.Filter, dto.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{"total": size})
}

// Exists 获取文档存在状态
// @router /dsl/:model/_exists [GET]
func (x *Controller) Exists(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model,required" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 筛选条件
		Filter M `query:"filter,required" vd:"len($)>0;msg:'筛选条件不能为空'"`
		// 筛选字段格式转换
		Format []string `query:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	size, err := x.DslService.Size(ctx, dto.Model, dto.Filter, dto.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{"exists": size != 0})
}

// Find 获取匹配文档
// @router /dsl/:model [GET]
func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model,required" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 筛选字段
		Filter M `query:"filter"`
		// 排序规则
		Sort M `query:"sort" vd:"range($,in(#v,-1,1));msg:'排序规则不规范'"`
		// 投影规则
		Keys M `query:"keys" vd:"range($,in(#v,0,1));msg:'投影规则不规范'"`
		// 最大返回数量
		Limit int64 `query:"limit" vd:"$>=0 && $<=1000;msg:'最大返回数量必须在 1~1000 之间'"`
		// 跳过数量
		Skip int64 `query:"skip" vd:"$>=0;msg:'跳过数量必须大于 0'"`
		// 筛选字段格式转换
		Format []string `query:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	data, err := x.DslService.Find(ctx, dto.Model, dto.Filter, dto.Format, FindOption{
		Sort:  dto.Sort,
		Keys:  dto.Keys,
		Limit: dto.Limit,
		Skip:  dto.Skip,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

// FindPages 获取匹配分页文档
// @router /dsl/:model/_pages [GET]
func (x *Controller) FindPages(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model,required" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 筛选条件
		Filter M `query:"filter"`
		// 排序规则
		Sort M `query:"sort" vd:"range($,in(#v,-1,1));msg:'排序规则不规范'"`
		// 投影规则
		Keys M `query:"keys" vd:"range($,in(#v,0,1));msg:'投影规则不规范'"`
		// 最大返回数量
		Limit int64 `query:"limit" vd:"$>=0 && $<=1000;msg:'最大返回数量必须在 0~1000 之间'"`
		// 分页页码
		Page int64 `query:"page" vd:"$>=0;msg:'分页页码必须大于 0'"`
		// 筛选条件格式转换
		Format []string `query:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	size, err := x.DslService.Size(ctx, dto.Model, dto.Filter, dto.Format)
	if err != nil {
		c.Error(err)
		return
	}

	data, err := x.DslService.Find(ctx, dto.Model, dto.Filter, dto.Format, FindOption{
		Sort:  dto.Sort,
		Keys:  dto.Keys,
		Limit: dto.Limit,
		Page:  dto.Page,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"total": size,
		"data":  data,
	})
}

// FindOne 获取单个文档
// @router /dsl/:model/_one [GET]
func (x *Controller) FindOne(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model,required" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 筛选条件
		Filter M `query:"filter,required" vd:"len($)>0;msg:'筛选条件不能为空'"`
		// 投影规则
		Keys M `query:"keys" vd:"range($,in(#v,0,1));msg:'投影规则不规范'"`
		// 筛选字段格式转换
		Format []string `query:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	data, err := x.DslService.FindOne(ctx, dto.Model, dto.Filter, dto.Format, FindOption{
		Keys: dto.Keys,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"data": data,
	})
}

// FindById 获取指定 Id 的文档
// @router /dsl/:model/:id [GET]
func (x *Controller) FindById(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model,required" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 文档 ID
		Id string `path:"id,required" vd:"mongoId($);msg:'文档 ID 不规范'"`
		// 投影规则
		Keys M `query:"keys" vd:"range($,in(#v,0,1));msg:'投影规则不规范'"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	data, err := x.DslService.FindById(ctx, dto.Model, dto.Id, FindOption{
		Keys: dto.Keys,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"data": data,
	})
}

// Update 局部更新多个匹配文档
// @router /dsl/:model [PATCH]
func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 筛选条件
		Filter M `query:"filter,required" vd:"len($)>0;msg:'筛选条件不能为空'"`
		// 筛选条件格式转换
		FFormat []string `query:"format"`
		// 更新数据
		Data M `json:"data,required" vd:"len($)>0;msg:'更新数据不能为空'"`
		// 文档字段格式转换
		DFormat []string `json:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.Update(ctx, dto.Model, dto.Filter, dto.FFormat, dto.Data, dto.DFormat)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// UpdateById 局部更新指定 Id 的文档
// @router /dsl/:model/:id [PATCH]
func (x *Controller) UpdateById(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model,required" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 文档 ID
		Id string `path:"id,required" vd:"mongoId($);msg:'文档 ID 不规范'"`
		// 更新数据
		Data M `json:"data,required" vd:"len($)>0;msg:'更新数据不能为空'"`
		// 文档字段格式转换
		Format []string `json:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.UpdateById(ctx, dto.Model, dto.Id, dto.Data, dto.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// Replace 替换指定 Id 的文档
// @router /dsl/:model/:id [PUT]
func (x *Controller) Replace(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model,required" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 文档 ID
		Id string `path:"id,required" vd:"mongoId($);msg:'文档 ID 不规范'"`
		// 文档数据
		Data M `json:"data,required" vd:"len($)>0;msg:'文档数据不能为空'"`
		// 文档字段格式转换
		Format []string `json:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.Replace(ctx, dto.Model, dto.Id, dto.Data, dto.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// Delete 删除指定 Id 的文档
// @router /dsl/:model/:id [DELETE]
func (x *Controller) Delete(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model,required" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 文档 ID
		Id string `path:"id,required" vd:"mongoId($);msg:'文档 ID 不规范'"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.Delete(ctx, dto.Model, dto.Id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// BulkDelete 批量删除匹配文档
// @router /dsl/:model/bulk-delete [POST]
func (x *Controller) BulkDelete(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model,required" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 筛选条件
		Data M `json:"data,required" vd:"len($)>0;msg:'筛选条件不能为空'"`
		// 筛选条件格式转换
		Format []string `json:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.BulkDelete(ctx, dto.Model, dto.Data, dto.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// Sort 通用排序
// @router /dsl/:model/sort [POST]
func (x *Controller) Sort(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model,required" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 文档 ID 数组
		Data []primitive.ObjectID `json:"data,required" vd:"len($)>0 && range($,mongoId(#v));msg:'数组必须均为文档 ID'"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.Sort(ctx, dto.Model, dto.Data)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}
