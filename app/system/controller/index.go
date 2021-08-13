package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/hash"
	"github.com/kainonly/go-bit/str"
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

func (x *Index) Test(c *gin.Context) interface{} {
	tokenString, err := x.auth.Create(str.Uuid().String(), str.Uuid().String(), nil)
	if err != nil {
		return err
	}
	return tokenString
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
	uid := data.UID.String()
	jti := str.Uuid().String()
	tokenString, err := x.auth.Create(jti, uid, nil)
	if err != nil {
		return err
	}
	if err := x.Session.Update(jti, uid); err != nil {
		return err
	}
	x.Cookie.Set(c, "system_access_token", tokenString)
	return "ok"
}

func (x *Index) Verify(c *gin.Context) interface{} {
	tokenString, err := x.Cookie.Get(c, "system_access_token")
	if err != nil {
		return err
	}
	if _, err := x.auth.Verify(tokenString); err != nil {
		return err
	}
	return "ok"
}

func (x *Index) Logout(c *gin.Context) interface{} {
	//x.Session.Destory()
	x.Cookie.Del(c, "system_access_token")

	return "ok"
}
