package pages

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Controller struct {
	Service *Service
}

// Navs 页面导航
func (x *Controller) Navs(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	navs, err := x.Service.Navs(ctx)
	if err != nil {
		return err
	}
	return navs
}

func (x *Controller) Dynamic(c *gin.Context) interface{} {
	var uri struct {
		Id string `uri:"id" binding:"required,objectId"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		return err
	}
	ctx := c.Request.Context()
	data, err := x.Service.FindOneById(ctx, uri.Id)
	if err != nil {
		return err
	}
	return data
}

func (x *Controller) GetIndexes(c *gin.Context) interface{} {
	var uri struct {
		Id string `uri:"id" binding:"required,objectId"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		return err
	}
	ctx := c.Request.Context()
	data, err := x.Service.FindOneById(ctx, uri.Id)
	if err != nil {
		return err
	}
	indexes, err := x.Service.GetIndexes(ctx, data.Schema.Model)
	if err != nil {
		return err
	}
	return indexes
}

func (x *Controller) SetIndex(c *gin.Context) interface{} {
	var uri struct {
		Id    string `uri:"id" binding:"required,objectId"`
		Index string `uri:"index" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		return err
	}
	var body struct {
		Keys   bson.D `json:"keys" binding:"required,gt=0"`
		Unique *bool  `json:"unique" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	page, err := x.Service.FindOneById(ctx, uri.Id)
	if err != nil {
		return err
	}
	if _, err = x.Service.SetIndex(ctx, page.Schema.Model, uri.Index, body.Keys, *body.Unique); err != nil {
		return err
	}
	return nil
}

func (x *Controller) DeleteIndex(c *gin.Context) interface{} {
	var uri struct {
		Id    string `uri:"id" binding:"required,objectId"`
		Index string `uri:"index" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		return err
	}
	ctx := c.Request.Context()
	page, err := x.Service.FindOneById(ctx, uri.Id)
	if err != nil {
		return err
	}
	if _, err = x.Service.DeleteIndex(ctx, page.Schema.Model, uri.Index); err != nil {
		return err
	}
	return nil
}
