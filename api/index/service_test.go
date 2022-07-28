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
