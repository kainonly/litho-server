package captcha_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common/captcha"
	"os"
	"testing"
	"time"
)

var x *captcha.Captcha

func TestMain(m *testing.M) {
	os.Chdir("../../")
	values, _ := bootstrap.LoadStaticValues()
	redis, _ := bootstrap.UseRedis(values)
	x = &captcha.Captcha{
		Values: values,
		Redis:  redis,
	}
	os.Exit(m.Run())
}

func TestCaptcha_Create(t *testing.T) {
	var err error
	err = x.Create(context.TODO(), "dev1", "abcd", time.Second*60)
	assert.NoError(t, err)
	var ttl time.Duration
	ttl, err = x.Redis.TTL(context.TODO(), x.Values.Key("captcha", "dev1")).Result()
	assert.NoError(t, err)
	t.Log(ttl.Seconds())
	err = x.Create(context.TODO(), "dev2", "abcd", time.Millisecond)
	assert.NoError(t, err)
}

func TestCaptcha_Verify(t *testing.T) {
	var err error
	err = x.Verify(context.TODO(), "dev1", "abc")
	assert.ErrorIs(t, err, captcha.ErrCaptchaInconsistent)
	err = x.Verify(context.TODO(), "dev1", "abcd")
	assert.NoError(t, err)
	err = x.Verify(context.TODO(), "dev2", "abcd")
	assert.ErrorIs(t, err, captcha.ErrCaptchaNotExists)
}

func TestCaptcha_Delete(t *testing.T) {
	var err error
	var exists bool
	exists, err = x.Exists(context.TODO(), "dev1")
	assert.NoError(t, err)
	assert.True(t, exists)
	err = x.Delete(context.TODO(), "dev1")
	assert.NoError(t, err)
	exists, err = x.Exists(context.TODO(), "dev1")
	assert.NoError(t, err)
	assert.False(t, exists)
}
