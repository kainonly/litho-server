package x

import (
	"api/app/x/admin"
	"api/app/x/page"
	"api/common"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/thoas/go-funk"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InjectController struct {
	*common.App
	Service      *Service
	AdminService *admin.Service
	PageService  *page.Service
}

type Controller struct {
	*InjectController
	Auth *passport.Auth
}

func NewController(i *InjectController) *Controller {
	return &Controller{
		InjectController: i,
		Auth:             i.Passport.Make("system"),
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
	data, err := x.AdminService.FindByUsername(c, body.Username)
	if err != nil {
		return err
	}
	match, err := argon2id.ComparePasswordAndHash(body.Password, data["password"].(string))
	if err != nil {
		return err
	}
	if match == false {
		return common.LoginInvalid
	}
	uid := data["_id"].(primitive.ObjectID).Hex()
	jti := uuid.New().String()
	tokenString, err := x.Auth.Create(jti, map[string]interface{}{
		"uid": uid,
	})
	if err != nil {
		return err
	}
	x.Cookie.Set(c, "system_access_token", tokenString)
	return "ok"
}

func (x *Controller) Verify(c *gin.Context) interface{} {
	tokenString, err := x.Cookie.Get(c, "system_access_token")
	if err != nil {
		return common.LoginExpired
	}
	if _, err := x.Auth.Verify(tokenString); err != nil {
		return err
	}
	return "ok"
}

func (x *Controller) Code(c *gin.Context) interface{} {
	claims, exists := c.Get("access_token")
	if !exists {
		return common.LoginExpired
	}
	jti := claims.(jwt.MapClaims)["jti"].(string)
	code := funk.RandomString(8)
	if err := x.Service.CreateVerifyCode(c, jti, code); err != nil {
		return err
	}
	return gin.H{
		"code": code,
	}
}

func (x *Controller) RefreshToken(c *gin.Context) interface{} {
	var body struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	claims, _ := c.Get("access_token")
	jti := claims.(jwt.MapClaims)["jti"].(string)
	result, err := x.Service.VerifyCode(c, jti, body.Code)
	if err != nil {
		return err
	}
	if !result {
		return common.LoginExpired
	}
	if err = x.Service.RemoveVerifyCode(c, jti); err != nil {
		return err
	}
	tokenString, err := x.Auth.Create(jti, map[string]interface{}{
		"uid": claims.(jwt.MapClaims)["uid"],
	})
	if err != nil {
		return err
	}
	x.Cookie.Set(c, "system_access_token", tokenString)
	return "ok"
}

func (x *Controller) Logout(c *gin.Context) interface{} {
	x.Cookie.Del(c, "system_access_token")
	return "ok"
}

func (x *Controller) Pages(c *gin.Context) interface{} {
	data, err := x.PageService.Get(c)
	if err != nil {
		return err
	}
	return data
}
