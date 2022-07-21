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
func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Params struct {
			// 模型命名
			Model string `path:"model" binding:"required,key"`
		}

		Body struct {
			// 数据
			Data M `json:"data,required" vd:"len($)>0"`
			// 文档字段格式转换
			Format []string `json:"format"`
		}
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.Create(ctx, dto.Params.Model, dto.Body.Data, dto.Body.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, any)
}

// BulkCreate 批量创建文档
func (x *Controller) BulkCreate(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Params struct {
			// 模型命名
			Model string `path:"model" binding:"required,key"`
		}

		Body struct {
			// 数据
			Data []M `json:"data,required" vd:"len($)>0 && range($, len(#v)>0)"`
			// 文档字段格式转换
			Format []string `json:"format"`
		}
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.BulkCreate(ctx, dto.Params.Model, dto.Body.Data, dto.Body.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, any)
}

// Size 获取文档总数
func (x *Controller) Size(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Params struct {
			// 模型命名
			Model string `path:"model" binding:"required,key"`
		}

		Query struct {
			// 筛选字段
			Filter M `query:"filter,required"`
			// 筛选字段格式转换
			Format []string `query:"format"`
		}
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	size, err := x.DslService.Size(ctx, dto.Params.Model, dto.Query.Filter, dto.Query.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{"total": size})
}

// Exists 获取文档存在状态
func (x *Controller) Exists(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Params struct {
			// 模型命名
			Model string `path:"model" binding:"required,key"`
		}
		Query struct {
			// 筛选字段
			Filter M `query:"filter,required" vd:"len($)>0"`
			// 筛选字段格式转换
			Format []string `query:"format"`
		}
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	size, err := x.DslService.Size(ctx, dto.Params.Model, dto.Query.Filter, dto.Query.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{"exists": size != 0})
}

// Find 获取匹配文档
func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Params struct {
			// 模型命名
			Model string `path:"model" binding:"required,key"`
		}
		Query struct {
			// 筛选字段
			Filter M `query:"filter"`
			// 排序规则
			Sort M `query:"sort" binding:"omitempty,gt=0"`
			// 投影规则
			Keys M `query:"keys" binding:"omitempty,gt=0"`
			// 最大返回数量
			Limit int64 `query:"limit" binding:"omitempty,max=1000,min=1"`
			// 跳过数量
			Skip int64 `query:"skip" binding:"omitempty,min=0"`
			// 筛选字段格式转换
			Format []string `query:"format"`
		}
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	data, err := x.DslService.Find(ctx, dto.Params.Model, dto.Query.Filter, dto.Query.Format, FindOption{
		Sort:  dto.Query.Sort,
		Keys:  dto.Query.Keys,
		Limit: dto.Query.Limit,
		Skip:  dto.Query.Skip,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

// FindPages 获取匹配分页文档
func (x *Controller) FindPages(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Params struct {
			// 模型命名
			Model string `path:"model" binding:"required,key"`
		}
		Query struct {
			// 筛选字段
			Filter M `query:"filter"`
			// 排序规则
			Sort M `query:"sort" binding:"omitempty,gt=0"`
			// 投影规则
			Keys M `query:"keys" binding:"omitempty,gt=0"`
			// 最大返回数量
			Limit int64 `query:"limit" binding:"omitempty,max=1000,min=1"`
			// 分页页码
			Page int64 `query:"page" binding:"omitempty,min=1"`
			// 筛选字段格式转换
			Format []string `query:"format"`
		}
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	size, err := x.DslService.Size(ctx, dto.Params.Model, dto.Query.Filter, dto.Query.Format)
	if err != nil {
		c.Error(err)
		return
	}

	data, err := x.DslService.Find(ctx, dto.Params.Model, dto.Query.Filter, dto.Query.Format, FindOption{
		Sort:  dto.Query.Sort,
		Keys:  dto.Query.Keys,
		Limit: dto.Query.Limit,
		Page:  dto.Query.Page,
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
func (x *Controller) FindOne(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Params struct {
			// 模型命名
			Model string `path:"model" binding:"required,key"`
		}
		Query struct {
			// 筛选字段
			Filter M `form:"filter" binding:"required,gt=0"`
			// 投影规则
			Keys M `form:"keys" binding:"omitempty,gt=0"`
			// 筛选字段格式转换
			Format []string `query:"format"`
		}
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	data, err := x.DslService.FindOne(ctx, dto.Params.Model, dto.Query.Filter, dto.Query.Format, FindOption{
		Keys: dto.Query.Keys,
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
func (x *Controller) FindById(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Params struct {
			// 模型命名
			Model string `path:"model" binding:"required,key"`
			// ObjectId
			Id string `path:"id" binding:"required,objectId"`
		}
		Query struct {
			// 投影规则
			Keys M `query:"keys" binding:"omitempty,gt=0"`
		}
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	data, err := x.DslService.FindById(ctx, dto.Params.Model, dto.Params.Id, FindOption{
		Keys: dto.Query.Keys,
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
func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Params struct {
			// 模型命名
			Model string `path:"model" binding:"required,key"`
		}
		Query struct {
			// 筛选字段
			Filter M `query:"filter" binding:"required,gt=0"`
			// 筛选字段格式转换
			Format []string `query:"format"`
		}
		Body struct {
			Data M `json:"data,required" vd:"len($)>0"`
			// 文档字段格式转换
			Format []string `json:"format"`
		}
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.Update(ctx,
		dto.Params.Model, dto.Query.Filter, dto.Query.Format, dto.Body.Data, dto.Body.Format,
	)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// UpdateById 局部更新指定 Id 的文档
func (x *Controller) UpdateById(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Params struct {
			// 模型命名
			Model string `path:"model" binding:"required,key"`
			// ObjectId
			Id string `path:"id" binding:"required,objectId"`
		}
		Body struct {
			Data M `json:"data,required" vd:"len($)>0"`
			// 文档字段格式转换
			Format []string `json:"format"`
		}
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.UpdateById(ctx, dto.Params.Model, dto.Params.Id, dto.Body.Data, dto.Body.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// Replace 替换指定 Id 的文档
func (x *Controller) Replace(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Params struct {
			// 模型命名
			Model string `path:"model" binding:"required,key"`
			// ObjectId
			Id string `path:"id" binding:"required,objectId"`
		}
		Body struct {
			Data M `json:"data,required" vd:"len($)>0"`
			// 文档字段格式转换
			Format []string `json:"format"`
		}
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.Replace(ctx, dto.Params.Model, dto.Params.Id, dto.Body.Data, dto.Body.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// Delete 删除指定 Id 的文档
func (x *Controller) Delete(ctx context.Context, c *app.RequestContext) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
		// ObjectId
		Id string `uri:"id" binding:"required,objectId"`
	}
	if err := c.BindAndValidate(&params); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.Delete(ctx, params.Model, params.Id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// BulkDelete 批量删除匹配文档
func (x *Controller) BulkDelete(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Params struct {
			// 模型命名
			Model string `path:"model" binding:"required,key"`
		}
		Body struct {
			// 筛选字段
			Data M `json:"filter" binding:"required,gt=0"`
			// 筛选字段格式转换
			Format []string `json:"format"`
		}
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.BulkDelete(ctx, dto.Params.Model, dto.Body.Data, dto.Body.Format)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// Sort 通用排序
func (x *Controller) Sort(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Params struct {
			// 模型命名
			Model string `uri:"model" binding:"required,key"`
		}
		Body struct {
			Data []primitive.ObjectID `json:"data,required" vd:"len($)>0"`
		}
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	any, err := x.DslService.Sort(ctx, dto.Params.Model, dto.Body.Data)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}
