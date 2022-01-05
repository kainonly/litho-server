package index

import (
	"api/app/pages"
	"api/app/users"
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/thoas/go-funk"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/password"
	"net/http"
)

type Controller struct {
	Service *Service
	Users   *users.Service
	Pages   *pages.Service
}

func (x *Controller) Index(c *gin.Context) interface{} {
	return gin.H{
		"name": x.Service.AppName(),
		"v":    "1.0.0",
	}
}

func (x *Controller) Login(c *gin.Context) interface{} {
	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	data, err := x.Users.FindByUsername(ctx, body.Username)
	if err != nil {
		return err
	}
	if err := password.Verify(body.Password, data.Password); err != nil {
		return err
	}
	uid := data.ID.Hex()
	jti := helper.Uuid()
	tokenString, _ := x.Service.Passport.Create(jti, map[string]interface{}{
		"uid": uid,
	})
	c.SetCookie("access_token", tokenString, 0, "", "", true, true)
	c.SetSameSite(http.SameSiteStrictMode)
	return nil
}

func (x *Controller) Verify(c *gin.Context) interface{} {
	tokenString, err := c.Cookie("access_token")
	if err != nil {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return common.LoginExpired
	}
	if _, err := x.Service.Passport.Verify(tokenString); err != nil {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return err
	}
	return nil
}

func (x *Controller) Code(c *gin.Context) interface{} {
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

func (x *Controller) RefreshToken(c *gin.Context) interface{} {
	var body struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	claims, exists := c.Get(common.TokenClaimsKey)
	if !exists {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return common.LoginExpired
	}
	jti := claims.(jwt.MapClaims)["jti"].(string)
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
	tokenString, _ := x.Service.Passport.Create(jti, map[string]interface{}{
		"uid": claims.(jwt.MapClaims)["uid"],
	})
	c.SetCookie("access_token", tokenString, 0, "", "", true, true)
	c.SetSameSite(http.SameSiteStrictMode)
	return nil
}

func (x *Controller) Logout(c *gin.Context) interface{} {
	c.SetCookie("access_token", "", 0, "", "", true, true)
	c.SetSameSite(http.SameSiteStrictMode)
	return nil
}

func (x *Controller) Api(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	pagesData, err := x.Pages.Fetch(ctx)
	if err != nil {
		return err
	}
	return gin.H{
		"pages": pagesData,
	}
}
