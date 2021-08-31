package service

import (
	"context"
	"testing"
)

func TestResource_GetFromCache(t *testing.T) {
	data, err := resource.GetFromCache(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestResource_RefreshCache(t *testing.T) {
	if err := resource.RefreshCache(context.Background()); err != nil {
		t.Error(err)
	}
}

func TestResource_RemoveCache(t *testing.T) {
	if err := resource.RemoveCache(context.Background()); err != nil {
		t.Error(err)
	}
}
