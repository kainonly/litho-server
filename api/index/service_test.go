package index_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/api/index"
	"testing"
	"time"
)

func TestService_CreateCaptcha(t *testing.T) {
	var err error
	err = x.IndexService.CreateCaptcha(context.TODO(), "dev1", "abcd", time.Second*60)
	assert.NoError(t, err)
	var ttl time.Duration
	ttl, err = x.Redis.TTL(context.TODO(), x.Values.Key("captcha", "dev1")).Result()
	assert.NoError(t, err)
	t.Log(ttl.Seconds())
	err = x.IndexService.CreateCaptcha(context.TODO(), "dev2", "abcd", time.Millisecond)
	assert.NoError(t, err)
}

func TestService_VerifyCaptcha(t *testing.T) {
	var err error
	err = x.IndexService.VerifyCaptcha(context.TODO(), "dev1", "abc")
	assert.ErrorIs(t, err, index.ErrCaptchaInconsistent)
	err = x.IndexService.VerifyCaptcha(context.TODO(), "dev1", "abcd")
	assert.NoError(t, err)
	err = x.IndexService.VerifyCaptcha(context.TODO(), "dev2", "abcd")
	assert.ErrorIs(t, err, index.ErrCaptchaNotExists)
}

func TestService_DeleteCaptcha(t *testing.T) {
	var err error
	var exists bool
	exists, err = x.IndexService.ExistsCaptcha(context.TODO(), "dev1")
	assert.NoError(t, err)
	assert.True(t, exists)
	err = x.IndexService.DeleteCaptcha(context.TODO(), "dev1")
	assert.NoError(t, err)
	exists, err = x.IndexService.ExistsCaptcha(context.TODO(), "dev1")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestService_UpdateLocker(t *testing.T) {
	var err error
	err = x.IndexService.UpdateLocker(context.TODO(), "dev", time.Second*60)
	assert.NoError(t, err)
	var ttl time.Duration
	ttl, err = x.Redis.TTL(context.TODO(), x.Values.Key("locker", "dev")).Result()
	assert.NoError(t, err)
	t.Log(ttl.Seconds())
}

func TestService_VerifyLocker(t *testing.T) {
	var err error
	var result bool
	result, err = x.IndexService.VerifyLocker(context.TODO(), "dev", 3)
	assert.NoError(t, err)
	assert.False(t, result)

	for i := 0; i < 3; i++ {
		err = x.IndexService.UpdateLocker(context.TODO(), "dev", time.Second*60)
		assert.NoError(t, err)
	}

	result, err = x.IndexService.VerifyLocker(context.TODO(), "dev", 3)
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestService_DeleteLocker(t *testing.T) {
	var err error
	err = x.IndexService.DeleteLocker(context.TODO(), "dev")
	assert.NoError(t, err)

	var exists int64
	exists, err = x.Redis.Exists(context.TODO(), x.Values.Key("locker", "dev")).Result()
	assert.NoError(t, err)
	assert.Equal(t, exists, int64(0))
}
