package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
	"van-api/helper"
	"van-api/types"
)

var (
	Key    []byte
	Method jwt.SigningMethod = jwt.SigningMethodHS256
)

func Make(scene string, claims jwt.MapClaims) (tokenString string, err error) {
	option := helper.Config.Token[scene]
	if option == (types.TokenOption{}) {
		return "", fmt.Errorf("the [%v] scene does not exist", scene)
	}
	claims["jti"] = uuid.New()
	claims["iat"] = time.Now().Unix()
	claims["iss"] = option.Issuer
	claims["aud"] = option.Audience
	claims["exp"] = time.Now().Add(time.Second * time.Duration(option.Expires)).Unix()
	token := jwt.NewWithClaims(Method, claims)
	tokenString, err = token.SignedString(Key)
	if err != nil {
		return
	}
	return
}

func Verify(tokenString string) (result bool, claims jwt.MapClaims, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return Key, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors == jwt.ValidationErrorExpired {
				// TODO: Refresh Token
				return
			}
		}
	} else {
		result = token.Valid
		claims = token.Claims.(jwt.MapClaims)
	}
	return
}
