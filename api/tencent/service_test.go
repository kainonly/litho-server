package tencent_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_GetCity(t *testing.T) {
	ctx := context.TODO()
	dto, err := x.TencentService.GetCity(ctx, "119.41.34.152")
	assert.NoError(t, err)
	t.Log(dto)
}
