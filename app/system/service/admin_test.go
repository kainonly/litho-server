package service

import (
	"context"
	"testing"
)

func TestAdmin_RefreshCache(t *testing.T) {
	if err := s.Admin.RefreshCache(context.Background()); err != nil {
		t.Error(err)
	}
}
