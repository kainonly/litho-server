package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"time"
)

type Controller struct {
	IndexService *Service
}

// Index 入口
// @router / [GET]
func (x *Controller) Index(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, utils.H{
		"ip":   c.ClientIP(),
		"time": time.Now(),
	})
}

// GetRefreshCode 获取刷新令牌验证码
// @router /code [GET]
func (x *Controller) GetRefreshCode(ctx context.Context, c *app.RequestContext) {
	//active := common.GetActive(c)
	//code, err := x.IndexService.GetRefreshCode(ctx, active.UID)
	//if err != nil {
	//	c.Error(err)
	//	return
	//}
	//
	//c.JSON(http.StatusOK, utils.H{
	//	"code": code,
	//})
}

// VerifyRefreshCode 校验刷新令牌验证码
// @router /refresh_token [POST]
func (x *Controller) VerifyRefreshCode(ctx context.Context, c *app.RequestContext) {
	//var dto struct {
	//	Code string `json:"code,required"`
	//}
	//if err := c.BindAndValidate(&dto); err != nil {
	//	c.Error(err)
	//	return
	//}
	//
	//active := common.GetActive(c)
	//if err := x.IndexService.Captcha.Verify(ctx, active.UID, dto.Code); err != nil {
	//	c.Error(err)
	//	return
	//}
	//
	//c.Next(ctx)
}

// GetNavs 导航数据
// @router /navs [GET]
func (x *Controller) GetNavs(ctx context.Context, c *app.RequestContext) {
	//active := common.GetActive(c)
	//
	//data, err := x.IndexService.GetNavs(ctx, active.UID)
	//if err != nil {
	//	return
	//}
	//
	//c.JSON(http.StatusOK, data)
}

// GetOptions 返回通用配置
// @router /options [GET]
func (x *Controller) GetOptions(ctx context.Context, c *app.RequestContext) {
	//var dto struct {
	//	// 类型
	//	Type string `query:"type,required"`
	//}
	//if err := c.BindAndValidate(&dto); err != nil {
	//	c.Error(err)
	//	return
	//}
	//
	//data := x.IndexService.GetOptions(dto.Type)
	//if data == nil {
	//	c.Status(http.StatusNoContent)
	//	return
	//}
	//
	//c.JSON(http.StatusOK, data)
}

// GetUser 获取授权用户信息
// @router /user [GET]
func (x *Controller) GetUser(ctx context.Context, c *app.RequestContext) {
	//active := common.GetActive(c)
	//data, err := x.IndexService.GetUser(ctx, active.UID)
	//if err != nil {
	//	c.Error(err)
	//	return
	//}
	//
	//c.JSON(http.StatusOK, data)
}

type SetUserDto struct {
	// 用户名
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	// 电子邮件
	Email string `json:"email,omitempty" bson:"email,omitempty" vd:"$=='' || email($)"`
	// 称呼
	Name string `json:"name" bson:"name,omitempty"`
	// 头像
	Avatar string `json:"avatar" bson:"avatar,omitempty"`
	// 密码
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	// 重置
	Reset string `json:"reset,omitempty" vd:"in($, 'feishu')" bson:"reset"`
	// 更新时间
	UpdateTime time.Time `json:"-" bson:"update_time"`
}

// SetUser 设置授权用户信息
// @router /user [PATCH]
func (x *Controller) SetUser(ctx context.Context, c *app.RequestContext) {
	//var dto SetUserDto
	//if err := c.BindAndValidate(&dto); err != nil {
	//	c.Error(err)
	//	return
	//}
	//
	//// 密码转散列
	//if dto.Password != "" {
	//	dto.Password, _ = passlib.Hash(dto.Password)
	//}
	//
	//dto.UpdateTime = time.Now()
	//active := common.GetActive(c)
	//if _, err := x.IndexService.SetUser(ctx, active.UID, dto); err != nil {
	//	c.Error(err)
	//	return
	//}
	//
	//// 用户名变更，注销登录状态
	//if dto.Username != "" {
	//	cookie := &protocol.Cookie{}
	//	cookie.SetKey("access_token")
	//	cookie.SetValue("")
	//	c.Response.Header.SetCookie(cookie)
	//
	//	if err := x.IndexService.LogoutSession(ctx, active.UID); err != nil {
	//		c.Error(err)
	//		return
	//	}
	//}
	//
	//c.Status(http.StatusNoContent)
}
