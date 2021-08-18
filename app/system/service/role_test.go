package service

import (
	"context"
	"testing"
)

//func TestRole_GetFromCache(t *testing.T) {
//	s.Role.GetFromCache()
//}

func TestRole_RefreshCache(t *testing.T) {
	if err := s.Role.RefreshCache(context.Background()); err != nil {
		t.Error(err)
	}
}
