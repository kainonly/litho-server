package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/errors"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/weplanx/server/api/sessions"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"github.com/weplanx/server/utils/captcha"
	"github.com/weplanx/server/utils/locker"
	"github.com/weplanx/server/utils/passlib"
	"github.com/weplanx/server/utils/passport"
	"go.mongodb.org/mongo-driver/bson"
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
func (x *Service) Login(ctx context.Context, identity string, password string) (ts string, err error) {
	var user model.User
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{
			"status": true,
			"$or": bson.A{
				bson.M{"username": identity},
				bson.M{"email": identity},
			},
		}).
		Decode(&user); err != nil {
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

//
//type Nav struct {
//	ID     primitive.ObjectID `bson:"_id" json:"_id"`
//	Parent interface{}        `json:"parent"`
//	Name   string             `json:"name"`
//	Icon   string             `json:"icon"`
//	Kind   string             `json:"kind"`
//	Sort   int64              `json:"sort"`
//}
//
//// GetNavs 筛选导航数据
//func (x *Service) GetNavs(ctx context.Context, uid string) (navs []Nav, err error) {
//	// TODO: 权限过滤...
//	//var user model.User
//	//if user, err = x.UsersService.GetActived(ctx, uid); err != nil {
//	//	return
//	//}
//	var cursor *mongo.Cursor
//	if cursor, err = x.Db.Collection("pages").
//		Find(ctx, bson.M{"status": true}); err != nil {
//		return
//	}
//	if err = cursor.All(ctx, &navs); err != nil {
//		return
//	}
//	return
//}
//

//// GetOptions 返回通用配置
//func (x *Service) GetOptions(v string) utils.H {
//	switch v {
//	// 上传类
//	case "upload":
//		switch x.Values.GetCloud() {
//		// 腾讯云
//		case "tencent":
//			// Cos 对象存储
//			return utils.H{
//				"type": "cos",
//				"url": fmt.Sprintf(`https://%s.cos.%s.myqcloud.com`,
//					x.Values.GetTencentCosBucket(), x.Values.GetTencentCosRegion(),
//				),
//				"limit": x.Values.GetTencentCosLimit(),
//			}
//		}
//	// 企业平台类
//	case "office":
//		switch x.Values.GetOffice() {
//		// 飞书
//		case "feishu":
//			return utils.H{
//				"url":      "https://open.feishu.cn/open-apis/authen/v1/index",
//				"redirect": x.Values.GetRedirectUrl(),
//				"app_id":   x.Values.GetFeishuAppId(),
//			}
//		}
//	}
//	return nil
//}
//
//// GetActived 获取登录用户数据
//func (x *Service) GetActived(ctx context.Context, id string) (data model.User, err error) {
//	key := x.Values.Name("users")
//	var exists int64
//	if exists, err = x.Redis.Exists(ctx, key).Result(); err != nil {
//		return
//	}
//
//	if exists == 0 {
//		option := options.Find().SetProjection(bson.M{"password": 0})
//		var cursor *mongo.Cursor
//		if cursor, err = x.Db.Collection("users").
//			Find(ctx, bson.M{"status": true}, option); err != nil {
//			return
//		}
//
//		values := make(map[string]string)
//		for cursor.Next(ctx) {
//			var user model.User
//			if err = cursor.Decode(&user); err != nil {
//				return
//			}
//
//			var value string
//			if value, err = sonic.MarshalString(user); err != nil {
//				return
//			}
//
//			values[user.ID.Hex()] = value
//		}
//		if err = cursor.Err(); err != nil {
//			return
//		}
//
//		if err = x.Redis.HSet(ctx, key, values).Err(); err != nil {
//			return
//		}
//	}
//
//	var result string
//	if result, err = x.Redis.HGet(ctx, key, id).Result(); err != nil {
//		return
//	}
//	if err = sonic.UnmarshalString(result, &data); err != nil {
//		return
//	}
//
//	return
//}
//
//// GetUser 获取登录用户信息
//func (x *Service) GetUser(ctx context.Context, uid string) (data map[string]interface{}, err error) {
//	var user model.User
//	if user, err = x.GetActived(ctx, uid); err != nil {
//		return
//	}
//
//	data = map[string]interface{}{
//		"username":    user.Username,
//		"email":       user.Email,
//		"name":        user.Name,
//		"avatar":      user.Avatar,
//		"sessions":    user.Sessions,
//		"last":        user.Last,
//		"create_time": user.CreateTime,
//	}
//
//	// 权限组名称
//	var cursor *mongo.Cursor
//	var roles []string
//	if cursor, err = x.Db.Collection("roles").
//		Find(ctx, bson.M{"_id": bson.M{"$in": user.Roles}}); err != nil {
//		return
//	}
//	for cursor.Next(ctx) {
//		var value model.Role
//		if err = cursor.Decode(&value); err != nil {
//			return
//		}
//
//		roles = append(roles, value.Name)
//	}
//	if err = cursor.Err(); err != nil {
//		return
//	}
//	data["roles"] = roles
//
//	// 部门名称
//	if user.Department != nil {
//		var department model.Department
//		if err = x.Db.Collection("departments").
//			FindOne(ctx, bson.M{"_id": *user.Department}).
//			Decode(&data); err != nil {
//			return
//		}
//		data["department"] = department.Name
//	}
//
//	return
//}
//
//// SetUser 设置登录用户信息
//func (x *Service) SetUser(ctx context.Context, id string, data SetUserDto) (interface{}, error) {
//	oid, _ := primitive.ObjectIDFromHex(id)
//	update := bson.M{
//		"$set": data,
//	}
//	if data.Reset != "" {
//		update["$unset"] = bson.M{data.Reset: ""}
//	}
//
//	return x.Db.Collection("users").
//		UpdateByID(ctx, oid, update)
//}
