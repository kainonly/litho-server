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
	data, err := x.Users.FindOneById(ctx,
		claimsContext["uid"].(string),
		options.FindOne().SetProjection(bson.M{
			"password": 0,
			"roles":    0,
			"pages":    0,
			"readonly": 0,
		}),
	)
	if err != nil {
		return err
	}
	return data
}
