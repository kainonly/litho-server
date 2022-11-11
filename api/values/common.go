package values

import (
	"github.com/google/wire"
	"github.com/weplanx/server/common"
	"time"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

var (
	// Secret 密文配置
	Secret = map[string]bool{
		"tencent_secret_key":        true,
		"feishu_app_secret":         true,
		"feishu_encrypt_key":        true,
		"feishu_verification_token": true,
		"email_password":            true,
		"openapi_secret":            true,
	}
	Default = common.DynamicValues{
		LoginTTL:        time.Minute * 15,
		LoginFailures:   5,
		IpLoginFailures: 10,
		IpWhitelist:     []string{},
		IpBlacklist:     []string{},
		PwdStrategy:     1,
		PwdTTL:          time.Hour * 24 * 365,
	}
)
