package common

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"strings"
	"time"
)

type Values struct {
	// 应用设置
	App `yaml:"app"`

	// 跨域设置
	Cors `yaml:"cors"`

	// MongoDB 配置
	Database `yaml:"database"`

	// Redis 配置
	Redis `yaml:"redis"`

	// NATS 配置
	Nats `yaml:"nats"`

	// 动态配置
	DynamicValues `yaml:"-"`
}

type App struct {
	// 地址
	Address string `yaml:"address"`
	// 命名空间
	Namespace string `yaml:"namespace"`
	// 密钥
	Key string `yaml:"key"`
}

// Name 生成空间名称
func (x App) Name(v ...string) string {
	return fmt.Sprintf(`%s:%s`, x.Namespace, strings.Join(v, ":"))
}

// Subject 生成主题名称
func (x App) Subject(v string) string {
	return fmt.Sprintf(`%s.events.%s`, x.Namespace, v)
}

// Queue 生成队列名称
func (x App) Queue(v string) string {
	return fmt.Sprintf(`%s:events:%s`, x.Namespace, v)
}

type Cors struct {
	AllowOrigins     []string `yaml:"allowOrigins"`
	AllowMethods     []string `yaml:"allowMethods"`
	AllowHeaders     []string `yaml:"allowHeaders"`
	ExposeHeaders    []string `yaml:"exposeHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
	MaxAge           int      `yaml:"maxAge"`
}

type Database struct {
	Uri string `yaml:"uri"`
	Db  string `yaml:"db"`
}

type Redis struct {
	Uri string `yaml:"uri"`
}

type Nats struct {
	Hosts []string `yaml:"hosts"`
	Nkey  string   `yaml:"nkey"`
}

// DynamicValues 动态配置
type DynamicValues map[string]interface{}

// GetSessionTTL  会话周期（秒）
// 用户在 1 小时 内没有操作，将结束会话。
func (x DynamicValues) GetSessionTTL() time.Duration {
	return time.Second * time.Duration(x["session_ttl"].(float64))
}

// GetLoginTTL 登录锁定时间
// 锁定 15 分钟。
func (x DynamicValues) GetLoginTTL() time.Duration {
	return time.Second * time.Duration(x["login_ttl"].(float64))
}

// GetLoginFailures 用户最大登录失败次数
// 有限时间（锁定时间）内连续登录失败 5 次，锁定帐号。
func (x DynamicValues) GetLoginFailures() int64 {
	return int64(x["login_failures"].(float64))
}

// GetIpLoginFailures IP 最大登录失败次数
// 同 IP 连续 10 次登录失败后，锁定 IP（周期为锁定时间）。
func (x DynamicValues) GetIpLoginFailures() int64 {
	return int64(x["ip_login_failures"].(float64))
}

// GetIpWhitelist  IP 白名单
// 白名单 IP 允许超出最大登录失败次数。
func (x DynamicValues) GetIpWhitelist() []string {
	return x["ip_whitelist"].([]string)
}

// GetIpBlacklist IP 黑名单
// 黑名单 IP 将禁止访问。
func (x DynamicValues) GetIpBlacklist() []string {
	return x["ip_blacklist"].([]string)
}

// GetPwdStrategy 密码强度
// 0：无限制；
// 1：需要大小写字母；
// 2：需要大小写字母、数字；
// 3：需要大小写字母、数字、特殊字符
func (x DynamicValues) GetPwdStrategy() int64 {
	return int64(x["pwd_strategy"].(float64))
}

// GetPwdTTL 密码有效期（天）
// 密码有效期（天）
// 密码过期后强制要求修改密码，0：永久有效
func (x DynamicValues) GetPwdTTL() time.Duration {
	return 24 * time.Hour * time.Duration(x["pwd_ttl"].(float64))
}

// GetCloud 云平台
// tencent：腾讯云；
func (x DynamicValues) GetCloud() string {
	return x["cloud"].(string)
}

// GetTencentSecretId 腾讯云 API 密钥 Id
// 建议用子账号分配需要的权限
func (x DynamicValues) GetTencentSecretId() string {
	return x["tencent_secret_id"].(string)
}

// GetTencentSecretKey 腾讯云 API 密钥 Key
func (x DynamicValues) GetTencentSecretKey() string {
	return x["tencent_secret_key"].(string)
}

// GetTencentCosBucket 腾讯云 COS 对象存储 Bucket（存储桶名称）
func (x DynamicValues) GetTencentCosBucket() string {
	return x["tencent_cos_bucket"].(string)
}

// GetTencentCosRegion 腾讯云 COS 对象存储所属地域
// 例如：ap-guangzhou
func (x DynamicValues) GetTencentCosRegion() string {
	return x["tencent_cos_region"].(string)
}

// GetTencentCosExpired 腾讯云 COS 对象存储预签名有效期，单位：秒
func (x DynamicValues) GetTencentCosExpired() time.Duration {
	return time.Second * time.Duration(x["tencent_cos_expired"].(float64))
}

// GetTencentCosLimit 腾讯云 COS 对象存储上传大小限制，单位：KB
func (x DynamicValues) GetTencentCosLimit() int64 {
	return int64(x["tencent_cos_limit"].(int))
}

// GetOffice 办公平台
// feishu：飞书；
func (x DynamicValues) GetOffice() string {
	return x["office"].(string)
}

// GetFeishuAppId 飞书应用 ID
func (x DynamicValues) GetFeishuAppId() string {
	return x["feishu_app_id"].(string)
}

// GetFeishuAppSecret 飞书应用密钥
func (x DynamicValues) GetFeishuAppSecret() string {
	return x["feishu_app_secret"].(string)
}

// GetFeishuEncryptKey 飞书事件订阅安全校验数据密钥
func (x DynamicValues) GetFeishuEncryptKey() string {
	return x["feishu_encrypt_key"].(string)
}

// GetFeishuVerificationToken 飞书事件订阅验证令牌
func (x DynamicValues) GetFeishuVerificationToken() string {
	return x["feishu_verification_token"].(string)
}

// GetRedirectUrl 第三方免登授权码跳转地址
func (x DynamicValues) GetRedirectUrl() string {
	return x["redirect_url"].(string)
}

// GetEmailHost 公共电子邮件服务 SMTP 地址
func (x DynamicValues) GetEmailHost() string {
	return x["email_host"].(string)
}

// GetEmailPort SMTP 端口号（SSL）
func (x DynamicValues) GetEmailPort() string {
	return x["email_port"].(string)
}

// GetEmailUsername 公共邮箱用户
// 例如：support@example.com
func (x DynamicValues) GetEmailUsername() string {
	return x["email_username"].(string)
}

// GetEmailPassword 公共邮箱用户密码
func (x DynamicValues) GetEmailPassword() string {
	return x["email_password"].(string)
}

// GetOpenapiUrl 开放服务地址
func (x DynamicValues) GetOpenapiUrl() string {
	return x["openapi_url"].(string)
}

// GetOpenapiKey 开放服务应用认证 Key
// API 网关应用认证方式 https://cloud.tencent.com/document/product/628/55088
func (x DynamicValues) GetOpenapiKey() string {
	return x["openapi_key"].(string)
}

// GetOpenapiSecret 开放服务应用认证密钥
func (x DynamicValues) GetOpenapiSecret() string {
	return x["openapi_secret"].(string)
}

// Active 授权用户标识
type Active struct {
	// Token ID
	JTI string

	// User ID
	UID string
}

// GetActive 获取授权用户标识
func GetActive(c *app.RequestContext) (data Active) {
	value, ok := c.Get("identity")
	if !ok {
		return
	}
	return value.(Active)
}
