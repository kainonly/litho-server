package dsl

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/weplanx/server/utils/passlib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	DslService *Service
}

func (x *Controller) In(r *route.RouterGroup) {
	r.POST("", x.Create)
	r.POST("bulk-create", x.BulkCreate)
	r.GET("_size", x.Size)
	r.GET("", x.Find)
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
		Format M `json:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	// 数据转换
	if err := x.Transform(dto.Data, dto.Format); err != nil {
		c.Error(errors.New(err, errors.ErrorTypePublic, nil))
		return
	}
	dto.Data["create_time"] = time.Now()
	dto.Data["update_time"] = time.Now()

	any, err := x.DslService.Create(ctx, dto.Model, dto.Data)
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
		Format M `json:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	// 数据转换
	docs := make([]interface{}, len(dto.Data))
	for i, doc := range dto.Data {
		if err := x.Transform(doc, dto.Format); err != nil {
			c.Error(errors.New(err, errors.ErrorTypePublic, nil))
			return
		}
		doc["create_time"] = time.Now()
		doc["update_time"] = time.Now()
		docs[i] = doc
	}

	any, err := x.DslService.BulkCreate(ctx, dto.Model, docs)
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
		Format M `query:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	// 数据转换
	if err := x.Transform(dto.Filter, dto.Format); err != nil {
		c.Error(errors.New(err, errors.ErrorTypePublic, nil))
		return
	}

	size, err := x.DslService.Size(ctx, dto.Model, dto.Filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-total", strconv.Itoa(int(size)))
	c.Status(http.StatusNoContent)
}

// Find 获取匹配文档
// @router /dsl/:model [GET]
func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model,required" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 分页数量
		Pagesize int64 `header:"x-pagesize" vd:"$>=0 && $<=1000;msg:'分页数量必须在 1~1000 之间'"`
		// 页码
		Page int64 `header:"x-page" vd:"$>=0;msg:'页码必须大于 0'"`
		// 筛选条件
		Filter M `query:"filter"`
		// 筛选条件格式转换
		Format M `query:"format"`
		// 排序规则
		Sort M `query:"sort" vd:"range($,in(#v,-1,1));msg:'排序规则不规范'"`
		// 投影规则
		Keys M `query:"keys" vd:"range($,in(#v,0,1));msg:'投影规则不规范'"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	// 数据转换
	if err := x.Transform(dto.Filter, dto.Format); err != nil {
		c.Error(errors.New(err, errors.ErrorTypePublic, nil))
		return
	}

	size, err := x.DslService.Size(ctx, dto.Model, dto.Filter)
	if err != nil {
		c.Error(err)
		return
	}

	// 默认分页数量 100
	if dto.Pagesize == 0 {
		dto.Pagesize = 100
	}

	// 默认页码 1
	if dto.Page == 0 {
		dto.Page = 1
	}

	var sort bson.D
	for key, value := range dto.Sort {
		sort = append(sort, bson.E{Key: key, Value: value})
	}
	// 默认倒序 ID
	if len(sort) == 0 {
		sort = bson.D{{Key: "_id", Value: -1}}
	}

	option := options.Find().
		SetLimit(dto.Pagesize).
		SetSkip((dto.Page - 1) * dto.Pagesize).
		SetProjection(dto.Keys).
		SetSort(sort).
		SetAllowDiskUse(true)
	data, err := x.DslService.Find(ctx, dto.Model, dto.Filter, option)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-total", strconv.Itoa(int(size)))
	c.JSON(http.StatusOK, data)
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
		// 筛选条件格式转换
		Format M `query:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	// 数据转换
	if err := x.Transform(dto.Filter, dto.Format); err != nil {
		c.Error(errors.New(err, errors.ErrorTypePublic, nil))
		return
	}

	option := options.FindOne().
		SetProjection(dto.Keys)
	data, err := x.DslService.FindOne(ctx, dto.Model, dto.Filter, option)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
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

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	option := options.FindOne().
		SetProjection(dto.Keys)
	data, err := x.DslService.FindOne(ctx, dto.Model, M{"_id": id}, option)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

// Update 局部更新匹配文档
// @router /dsl/:model [PATCH]
func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 模型命名
		Model string `path:"model" vd:"regexp('^[a-z_]+$');msg:'模型名称必须是小写字母与下划线'"`
		// 筛选条件
		Filter M `query:"filter,required" vd:"len($)>0;msg:'筛选条件不能为空'"`
		// 筛选条件格式转换
		FFormat M `query:"format"`
		// 更新数据
		Data M `json:"data,required" vd:"len($)>0;msg:'更新数据不能为空'"`
		// 文档字段格式转换
		DFormat M `json:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	// 数据转换
	if err := x.Transform(dto.Filter, dto.FFormat); err != nil {
		c.Error(errors.New(err, errors.ErrorTypePublic, nil))
		return
	}
	if err := x.Transform(dto.Data, dto.DFormat); err != nil {
		c.Error(errors.New(err, errors.ErrorTypePublic, nil))
		return
	}
	if _, ok := dto.Data["$set"]; !ok {
		dto.Data["$set"] = M{}
	}
	dto.Data["$set"].(M)["update_time"] = time.Now()

	any, err := x.DslService.Update(ctx, dto.Model, dto.Filter, dto.Data)
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
		Format M `json:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	// 数据转换
	if err := x.Transform(dto.Data, dto.Format); err != nil {
		c.Error(errors.New(err, errors.ErrorTypePublic, nil))
		return
	}
	if _, ok := dto.Data["$set"]; !ok {
		dto.Data["$set"] = M{}
	}
	dto.Data["$set"].(M)["update_time"] = time.Now()

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	any, err := x.DslService.UpdateById(ctx, dto.Model, id, dto.Data)
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
		Format M `json:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	// 数据转换
	if err := x.Transform(dto.Data, dto.Format); err != nil {
		c.Error(errors.New(err, errors.ErrorTypePublic, nil))
		return
	}
	dto.Data["update_time"] = time.Now()

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	any, err := x.DslService.Replace(ctx, dto.Model, id, dto.Data)
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

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	any, err := x.DslService.Delete(ctx, dto.Model, id)
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
		Format M `json:"format"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	// 数据转换
	if err := x.Transform(dto.Data, dto.Format); err != nil {
		c.Error(errors.New(err, errors.ErrorTypePublic, nil))
		return
	}

	any, err := x.DslService.BulkDelete(ctx, dto.Model, dto.Data)
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

// Transform 格式转换
func (x *Controller) Transform(data M, format M) (err error) {
	for path, spec := range format {
		keys, cursor := strings.Split(path, "."), data
		n := len(keys) - 1
		for _, key := range keys[:n] {
			if v, ok := cursor[key].(M); ok {
				cursor = v
			}
		}
		key := keys[n]
		if cursor[key] == nil {
			continue
		}
		switch spec {
		case "oid":
			// 转换为 ObjectId
			if cursor[key], err = primitive.ObjectIDFromHex(cursor[key].(string)); err != nil {
				return
			}
			break

		case "oids":
			// 转换为 ObjectId 数组
			oids := cursor[key].([]interface{})
			for i, id := range oids {
				if oids[i], err = primitive.ObjectIDFromHex(id.(string)); err != nil {
					return
				}
			}
			break
		case "date":
			// 转换为 ISODate
			if cursor[key], err = time.Parse(time.RFC1123, cursor[key].(string)); err != nil {
				return
			}
			break

		case "password":
			// 密码类型，转换为 Argon2id
			if cursor[key], err = passlib.Hash(cursor[key].(string)); err != nil {
				return
			}
			break
		}
	}
	return
}
