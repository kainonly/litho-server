package tencent_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_GetCity(t *testing.T) {
	dto, err := x.TencentSerice.GetCity("119.41.34.152")
	assert.NoError(t, err)
	t.Log(dto)
}
