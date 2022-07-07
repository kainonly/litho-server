package pages

import (
	"github.com/gin-gonic/gin"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"server/app/roles"
	"server/app/users"
	"server/common"
	"server/model"
)

type Controller struct {
	Pages    *Service
	Users    *users.Service
	Roles    *roles.Service
	Passport *passport.Passport
}

// Navs 页面导航
func (x *Controller) Navs(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	claims, err := x.Passport.GetClaims(c, common.TokenClaimsKey)
	if err != nil {
		return err
	}
	uid := claims.Context["uid"].(string)
	id, _ := primitive.ObjectIDFromHex(uid)
	var user model.User
	if err = x.Users.FindOneById(ctx, id, &user); err != nil {
		return err
	}
	var rolesList []model.Role
	if err = x.Roles.FindByIds(ctx, user.Roles, &rolesList); err != nil {
		return err
	}
	navs, err := x.Pages.Navs(ctx, rolesList)
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
	id, _ := primitive.ObjectIDFromHex(uri.Id)
	var page model.Page
	if err := x.Pages.FindOneById(ctx, id, &page); err != nil {
		return err
	}
	return page
}

func (x *Controller) GetIndexes(c *gin.Context) interface{} {
	var uri struct {
		Id string `uri:"id" binding:"required,objectId"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		return err
	}
	ctx := c.Request.Context()
	id, _ := primitive.ObjectIDFromHex(uri.Id)
	var page model.Page
	if err := x.Pages.FindOneById(ctx, id, &page); err != nil {
		return err
	}
	indexes, err := x.Pages.GetIndexes(ctx, page.Schema.Key)
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
	id, _ := primitive.ObjectIDFromHex(uri.Id)
	var page model.Page
	if err := x.Pages.FindOneById(ctx, id, &page); err != nil {
		return err
	}
	if _, err := x.Pages.SetIndex(ctx, page.Schema.Key, uri.Index, body.Keys, *body.Unique); err != nil {
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
	id, _ := primitive.ObjectIDFromHex(uri.Id)
	var page model.Page
	if err := x.Pages.FindOneById(ctx, id, &page); err != nil {
		return err
	}
	if _, err := x.Pages.DeleteIndex(ctx, page.Schema.Key, uri.Index); err != nil {
		return err
	}
	return nil
}
