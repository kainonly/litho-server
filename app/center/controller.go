package center

import (
	"api/app/users"
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Controller struct {
	Service *Service
	Users   *users.Service
}

// GetUserInfo 获取用户信息
func (x *Controller) GetUserInfo(c *gin.Context) interface{} {
	value, _ := c.Get(common.TokenClaimsKey)
	claimsContext := value.(jwt.MapClaims)["context"].(map[string]interface{})
	ctx := c.Request.Context()
	var data map[string]interface{}
	if err := x.Users.FindOneById(ctx,
		claimsContext["uid"].(string),
		&data,
		options.FindOne().SetProjection(bson.M{
			"_id":      0,
			"password": 0,
			"roles":    0,
			"pages":    0,
			"readonly": 0,
			"labels":   0,
			"status":   0,
		}),
	); err != nil {
		return err
	}
	return data
}

type SetUserInfoDto struct {
	Name         string `json:"name"`
	Avatar       string `json:"avatar"`
	Region       string `json:"region"`
	City         string `json:"city"`
	Address      string `json:"address"`
	Introduction string `json:"introduction"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

// SetUserInfo 更新用户信息
func (x *Controller) SetUserInfo(c *gin.Context) interface{} {
	var body SetUserInfoDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	value, _ := c.Get(common.TokenClaimsKey)
	claimsContext := value.(jwt.MapClaims)["context"].(map[string]interface{})
	ctx := c.Request.Context()
	if err := x.Users.UpdateOneById(ctx,
		claimsContext["uid"].(string),
		bson.M{"$set": body},
	); err != nil {
		return err
	}
	return nil
}
