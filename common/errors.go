package common

import "github.com/cloudwego/hertz/pkg/common/errors"

var ErrAuthenticationExpired = errors.NewPublic("认证过期，请重新登录")
var ErrLoginNotExists = errors.NewPublic("登录账号不存在或被冻结")
var ErrLoginMaxFailures = errors.NewPublic("登录失败超出最大次数")
var ErrLoginInvalid = errors.NewPublic("登录验证无效")
var ErrSession = errors.NewPrivate("会话建立失败")
var ErrSessionInconsistent = errors.NewPublic("会话令牌不一致")
