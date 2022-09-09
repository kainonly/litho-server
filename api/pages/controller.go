package pages

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Controller struct {
	PagesService *Service
}

// GetNavs 导航数据
// @router /navs [GET]
func (x *Controller) GetNavs(ctx context.Context, c *app.RequestContext) {
	active := common.GetActive(c)

	data, err := x.PagesService.GetNavs(ctx, active.UID)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetOne 获取页面数据
// @router /:id
func (x *Controller) GetOne(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 页面 ID
		Id string `path:"id,required" vd:"mongoId($);msg:'页面 ID 不规范'"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	page, err := x.PagesService.FindOneById(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, page)
}

// GetIndexes 获取页面的模型索引
// @router /:id/indexes [GET]
func (x *Controller) GetIndexes(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 页面 ID
		Id string `path:"id,required" vd:"mongoId($);msg:'页面 ID 不规范'"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	page, err := x.PagesService.FindOneById(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}
	indexes, err := x.PagesService.GetIndexes(ctx, page.Schema.Key)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, indexes)
}

// SetIndex 所属页面的模型设置索引
// @router /:id/indexes/:index [PUT]
func (x *Controller) SetIndex(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 页面 ID
		Id string `path:"id,required" vd:"mongoId($);msg:'页面 ID 不规范'"`
		// 索引名称
		Index string `path:"index,required" vd:"regexp('^[a-z_]+$');msg:'索引名称必须是小写字母与下划线'"`
		// 索引设置
		Keys bson.D `json:"keys" vd:"len($)>0;msg:'索引设置不能为空'"`
		// 唯一索引
		Unique *bool `json:"unique,omitempty"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	page, err := x.PagesService.FindOneById(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}
	unique := false
	if dto.Unique != nil {
		unique = *dto.Unique
	}
	if _, err = x.PagesService.SetIndex(ctx, page.Schema.Key, dto.Index, dto.Keys, unique); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteIndex 所属页面的模型删除索引
// @router /:id/indexes/:index [DELETE]
func (x *Controller) DeleteIndex(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 页面 ID
		Id string `path:"id,required" vd:"mongoId($);msg:'页面 ID 不规范'"`
		// 索引名称
		Index string `path:"index,required" vd:"regexp('^[a-z_]+$');msg:'索引名称必须是小写字母与下划线'"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	page, err := x.PagesService.FindOneById(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}
	if _, err = x.PagesService.DeleteIndex(ctx, page.Schema.Key, dto.Index); err != nil {
		return
	}

	c.Status(http.StatusNoContent)
}
