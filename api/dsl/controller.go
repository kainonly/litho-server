package dsl

import (
	"github.com/gin-gonic/gin"
	"github.com/weplanx/server/utils/helper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

type Controller struct {
	DslService *Service
}

func (x *Controller) In(r *gin.RouterGroup) {
	r.POST("", x.Create)
	r.POST("bulk-create", x.BulkCreate)
	r.HEAD("_size", x.Size)
	r.HEAD("_exists", x.Exists)
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
func (x *Controller) Create(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var header struct {
		// 文档字段格式转换
		Doc string `header:"wpx-doc"`
	}
	if err := c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		return
	}
	xdoc := helper.ParseArray(header.Doc)

	var body M
	if err := helper.BindAndValidate(c.Request.Body, &body, `required,gt=0`); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	any, err := x.DslService.Create(ctx, params.Model, body, xdoc)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, any)
}

// BulkCreate 批量创建文档
func (x *Controller) BulkCreate(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var header struct {
		// 文档字段格式转换
		Doc string `header:"wpx-doc"`
	}
	if err := c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		return
	}
	xdoc := helper.ParseArray(header.Doc)

	var body []M
	if err := helper.BindAndValidate(c.Request.Body, &body, `required,gt=0,dive,gt=0`); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	any, err := x.DslService.BulkCreate(ctx, params.Model, body, xdoc)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, any)
}

// Size 获取文档总数
func (x *Controller) Size(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var header struct {
		// 筛选字段格式转换
		Filter string `header:"wpx-filter"`
	}
	if err := c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		return
	}
	xfilter := helper.ParseArray(header.Filter)

	var query struct {
		// 筛选字段
		Filter M `form:"filter"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	size, err := x.DslService.Size(ctx, params.Model, query.Filter, xfilter)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("wpx-total", strconv.Itoa(int(size)))
	c.Status(http.StatusNoContent)
}

// Exists 获取文档存在状态
func (x *Controller) Exists(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var header struct {
		// 筛选字段格式转换
		Filter string `header:"wpx-filter"`
	}
	if err := c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		return
	}
	xfilter := helper.ParseArray(header.Filter)

	var query struct {
		// 筛选字段
		Filter M `form:"filter" binding:"required,gt=0"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	size, err := x.DslService.Size(ctx, params.Model, query.Filter, xfilter)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("wpx-exists", strconv.FormatBool(size != 0))
	c.Status(http.StatusNoContent)
}

// Find 获取匹配文档
func (x *Controller) Find(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var header struct {
		// 筛选字段格式转换
		Filter string `header:"wpx-filter"`
	}
	if err := c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		return
	}
	xfilter := helper.ParseArray(header.Filter)

	var query struct {
		// 筛选字段
		Filter M `form:"filter"`
		// 排序规则
		Sort M `form:"sort" binding:"omitempty,gt=0"`
		// 投影规则
		Keys M `form:"keys" binding:"omitempty,gt=0"`
		// 最大返回数量
		Limit int64 `form:"limit" binding:"omitempty,max=1000,min=1"`
		// 跳过数量
		Skip int64 `form:"skip" binding:"omitempty,min=0"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	data, err := x.DslService.Find(ctx, params.Model, query.Filter, xfilter, FindOption{
		Sort:  query.Sort,
		Keys:  query.Keys,
		Limit: query.Limit,
		Skip:  query.Skip,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

// FindPages 获取匹配分页文档
func (x *Controller) FindPages(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var header struct {
		// 筛选字段格式转换
		Filter string `header:"wpx-filter"`
	}
	if err := c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		return
	}
	xfilter := helper.ParseArray(header.Filter)

	var query struct {
		// 筛选字段
		Filter M `form:"filter"`
		// 排序规则
		Sort M `form:"sort" binding:"omitempty,gt=0"`
		// 投影规则
		Keys M `form:"keys" binding:"omitempty,gt=0"`
		// 最大返回数量
		Limit int64 `form:"limit" binding:"omitempty,max=1000,min=1"`
		// 分页页码
		Page int64 `form:"page" binding:"omitempty,min=1"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(err)
		return
	}
	if query.Page == 0 {
		query.Page = 1
	}

	ctx := c.Request.Context()
	size, err := x.DslService.Size(ctx, params.Model, query.Filter, xfilter)
	if err != nil {
		c.Error(err)
		return
	}
	c.Header("wpx-total", strconv.Itoa(int(size)))

	data, err := x.DslService.Find(ctx, params.Model, query.Filter, xfilter, FindOption{
		Sort:  query.Sort,
		Keys:  query.Keys,
		Limit: query.Limit,
		Page:  query.Page,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

// FindOne 获取单个文档
func (x *Controller) FindOne(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var header struct {
		// 筛选字段格式转换
		Filter string `header:"wpx-filter"`
	}
	if err := c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		return
	}
	xfilter := helper.ParseArray(header.Filter)

	var query struct {
		// 筛选字段
		Filter M `form:"filter" binding:"required,gt=0"`
		// 投影规则
		Keys M `form:"keys" binding:"omitempty,gt=0"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	data, err := x.DslService.FindOne(ctx, params.Model, query.Filter, xfilter, FindOption{
		Keys: query.Keys,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

// FindById 获取指定 Id 的文档
func (x *Controller) FindById(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
		// ObjectId
		Id string `uri:"id" binding:"required,objectId"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var query struct {
		// 投影规则
		Keys M `form:"keys" binding:"omitempty,gt=0"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	data, err := x.DslService.FindById(ctx, params.Model, params.Id, FindOption{
		Keys: query.Keys,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

// Update 局部更新多个匹配文档
func (x *Controller) Update(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var header struct {
		// 筛选字段格式转换
		Filter string `header:"wpx-filter"`
		// 文档字段格式转换
		Doc string `header:"wpx-doc"`
	}
	if err := c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		return
	}
	xfilter := helper.ParseArray(header.Filter)
	xdoc := helper.ParseArray(header.Doc)

	var query struct {
		// 筛选字段
		Filter M `form:"filter" binding:"required,gt=0"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(err)
		return
	}

	var body M
	if err := helper.BindAndValidate(c.Request.Body, &body, `required,gt=0`); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	any, err := x.DslService.Update(ctx, params.Model, query.Filter, xfilter, body, xdoc)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// UpdateById 局部更新指定 Id 的文档
func (x *Controller) UpdateById(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
		// ObjectId
		Id string `uri:"id" binding:"required,objectId"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var header struct {
		// 文档字段格式转换
		Doc string `header:"wpx-doc"`
	}
	if err := c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		return
	}
	xdoc := helper.ParseArray(header.Doc)

	var body M
	if err := helper.BindAndValidate(c.Request.Body, &body, `required,gt=0`); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	any, err := x.DslService.UpdateById(ctx, params.Model, params.Id, body, xdoc)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// Replace 替换指定 Id 的文档
func (x *Controller) Replace(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
		// ObjectId
		Id string `uri:"id" binding:"required,objectId"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var header struct {
		// 文档字段格式转换
		Doc string `header:"wpx-doc"`
	}
	if err := c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		return
	}
	xdoc := helper.ParseArray(header.Doc)

	var body M
	if err := helper.BindAndValidate(c.Request.Body, &body, `required,gt=0`); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	any, err := x.DslService.Replace(ctx, params.Model, params.Id, body, xdoc)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// Delete 删除指定 Id 的文档
func (x *Controller) Delete(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
		// ObjectId
		Id string `uri:"id" binding:"required,objectId"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	any, err := x.DslService.Delete(ctx, params.Model, params.Id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// BulkDelete 批量删除匹配文档
func (x *Controller) BulkDelete(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var header struct {
		// 筛选字段格式转换
		Filter string `header:"wpx-filter"`
		// 文档字段格式转换
		Doc string `header:"wpx-doc"`
	}
	if err := c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		return
	}
	xfilter := helper.ParseArray(header.Filter)

	var body M
	if err := helper.BindAndValidate(c.Request.Body, &body, `required,gt=0`); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	any, err := x.DslService.BulkDelete(ctx, params.Model, body, xfilter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}

// Sort 通用排序
func (x *Controller) Sort(c *gin.Context) {
	var params struct {
		// 模型命名
		Model string `uri:"model" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var body []primitive.ObjectID
	if err := helper.BindAndValidate(c.Request.Body, &body, `required,gt=0,dive,gt=0`); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	any, err := x.DslService.Sort(ctx, params.Model, body)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, any)
}
