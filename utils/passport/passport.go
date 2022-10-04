package passport

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v4"
	"github.com/weplanx/server/common"
	"time"
)

type Passport struct {
	Values *common.Values
}

type Claims struct {
	UserId string
	jwt.RegisteredClaims
}

// Create 生成令牌
func (x *Passport) Create(userId string, jti string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    x.Values.Namespace,
			ID:        jti,
		},
	})
	return token.SignedString([]byte(x.Values.Key))
}

// Verify 验证令牌
func (x *Passport) Verify(tokenString string) (claims Claims, err error) {
	if _, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(x.Values.Key), nil
	}); err != nil {
		return
	}
	return
}

// GetClaims 获取授权标识
func GetClaims(c *app.RequestContext) (claims Claims) {
	value, ok := c.Get("identity")
	if !ok {
		return
	}
	return value.(Claims)
}
