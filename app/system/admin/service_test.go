package admin

import (
	"context"
	"testing"
)

func TestService_FindByUsername(t *testing.T) {
	data, err := s.FindByUsername(context.TODO(), "admin")
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestService_GetFromCache(t *testing.T) {
	data, err := s.GetFromCache(context.Background(), "426597eb-fd6a-4ea4-a734-da248361eb9b")
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
	if err := s.RemoveCache(context.TODO()); err != nil {
		t.Error(err)
	}
}
