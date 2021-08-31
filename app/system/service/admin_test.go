package service

import (
	"context"
	"testing"
)

func TestAdmin_FindByUsername(t *testing.T) {
	data, err := admin.FindByUsername(context.Background(), "kain")
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestAdmin_GetFromCache(t *testing.T) {
	data, err := admin.GetFromCache(context.Background(), "6ee951e5-f353-4111-abe3-d5a90fa9e574")
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestAdmin_RefreshCache(t *testing.T) {
	if err := admin.RefreshCache(context.Background()); err != nil {
		t.Error(err)
	}
}

func TestAdmin_RemoveCache(t *testing.T) {
	if err := admin.RemoveCache(context.Background()); err != nil {
		t.Error(err)
	}
}
