package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

var (
	Key    string
	Method jwt.SigningMethod = jwt.SigningMethodHS256
	Exp    time.Duration     = time.Hour
)

func Make(claims jwt.MapClaims) (tokenString string, err error) {
	claims["jti"] = uuid.New()
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(Exp).Unix()
	token := jwt.NewWithClaims(Method, claims)
	tokenString, err = token.SignedString([]byte(Key))
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
