package common

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

// LoadStaticValues 加载静态配置
func LoadStaticValues(path string) (values *Values, err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("静态配置不存在，请检查路径 [%s]", path)
	}
	var b []byte
	if b, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if err = yaml.Unmarshal(b, &values); err != nil {
		return
	}
	return
}

type Values struct {
	// 应用设置
	App struct {
		// 命名空间
		Namespace string `yaml:"namespace"`
		// 密钥
		Key string `yaml:"key"`
	} `yaml:"app"`

	// 跨域设置
	Cors struct {
		AllowOrigins     []string `yaml:"allowOrigins"`
		AllowMethods     []string `yaml:"allowMethods"`
		AllowHeaders     []string `yaml:"allowHeaders"`
		ExposeHeaders    []string `yaml:"exposeHeaders"`
		AllowCredentials bool     `yaml:"allowCredentials"`
		MaxAge           int      `yaml:"maxAge"`
	} `yaml:"cors"`

	// MongoDB 配置
	Database struct {
		Uri string `yaml:"uri"`
		Db  string `yaml:"db"`
	} `yaml:"database"`

	// Redis 配置
	Redis struct {
		Uri string `yaml:"uri"`
	} `yaml:"redis"`

	// NATS 配置
	Nats struct {
		Hosts []string `yaml:"hosts"`
		Nkey  string   `yaml:"nkey"`
	} `yaml:"nats"`

	// 动态配置
	//*DynamicValues `yaml:"-"`
}

// Key 空间命名
func (x *Values) Key(v ...string) string {
	return fmt.Sprintf(`%s:%s`, x.App.Namespace, strings.Join(v, ":"))
}

// Event 事件主题
func (x *Values) Event(v string) string {
	return fmt.Sprintf(`%s.events.%s`, x.App.Namespace, v)
}

// EventName 事件名称
func (x *Values) EventName(v string) string {
	return fmt.Sprintf(`%s:events:%s`, x.App.Namespace, v)
}

//type DynamicValues struct {
//	// 用户无操作的最大时长，超出将结束会话
//	UserSessionExpire time.Duration `json:"user_session_expire"`
//
//	// 用户有限时间（等同锁定时间）内连续登录失败的次数，超出锁定帐号
//	UserLoginFailedTimes int64 `json:"user_login_failed_times"`
//
//	// 锁定账户时间
//	UserLockTime time.Duration `json:"user_lock_time"`
//
//	// IP 连续登录失败后的最大次数（无白名单时启用），锁定 IP
//	IpLoginFailedTimes int64 `json:"ip_login_failed_times"`
//
//	// IP 白名单
//	IpWhitelist []string `json:"ip_whitelist"`
//
//	// IP 黑名单
//	IpBlacklist []string `json:"ip_blacklist"`
//
//	// 密码强度
//	// 0：无限制；1：需要大小写字母；2：需要大小写字母、数字；3：需要大小写字母、数字、特殊字符
//	PasswordStrength int64 `json:"password_strength"`
//
//	// 密码有效期（天）
//	// 密码过期后强制要求修改密码，0：永久有效
//	PasswordExpire int64 `json:"password_expire"`
//
//	// 云厂商
//	// 腾讯云
//	CloudPlatform string `json:"cloud_platform"`
//
//	// 腾讯云 API 密钥 Id，建议用子账号分配需要的权限
//	TencentSecretId string `json:"tencent_secret_id"`
//
//	// 腾讯云 API 密钥 Key
//	TencentSecretKey string `json:"tencent_secret_key"`
//
//	// 腾讯云 COS 对象存储 Bucket（存储桶名称）
//	TencentCosBucket string `json:"tencent_cos_bucket"`
//
//	// 腾讯云 COS 对象存储所属地域，例如：ap-guangzhou
//	TencentCosRegion string `json:"tencent_cos_region"`
//
//	// 腾讯云 COS 对象存储预签名有效期，单位：秒
//	TencentCosExpired time.Duration `json:"tencent_cos_expired"`
//
//	// 腾讯云 COS 对象存储上传大小限制，单位：KB
//	TencentCosLimit int `json:"tencent_cos_limit"`
//
//	// 办公平台
//	// 飞书
//	OfficePlatform string `json:"office_platform"`
//
//	// 飞书应用 ID
//	FeishuAppId string `json:"feishu_app_id"`
//
//	// 飞书应用密钥
//	FeishuAppSecret string `json:"feishu_app_secret"`
//
//	// 飞书事件订阅安全校验数据密钥
//	FeishuEncryptKey string `json:"feishu_encrypt_key"`
//
//	// 飞书事件订阅验证令牌
//	FeishuVerificationToken string `json:"feishu_verification_token"`
//
//	// 第三方免登授权码跳转地址
//	RedirectUrl string `json:"redirect_url"`
//
//	// 公共电子邮件服务 SMTP 地址
//	EmailHost string `json:"email_host"`
//
//	// SMTP 端口号（SSL）
//	EmailPort string `json:"email_port"`
//
//	// 公共邮箱用户，例如：support@example.com
//	EmailUsername string `json:"email_username"`
//
//	// 公共邮箱用户密码
//	EmailPassword string `json:"email_password"`
//
//	// 开放服务地址
//	OpenapiUrl string `json:"openapi_url"`
//
//	// 开放服务应用认证 Key
//	// API 网关应用认证方式 https://cloud.tencent.com/document/product/628/55088
//	OpenapiKey string `json:"openapi_key"`
//
//	// 开放服务应用认证密钥
//	OpenapiSecret string `json:"openapi_secret"`
//}
