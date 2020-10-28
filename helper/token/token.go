package token

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kataras/jwt"
	"time"
)

type Option struct {
	Issuer   string   `yaml:"issuer"`
	Audience []string `yaml:"audience"`
	Expires  uint     `yaml:"expires"`
}
type Handle func(option Option) (claims map[string]interface{}, err error)

var (
	Key     []byte
	Options map[string]Option
	Method  jwt.Alg = jwt.HS256
)

func Make(scene string, claims map[string]interface{}) (token []byte, err error) {
	option, exists := Options[scene]
	if !exists {
		err = fmt.Errorf("the [%v] scene does not exist", scene)
		return
	}
	claims["jti"] = uuid.New()
	claims["iss"] = option.Issuer
	claims["aud"] = option.Audience
	token, err = jwt.Sign(Method, Key, claims, jwt.MaxAge(time.Second*time.Duration(option.Expires)))
	if err != nil {
		return
	}
	return
}

func Verify(scene string, token []byte, refresh Handle) (claims map[string]interface{}, err error) {
	option, exists := Options[scene]
	if !exists {
		err = fmt.Errorf("the [%v] scene does not exist", scene)
		return
	}
	var verifiedToken *jwt.VerifiedToken
	verifiedToken, err = jwt.Verify(Method, Key, token)
	if err != nil {
		if err == jwt.ErrExpired && refresh != nil {
			return refresh(option)
		}
		return
	}
	err = verifiedToken.Claims(&claims)
	if err != nil {
		return
	}
	return
}
