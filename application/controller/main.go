package controller

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/authx"
	"github.com/kainonly/gin-extra/hash"
	"github.com/kainonly/gin-extra/mvcx"
	"github.com/kainonly/gin-extra/storage/cos"
	"lab-api/application/common"
	"lab-api/application/model"
	"mime/multipart"
)

type Controller struct {
	common.Dependency
}

type LoginBody struct {
	Username string `binding:"required,min=4,max=20"`
	Password string `binding:"required,min=12,max=20"`
}

func (c *Controller) Login(ctx *gin.Context) interface{} {
	var body LoginBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	redisCtx := context.Background()
	data := c.Redis.Admin.Get(redisCtx, body.Username)
	if data == nil {
		return errors.New("当前用户不存在或已被冻结")
	}
	userKey := "admin:"
	if !c.Redis.Lock.Check(redisCtx, userKey+body.Username) {
		c.Redis.Lock.Lock(redisCtx, userKey+body.Username)
		return mvcx.Response{
			Code: 2,
			Msg:  "当前用户登录失败次数已上限，请稍后再试",
		}
	}
	if ok, _ := hash.Verify(body.Password, data["password"].(string)); !ok {
		c.Redis.Lock.Inc(redisCtx, userKey+body.Username)
		return errors.New("当前用户认证不成功")
	}
	c.Redis.Lock.Remove(redisCtx, userKey+body.Username)
	if err := authx.Create(ctx, common.SystemCookie, jwt.MapClaims{
		"user": data["username"],
	}, c.Redis.RefreshToken); err != nil {
		return err
	}
	return true
}

func (c *Controller) Verify(ctx *gin.Context) interface{} {
	if err := authx.Verify(ctx, common.SystemCookie, c.Redis.RefreshToken); err != nil {
		return err
	}
	return true
}

func (c *Controller) Logout(ctx *gin.Context) interface{} {
	if err := authx.Destory(ctx, "system", c.Redis.RefreshToken); err != nil {
		return err
	}
	return true
}

func (c *Controller) Resource(ctx *gin.Context) interface{} {
	var err error
	var auth jwt.MapClaims
	if auth, err = authx.Get(ctx); err != nil {
		return err
	}
	redisCtx := context.Background()
	user := auth["user"].(string)
	data := c.Redis.Admin.Get(redisCtx, user)
	roles := data["role"].([]interface{})
	roleKeys := make([]string, len(roles))
	for index, value := range roles {
		roleKeys[index] = value.(string)
	}
	roleSet := c.Redis.Role.Get(redisCtx, roleKeys, "resource")
	if data["resource"] != nil {
		roleSet.Add(data["resource"].([]interface{})...)
	}
	lists := make([]interface{}, 0)
	resource := c.Redis.Resource.Get(redisCtx)
	for _, item := range resource {
		if roleSet.Contains(item["key"].(string)) {
			lists = append(lists, item)
		}
	}
	return lists
}

func (c *Controller) Information(ctx *gin.Context) interface{} {
	auth, exists := ctx.Get("auth")
	if !exists {
		return false
	}
	data := make(map[string]interface{})
	c.Db.Model(&model.Admin{}).
		Where("username = ?", auth.(jwt.MapClaims)["user"]).
		Omit("password", "permission", "status", "create_time", "update_time").
		First(&data)
	return data
}

type updateBody struct {
	OldPassword string
	NewPassword string
	Email       string
	Phone       string
	Call        string
	Avatar      string
}

func (c *Controller) Update(ctx *gin.Context) interface{} {
	var body updateBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	auth, exists := ctx.Get("auth")
	if !exists {
		return false
	}
	username := auth.(jwt.MapClaims)["user"]
	data := model.Admin{
		Email:  body.Email,
		Phone:  body.Phone,
		Call:   body.Call,
		Avatar: body.Avatar,
	}
	if body.NewPassword != "" || body.OldPassword != "" {
		if body.NewPassword != body.OldPassword {
			return errors.New("the old and new password verification is inconsistent")
		}
		var adminData model.Admin
		c.Db.Model(&model.Admin{}).
			Where("username = ?", username).
			First(&adminData)
		if ok, _ := hash.Verify(body.OldPassword, adminData.Password); !ok {
			return mvcx.Response{
				Code: 2,
				Msg:  "password verification failed",
			}
		}
		data.Password, _ = hash.Make(body.NewPassword)
	}
	c.Db.Model(&model.Admin{}).
		Where("username = ?", username).
		Updates(data)
	c.Redis.Admin.Clear(context.Background())
	return true
}

func (c *Controller) Uploads(ctx *gin.Context) interface{} {
	var err error
	var fileHeader *multipart.FileHeader
	if fileHeader, err = ctx.FormFile("file"); err != nil {
		return err
	}
	var fileName string
	if fileName, err = cos.Put(fileHeader); err != nil {
		return err
	}
	return gin.H{
		"savename": fileName,
	}
}

func (c *Controller) CosPresigned(ctx *gin.Context) interface{} {
	data, err := cos.GeneratePostPresigned(600, []interface{}{
		"content-length-range", 0, 104857600,
	})
	if err != nil {
		return err
	}
	return mvcx.Raw(data)
}
