package service

import (
	"context"
	"testing"
)

func TestRole_GetFromCache(t *testing.T) {
	data, err := role.GetFromCache(context.Background(), "admin")
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestRole_RefreshCache(t *testing.T) {
	if err := role.RefreshCache(context.Background()); err != nil {
		t.Error(err)
	}
}
