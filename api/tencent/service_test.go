package tencent_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_GetIpv4(t *testing.T) {
	ctx := context.TODO()
	dto, err := x.TencentService.GetIpv4(ctx, "119.41.34.152")
	assert.NoError(t, err)
	t.Log(dto)
}

func TestService_GetIpv6(t *testing.T) {
	ctx := context.TODO()
	dto, err := x.TencentService.GetIpv6(ctx, "240e:314:e441:9000:2d47:2c35:4fb:a883")
	assert.NoError(t, err)
	t.Log(dto)
}
