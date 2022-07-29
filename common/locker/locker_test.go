package locker_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common/locker"
	"os"
	"testing"
	"time"
)

var x *locker.Locker

func TestMain(m *testing.M) {
	os.Chdir("../../")
	values, _ := bootstrap.LoadStaticValues()
	redis, _ := bootstrap.UseRedis(values)
	x = &locker.Locker{
		Values: values,
		Redis:  redis,
	}
	os.Exit(m.Run())
}

func TestLocker_Update(t *testing.T) {
	var err error
	err = x.Update(context.TODO(), "dev", time.Second*60)
	assert.NoError(t, err)
	var ttl time.Duration
	ttl, err = x.Redis.TTL(context.TODO(), x.Values.Key("locker", "dev")).Result()
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
	exists, err = x.Redis.Exists(context.TODO(), x.Values.Key("locker", "dev")).Result()
	assert.NoError(t, err)
	assert.Equal(t, exists, int64(0))
}
