package pages

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Controller struct {
	Service *Service
}

func (x *Controller) Indexes(c *gin.Context) interface{} {
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
	indexes, err := x.Service.Indexes(ctx, data.Schema.Key)
	if err != nil {
		return err
	}
	return indexes
}

func (x *Controller) CreateIndex(c *gin.Context) interface{} {
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
	if _, err = x.Service.CreateIndex(ctx, page.Schema.Key, uri.Index, body.Keys, *body.Unique); err != nil {
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
	if _, err = x.Service.DeleteIndex(ctx, page.Schema.Key, uri.Index); err != nil {
		return err
	}
	return nil
}
