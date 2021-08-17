package service

import (
	"context"
	"testing"
)

func TestIndex_GenerateCode(t *testing.T) {
	if err := s.Index.GenerateCode(context.Background(), "test", "abc"); err != nil {
		t.Error(err)
	}
}

func TestIndex_VerifyCode(t *testing.T) {
	result, err := s.Index.VerifyCode(context.Background(), "test", "abc")
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func TestIndex_RemoveCode(t *testing.T) {
	if err := s.Index.RemoveCode(context.Background(), "test"); err != nil {
		t.Error(err)
	}
}
