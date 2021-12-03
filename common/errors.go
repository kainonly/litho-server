package common

import "errors"

var (
	LoginInvalid = errors.New("登录认证失败")
	LoginExpired = errors.New("登录令牌已失效")
)
