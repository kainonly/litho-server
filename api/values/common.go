package values

import (
	"github.com/google/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

// secret 密文配置
var secret = map[string]bool{
	"tencent_secret_key":        true,
	"feishu_app_secret":         true,
	"feishu_encrypt_key":        true,
	"feishu_verification_token": true,
	"email_password":            true,
	"openapi_secret":            true,
}
