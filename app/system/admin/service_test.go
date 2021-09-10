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
	data, err := s.GetFromCache(context.Background(), "48ee5eda-4117-402a-88db-9d12df42aae9")
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
