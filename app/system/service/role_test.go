package service

import (
	"context"
	"testing"
)

func TestRole_GetFromCache(t *testing.T) {
	data, err := s.Role.GetFromCache(context.Background(), 2)
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestRole_RefreshCache(t *testing.T) {
	if err := s.Role.RefreshCache(context.Background()); err != nil {
		t.Error(err)
	}
}
