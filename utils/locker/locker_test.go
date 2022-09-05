package locker_test

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/utils/locker"
	"os"
	"testing"
	"time"
)

var x *locker.Locker

func TestMain(m *testing.M) {
	opts, err := redis.ParseURL(os.Getenv("REDIS_URI"))
	if err != nil {
		return
	}

	x = &locker.Locker{
		Values: &common.Values{
			App: common.App{Namespace: "dev"},
		},
		Redis: redis.NewClient(opts),
	}
	os.Exit(m.Run())
}

func TestLocker_Update(t *testing.T) {
	var err error
	err = x.Update(context.TODO(), "dev", time.Second*60)
	assert.NoError(t, err)
	var ttl time.Duration
	ttl, err = x.Redis.TTL(context.TODO(), x.Key("dev")).Result()
	assert.NoError(t, err)
	t.Log(ttl.Seconds())
}

func TestLocker_Verify(t *testing.T) {
	var err error
	var result bool
	result, err = x.Verify(context.TODO(), "dev", 3)
	assert.NoError(t, err)
	assert.False(t, result)

	for i := 0; i < 3; i++ {
		err = x.Update(context.TODO(), "dev", time.Second*60)
		assert.NoError(t, err)
	}

	result, err = x.Verify(context.TODO(), "dev", 3)
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestLocker_Delete(t *testing.T) {
	var err error
	err = x.Delete(context.TODO(), "dev")
	assert.NoError(t, err)

	var exists int64
	exists, err = x.Redis.Exists(context.TODO(), x.Key("dev")).Result()
	assert.NoError(t, err)
	assert.Equal(t, exists, int64(0))
}
