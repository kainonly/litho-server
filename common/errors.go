package common

import "github.com/cloudwego/hertz/pkg/common/errors"

var ErrAuthenticationExpired = errors.NewPublic("认证过期，请重新登录")
var ErrLoginNotExists = errors.NewPublic("登录账号不存在或被冻结")
var ErrLoginMaxFailures = errors.NewPublic("登录失败超出最大次数")
var ErrLoginInvalid = errors.NewPublic("登录验证无效")
var ErrSession = errors.NewPrivate("会话建立失败")
var ErrSessionInconsistent = errors.NewPublic("会话令牌不一致")
var ErrTotpInvalid = errors.NewPublic("口令验证码无效")
var ErrSmsInvalid = errors.NewPublic("短信验证码无效")
var ErrSmsNotExists = errors.NewPublic("该账户不存在或被冻结")
var ErrEmailInvalid = errors.NewPublic("邮箱验证码无效")
var ErrEmailNotExists = errors.NewPublic("该账户不存在或被冻结")
var ErrCodeFrequently = errors.NewPublic("您的验证码请求频繁，请稍后再试")
