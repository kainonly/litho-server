package values

import (
	"context"
	"testing"
)

func TestService_Get(t *testing.T) {
	data, err := service.Get(context.TODO(), []string{
		"user_session_expire",
		"tencent_secret_key",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestService_Set(t *testing.T) {

}

func TestService_Delete(t *testing.T) {

}
