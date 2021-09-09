package role

import (
	"context"
	"testing"
)

func TestService_GetFromCache(t *testing.T) {
	data, err := s.GetFromCache(context.TODO(), "admin")
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestService_RefreshCache(t *testing.T) {
	if err := s.RefreshCache(context.TODO()); err != nil {
		t.Error(err)
	}
}

func TestService_RemoveCache(t *testing.T) {

}
