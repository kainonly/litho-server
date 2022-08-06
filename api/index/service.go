package index

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/utils"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/weplanx/server/api/departments"
	"github.com/weplanx/server/api/pages"
	"github.com/weplanx/server/api/roles"
	"github.com/weplanx/server/api/sessions"
	"github.com/weplanx/server/api/users"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/common/captcha"
	"github.com/weplanx/server/common/locker"
	"github.com/weplanx/server/common/passlib"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Service struct {
	*common.Inject

	SessionService     *sessions.Service
	PagesService       *pages.Service
	UsersService       *users.Service
	RolesService       *roles.Service
	DepartmentsService *departments.Service

	Captcha *captcha.Captcha
	Locker  *locker.Locker
}

// Login 登录
func (x *Service) Login(ctx context.Context, identity string, password string) (_ common.Active, err error) {
	var user model.User
	if user, err = x.UsersService.FindByIdentity(ctx, identity); err != nil {
		return
	}
	uid := user.ID.Hex()

	// 锁定上限验证
	var maxFailures bool
	if maxFailures, err = x.Locker.Verify(ctx, uid, x.Values.GetLoginFailures()); err != nil {
		return
	}
	if maxFailures {
		err = errors.NewPublic("用户登录失败已超出最大次数")
		return
	}

	// 验证密码正确性
	if err = passlib.Verify(password, user.Password); err != nil {
		// 锁定更新
		if err = x.Locker.Update(ctx, uid, x.Values.GetLoginTTL()); err != nil {
			return
		}
		if err == passlib.ErrNotMatch {
			err = errors.New(err, errors.ErrorTypePublic, nil)
		}
		return
	}

	// 令牌 ID
	jti, _ := gonanoid.Nanoid()
	return common.Active{
		JTI: jti,
		UID: uid,
	}, nil
}

// LoginSession 建立登录会话，移除锁定
func (x *Service) LoginSession(ctx context.Context, uid string, jti string) (err error) {
	if err = x.Locker.Delete(ctx, jti); err != nil {
		return
	}
	if err = x.SessionService.Set(ctx, uid, jti); err != nil {
		return
	}
	return
}

// AuthVerify 认证鉴权、权限验证、会话续约
func (x *Service) AuthVerify(ctx context.Context, uid string, jti string) (err error) {
	var result bool
	// 检测会话
	if result, err = x.SessionService.Verify(ctx, uid, jti); err != nil {
		return
	}
	if !result {
		err = errors.NewPublic("会话令牌不一致")
		return
	}

	// TODO: Check User Status

	// 会话续约
	return x.SessionService.Renew(ctx, uid)
}

// LogoutSession 注销登录会话
func (x *Service) LogoutSession(ctx context.Context, uid string) (err error) {
	return x.SessionService.Remove(ctx, uid)
}

// GetNavs 导航数据
func (x *Service) GetNavs(ctx context.Context, uid string) (_ []pages.Nav, err error) {
	// TODO: 权限过滤...
	//var user model.User
	//if user, err = x.UsersService.GetActived(ctx, uid); err != nil {
	//	return
	//}

	return x.PagesService.GetNavs(ctx)
}

// GetRefreshCode 获取刷新令牌验证码
func (x *Service) GetRefreshCode(ctx context.Context, uid string) (code string, err error) {
	if code, err = gonanoid.Nanoid(); err != nil {
		return
	}
	if err = x.Captcha.Create(ctx, uid, code, 15*time.Second); err != nil {
		return
	}
	return
}

// GetOptions 返回通用配置
func (x *Service) GetOptions(v string) utils.H {
	switch v {
	// 上传类
	case "upload":
		switch x.Values.GetCloud() {
		// 腾讯云
		case "tencent":
			// Cos 对象存储
			return utils.H{
				"type": "cos",
				"url": fmt.Sprintf(`https://%s.cos.%s.myqcloud.com`,
					x.Values.GetTencentCosBucket(), x.Values.GetTencentCosRegion(),
				),
				"limit": x.Values.GetTencentCosLimit(),
			}
		}
	// 企业平台类
	case "office":
		switch x.Values.GetOffice() {
		// 飞书
		case "feishu":
			return utils.H{
				"url":      "https://open.feishu.cn/open-apis/authen/v1/index",
				"redirect": x.Values.GetRedirectUrl(),
				"app_id":   x.Values.GetFeishuAppId(),
			}
		}
	}
	return nil
}

// GetUser 获取授权用户信息
func (x *Service) GetUser(ctx context.Context, uid string) (data map[string]interface{}, err error) {
	var user model.User
	if user, err = x.UsersService.GetActived(ctx, uid); err != nil {
		return
	}

	data = map[string]interface{}{
		"username":    user.Username,
		"email":       user.Email,
		"name":        user.Name,
		"avatar":      user.Avatar,
		"sessions":    user.Sessions,
		"last":        user.Last,
		"create_time": user.CreateTime,
	}

	// 权限组名称
	if data["roles"], err = x.RolesService.FindNamesByIds(ctx, user.Roles); err != nil {
		return
	}

	// 部门名称
	if user.Department != nil {
		var department model.Department
		if department, err = x.DepartmentsService.FindOneById(ctx, *user.Department); err != nil {
			return
		}
		data["department"] = department.Name
	}

	return
}

// SetUser 设置授权用户信息
func (x *Service) SetUser(ctx context.Context, id string, data interface{}) (interface{}, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	return x.UsersService.UpdateOneById(ctx, oid, bson.M{"$set": data})
}

// UnsetUser 取消授权用户信息
func (x *Service) UnsetUser(ctx context.Context, id string, mate string) (interface{}, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	update := bson.M{
		"$set":   bson.M{"update_time": time.Now()},
		"$unset": bson.M{mate: ""},
	}
	return x.UsersService.UpdateOneById(ctx, oid, update)
}
