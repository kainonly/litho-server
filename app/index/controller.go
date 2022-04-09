package index

import (
	"api/app/pages"
	"api/app/users"
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/thoas/go-funk"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/password"
	"net/http"
	"time"
)

type Controller struct {
	Service *Service
	Users   *users.Service
	Pages   *pages.Service
}

func (x *Controller) Index(c *gin.Context) interface{} {
	return gin.H{
		"name": x.Service.AppName(),
		"time": time.Now(),
	}
}

//type InstallDto struct {
//	Username string `json:"username" binding:"required"`
//	Password string `json:"password" binding:"required"`
//	Email    string `json:"email" binding:"required"`
//	Template string `json:"template" binding:"omitempty,url"`
//}
//
//func (x *Controller) Install(c *gin.Context) interface{} {
//	var body InstallDto
//	if err := c.ShouldBindJSON(&body); err != nil {
//		return err
//	}
//	ctx := c.Request.Context()
//	if err := x.Service.Install(ctx, body); err != nil {
//		return err
//	}
//	// 载入自定义 Pages
//	if body.Template != "" {
//		if err := x.Service.UseTemplate(ctx, body.Template); err != nil {
//			return err
//		}
//	}
//	return nil
//}

type LoginDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (x *Controller) AuthLogin(c *gin.Context) interface{} {
	var body LoginDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	data, err := x.Users.FindOneByUsername(ctx, body.Username)
	if err != nil {
		c.Set("code", "AUTH_INCORRECT")
		return err
	}
	if err := password.Verify(body.Password, data.Password); err != nil {
		c.Set("code", "AUTH_INCORRECT")
		return err
	}
	jti := helper.Uuid()
	var ts string
	if ts, err = x.Service.Passport.Create(jti, map[string]interface{}{
		"uid": data.ID.Hex(),
	}); err != nil {
		return err
	}
	c.SetCookie("access_token", ts, 0, "", "", true, true)
	c.SetSameSite(http.SameSiteStrictMode)
	return gin.H{
		"username": data.Username,
		"name":     data.Name,
		"avatar":   data.Avatar,
		"time":     time.Now(),
	}
}

func (x *Controller) AuthVerify(c *gin.Context) interface{} {
	ts, err := c.Cookie("access_token")
	if err != nil {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return common.LoginExpired
	}
	if _, err = x.Service.Passport.Verify(ts); err != nil {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return err
	}
	return nil
}

func (x *Controller) AuthCode(c *gin.Context) interface{} {
	claims, exists := c.Get(common.TokenClaimsKey)
	if !exists {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return common.LoginExpired
	}
	jti := claims.(jwt.MapClaims)["jti"].(string)
	code := funk.RandomString(8)
	ctx := c.Request.Context()
	if err := x.Service.CreateVerifyCode(ctx, jti, code); err != nil {
		return err
	}
	return gin.H{"code": code}
}

type RefreshTokenDto struct {
	Code string `json:"code" binding:"required"`
}

func (x *Controller) AuthRefresh(c *gin.Context) interface{} {
	var body RefreshTokenDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	value, exists := c.Get(common.TokenClaimsKey)
	if !exists {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return common.LoginExpired
	}
	claims := value.(jwt.MapClaims)
	jti := claims["jti"].(string)
	ctx := c.Request.Context()
	result, err := x.Service.VerifyCode(ctx, jti, body.Code)
	if err != nil {
		return err
	}
	if !result {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return common.LoginExpired
	}
	if err = x.Service.RemoveVerifyCode(ctx, jti); err != nil {
		return err
	}
	var ts string
	if ts, err = x.Service.Passport.Create(
		jti,
		claims["context"].(map[string]interface {
		}),
	); err != nil {
		return err
	}
	c.SetCookie("access_token", ts, 0, "", "", true, true)
	c.SetSameSite(http.SameSiteStrictMode)
	return nil
}

func (x *Controller) AuthLogout(c *gin.Context) interface{} {
	c.SetCookie("access_token", "", 0, "", "", true, true)
	c.SetSameSite(http.SameSiteStrictMode)
	return nil
}

func (x *Controller) Uploader(c *gin.Context) interface{} {
	data, err := x.Service.Uploader()
	if err != nil {
		return err
	}
	return data
}

func (x *Controller) Navs(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	navs, err := x.Pages.Navs(ctx)
	if err != nil {
		return err
	}
	return navs
}

type DynamicParams struct {
	Id string `uri:"id" binding:"required,objectId"`
}

func (x *Controller) Dynamic(c *gin.Context) interface{} {
	var params DynamicParams
	if err := c.ShouldBindUri(&params); err != nil {
		return err
	}
	ctx := c.Request.Context()
	data, err := x.Pages.FindOneById(ctx, params.Id)
	if err != nil {
		return err
	}
	return data
}
