package passport

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Passport struct {

	// 应用命名
	Key string `yaml:"key"`

	// 设置
	Option `yaml:"option"`
}

type Option struct {
	// 签发人
	Iss string `yaml:"iss"`

	// 受众
	Aud []string `yaml:"aud"`

	// 主题
	Sub string `yaml:"sub"`

	// 未来生效的时间
	Nbf int64 `yaml:"nbf"`

	// 过期时间
	Exp int64 `yaml:"exp"`
}

// Create authentication token
func (x *Passport) Create(jti string, context map[string]interface{}) (tokenString string, err error) {
	claims := jwt.MapClaims{
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Add(time.Second * time.Duration(x.Nbf)).Unix(),
		"exp":     time.Now().Add(time.Second * time.Duration(x.Exp)).Unix(),
		"jti":     jti,
		"iss":     x.Iss,
		"aud":     x.Aud,
		"sub":     x.Sub,
		"context": context,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(x.Key))
}

// Verify authentication token
func (x *Passport) Verify(tokenString string) (claims jwt.MapClaims, err error) {
	var token *jwt.Token
	if token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("验证失败，签名方式不一致")
		}
		return []byte(x.Key), nil
	}); err != nil {
		return
	}
	return token.Claims.(jwt.MapClaims), nil
}

type Claims struct {
	jwt.MapClaims
	Context map[string]interface{}
}

//func (x *Passport) GetClaims(c *gin.Context, key string) (*Claims, error) {
//	value, exists := c.Get(key)
//	if !exists {
//		c.Set("status_code", 401)
//		c.Set("code", "AUTH_EXPIRED")
//		return nil, ErrAuthExpired
//	}
//	claims := value.(jwt.MapClaims)
//	return &Claims{
//		MapClaims: claims,
//		Context:   claims["context"].(map[string]interface{}),
//	}, nil
//}
