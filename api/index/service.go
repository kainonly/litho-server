package index

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/errors"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"github.com/weplanx/utils/captcha"
	"github.com/weplanx/utils/locker"
	"github.com/weplanx/utils/passlib"
	"github.com/weplanx/utils/passport"
	"github.com/weplanx/utils/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Service struct {
	*common.Inject
	Passport        *passport.Passport
	Locker          *locker.Locker
	Captcha         *captcha.Captcha
	SessionsService *sessions.Service
}

// Login 登录
func (x *Service) Login(ctx context.Context, email string, password string) (ts string, err error) {
	var user model.User
	if err = x.Db.Collection("users").FindOne(ctx, bson.M{
		"email":  email,
		"status": true,
	}).Decode(&user); err != nil {
		return
	}

	userId := user.ID.Hex()

	// 锁定上限验证
	var maxLoginFailures bool
	if maxLoginFailures, err = x.Locker.Verify(ctx, userId, x.Values.LoginFailures); err != nil {
		return
	}
	if maxLoginFailures {
		err = errors.NewPublic("用户登录失败已超出最大次数")
		return
	}

	// 验证密码正确性
	var match bool
	if match, err = passlib.Verify(password, user.Password); err != nil {
		return
	}
	if !match {
		// 锁定更新
		if err = x.Locker.Update(ctx, userId, x.Values.LoginTTL); err != nil {
			return
		}
		err = errors.NewPublic("用户名称或用户密码不正确")
		return
	}

	// 生成令牌
	jti, _ := gonanoid.Nanoid()
	if ts, err = x.Passport.Create(userId, jti); err != nil {
		return
	}

	// 锁定清零
	if err = x.Locker.Delete(ctx, userId); err != nil {
		return
	}

	// 设置会话
	if err = x.SessionsService.Set(ctx, userId, jti); err != nil {
		return
	}

	// 用户缓存刷新
	key := x.Values.Name("users", userId)
	if _, err = x.Redis.Del(ctx, key).Result(); err != nil {
		return
	}

	return
}

// Verify 认证鉴权
func (x *Service) Verify(ctx context.Context, ts string) (claims passport.Claims, err error) {
	if claims, err = x.Passport.Verify(ts); err != nil {
		return
	}
	var result bool
	// 检测会话
	if result, err = x.SessionsService.Verify(ctx, claims.UserId, claims.ID); err != nil {
		return
	}
	if !result {
		err = errors.NewPublic("会话令牌不一致")
		return
	}

	// TODO: 检查用户状态

	// 会话续约
	if err = x.SessionsService.Renew(ctx, claims.UserId); err != nil {
		return
	}

	return
}

// GetRefreshCode 获取刷新令牌验证码
func (x *Service) GetRefreshCode(ctx context.Context, userId string) (code string, err error) {
	if code, err = gonanoid.Nanoid(); err != nil {
		return
	}
	if err = x.Captcha.Create(ctx, userId, code, 15*time.Second); err != nil {
		return
	}
	return
}

// RefreshToken 刷新令牌
func (x *Service) RefreshToken(ctx context.Context, claims passport.Claims, code string) (ts string, err error) {
	// 验证随机码
	if err = x.Captcha.Verify(ctx, claims.UserId, code); err != nil {
		return
	}
	// 创建令牌
	if ts, err = x.Passport.Create(claims.UserId, claims.ID); err != nil {
		return
	}
	return
}

// Logout 注销登录
func (x *Service) Logout(ctx context.Context, userId string) (err error) {
	return x.SessionsService.Remove(ctx, userId)
}

// // GetOptions 返回通用配置
//	func (x *Service) GetOptions(v string) utils.H {
//		switch v {
//		// 上传类
//		case "upload":
//			switch x.Values.Cloud {
//			// 腾讯云
//			case "tencent":
//				// Cos 对象存储
//				return utils.H{
//					"type": "cos",
//					"url": fmt.Sprintf(`https://%s.cos.%s.myqcloud.com`,
//						x.Values.TencentCosBucket, x.Values.TencentCosRegion,
//					),
//					"limit": x.Values.TencentCosLimit,
//				}
//			}
//		// 企业平台类
//		case "office":
//			switch x.Values.Office {
//			// 飞书
//			case "feishu":
//				return utils.H{
//					"url":      "https://open.feishu.cn/open-apis/authen/v1/index",
//					"redirect": x.Values.RedirectUrl,
//					"app_id":   x.Values.FeishuAppId,
//				}
//			}
//		}
//		return nil
//	}
//

// GetIdentity 获取用户缓存
func (x *Service) GetIdentity(ctx context.Context, userId string) (data model.User, err error) {
	key := x.Values.Name("users", userId)
	var exists int64
	if exists, err = x.Redis.Exists(ctx, key).Result(); err != nil {
		return
	}

	if exists == 0 {
		id, _ := primitive.ObjectIDFromHex(userId)
		option := options.FindOne().SetProjection(bson.M{"password": 0})
		if err = x.Db.Collection("users").
			FindOne(ctx, bson.M{
				"_id":    id,
				"status": true,
			}, option).
			Decode(&data); err != nil {
			return
		}

		var value string
		if value, err = sonic.MarshalString(data); err != nil {
			return
		}

		if err = x.Redis.Set(ctx, key, value, 0).Err(); err != nil {
			return
		}

		return
	}

	var result string
	if result, err = x.Redis.Get(ctx, key).Result(); err != nil {
		return
	}
	if err = sonic.UnmarshalString(result, &data); err != nil {
		return
	}

	return
}

// GetUser 获取登录用户信息
func (x *Service) GetUser(ctx context.Context, userId string) (data map[string]interface{}, err error) {
	var user model.User
	if user, err = x.GetIdentity(ctx, userId); err != nil {
		return
	}

	data = map[string]interface{}{
		"email":  user.Email,
		"name":   user.Name,
		"avatar": user.Avatar,
	}

	return
}

// SetUser 设置登录用户信息
func (x *Service) SetUser(ctx context.Context, userId string, data map[string]interface{}) (result interface{}, err error) {
	id, _ := primitive.ObjectIDFromHex(userId)
	update := bson.M{
		"$set": data,
	}
	if result, err = x.Db.Collection("users").
		UpdateByID(ctx, id, update); err != nil {
		return
	}

	// 用户缓存刷新
	key := x.Values.Name("users", userId)
	if _, err = x.Redis.Del(ctx, key).Result(); err != nil {
		return
	}

	return
}
