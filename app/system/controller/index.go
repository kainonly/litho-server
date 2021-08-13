package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/hash"
)

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
	data, err := x.AdminService.FindByUsername(body.Username)
	if err != nil {
		return err
	}
	if err := hash.Verify(body.Password, data.Password); err != nil {
		return err
	}
	tokenString, err := x.auth.Create(data.UID.String(), nil)
	if err != nil {
		return err
	}
	x.Cookie.Set(c, "system_access_token", tokenString)
	return "ok"
}

func (x *Index) Verify(c *gin.Context) interface{} {
	return gin.H{}
}

func (x *Index) Logout(c *gin.Context) interface{} {
	return gin.H{}
}
