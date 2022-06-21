package values

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Set(t *testing.T) {
	err := service.Set(context.TODO(), map[string]interface{}{
		"awesome": "test",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestService_Get(t *testing.T) {
	data, err := service.Get(context.TODO(), []string{
		"awesome",
		"tencent_secret_key",
	})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "test", data["awesome"])
	assert.Equal(t, "*", data["tencent_secret_key"])
}

func TestService_Delete(t *testing.T) {
	err := service.Delete(context.TODO(), "awesome")
	if err != nil {
		t.Error(err)
	}
}
