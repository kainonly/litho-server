package service

import (
	"context"
	"testing"
)

func TestAcl_GetFromCache(t *testing.T) {
	data, err := s.Acl.GetFromCache(context.Background(), "acl", false)
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestAcl_RefreshCache(t *testing.T) {
	if err := s.Acl.RefreshCache(context.Background()); err != nil {
		t.Error(err)
	}
}
