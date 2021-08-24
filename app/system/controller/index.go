package controller

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/hash"
	"github.com/kainonly/go-bit/stringx"
)

var LoginExpired = errors.New("login authentication has expired")

type Index struct {
	*Dependency
	auth *authx.Auth
}

func NewIndex(d Dependency, authx *authx.Authx) *Index {
	return &Index{
		Dependency: &d,
		auth:       authx.Make("system"),
	}
}

func (x *Index) Login(c *gin.Context) interface{} {
	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	data, err := x.AdminService.FindByUsername(c, body.Username)
	if err != nil {
		return err
	}
	if err := hash.Verify(body.Password, data.Password); err != nil {
		return err
	}
	uid := data.UUID.String()
	jti := stringx.Uuid()
	tokenString, err := x.auth.Create(jti, map[string]interface{}{
		"uid": uid,
	})
	if err != nil {
		return err
	}
	x.Cookie.Set(c, "system_access_token", tokenString)
	return "ok"
}

func (x *Index) Verify(c *gin.Context) interface{} {
	tokenString, err := x.Cookie.Get(c, "system_access_token")
	if err != nil {
		return LoginExpired
	}
	if _, err := x.auth.Verify(tokenString); err != nil {
		return err
	}
	return "ok"
}

func (x *Index) Code(c *gin.Context) interface{} {
	claims, exists := c.Get("access_token")
	if !exists {
		return LoginExpired
	}
	jti := claims.(jwt.MapClaims)["jti"].(string)
	code := stringx.Random(8)
	if err := x.IndexService.GenerateCode(c, jti, code); err != nil {
		return err
	}
	return gin.H{
		"code": code,
	}
}

func (x *Index) RefreshToken(c *gin.Context) interface{} {
	var body struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	claims, _ := c.Get("access_token")
	jti := claims.(jwt.MapClaims)["jti"].(string)
	result, err := x.IndexService.VerifyCode(c, jti, body.Code)
	if err != nil {
		return err
	}
	if !result {
		return LoginExpired
	}
	if err = x.IndexService.RemoveCode(c, jti); err != nil {
		return err
	}
	tokenString, err := x.auth.Create(jti, map[string]interface{}{
		"uid": claims.(jwt.MapClaims)["uid"],
	})
	if err != nil {
		return err
	}
	x.Cookie.Set(c, "system_access_token", tokenString)
	return "ok"
}

func (x *Index) Logout(c *gin.Context) interface{} {
	x.Cookie.Del(c, "system_access_token")
	return "ok"
}

func (x *Index) Resource(c *gin.Context) interface{} {
	data, err := x.ResourceService.GetFromCache(context.Background())
	if err != nil {
		return err
	}
	return data
}
