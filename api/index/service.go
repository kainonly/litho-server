package index

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/weplanx/server/api/users"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"github.com/weplanx/server/utils/passlib"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	Values *common.Values
	Mongo  *mongo.Client
	Db     *mongo.Database
	Redis  *redis.Client
	Users  *users.Service
}

func (x *Service) Index() time.Time {
	return time.Now()
}

// Login 登录
func (x *Service) Login(ctx context.Context, identity string, password string) (user model.User, err error) {
	if user, err = x.Users.FindByIdentity(ctx, identity); err != nil {
		return
	}

	// 验证密码正确性
	if err = passlib.Verify(password, user.Password); err != nil {
		return
	}

	return
}
