package index

import (
	"api/app/admin"
	"api/app/page"
	"api/common"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/thoas/go-funk"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/password"
	"github.com/weplanx/go/route"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	*InjectController
}

type InjectController struct {
	common.Inject
	Service      *Service
	AdminService *admin.Service
	PageService  *page.Service
}

func (x *Controller) Login(c *fiber.Ctx) interface{} {
	var body struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	ctx := c.UserContext()
	data, err := x.AdminService.FindByUsername(ctx, body.Username)
	if err != nil {
		return err
	}
	if err := password.Verify(body.Password, data["password"].(string)); err != nil {
		return err
	}
	uid := data["_id"].(primitive.ObjectID).Hex()
	jti := helper.Uuid()
	tokenString, err := x.Passport.Create(jti, map[string]interface{}{
		"uid": uid,
	})
	if err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	})
	return "ok"
}

func (x *Controller) Verify(c *fiber.Ctx) interface{} {
	tokenString := c.Cookies("access_token")
	if tokenString == "" {
		return route.E{Code: 401, Message: common.LoginExpired.Error()}
	}
	if _, err := x.Passport.Verify(tokenString); err != nil {
		return err
	}
	return "ok"
}

func (x *Controller) Code(c *fiber.Ctx) interface{} {
	claims := c.Locals(common.TokenClaimsKey)
	if claims == nil {
		return common.LoginExpired
	}
	jti := claims.(jwt.MapClaims)["jti"].(string)
	code := funk.RandomString(8)
	ctx := c.UserContext()
	if err := x.Service.CreateVerifyCode(ctx, jti, code); err != nil {
		return err
	}
	return fiber.Map{"code": code}
}

func (x *Controller) RefreshToken(c *fiber.Ctx) interface{} {
	var body struct {
		Code string `json:"code" validate:"required"`
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	claims := c.Locals(common.TokenClaimsKey)
	if claims == nil {
		return common.LoginExpired
	}
	jti := claims.(jwt.MapClaims)["jti"].(string)
	ctx := c.UserContext()
	result, err := x.Service.VerifyCode(ctx, jti, body.Code)
	if err != nil {
		return err
	}
	if !result {
		return common.LoginExpired
	}
	if err = x.Service.RemoveVerifyCode(ctx, jti); err != nil {
		return err
	}
	tokenString, err := x.Passport.Create(jti, map[string]interface{}{
		"uid": claims.(jwt.MapClaims)["uid"],
	})
	if err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	})
	return "ok"
}

func (x *Controller) Logout(c *fiber.Ctx) interface{} {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	})
	return "ok"
}

func (x *Controller) Pages(c *fiber.Ctx) interface{} {
	ctx := c.UserContext()
	data, err := x.PageService.Get(ctx)
	if err != nil {
		return err
	}
	return data
}
