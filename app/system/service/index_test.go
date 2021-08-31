package service

import (
	"context"
	"testing"
)

func TestIndex_GenerateCode(t *testing.T) {
	if err := index.GenerateCode(context.Background(), "test", "abc"); err != nil {
		t.Error(err)
	}
	t.Log("ok")
}

func TestIndex_VerifyCode(t *testing.T) {
	result, err := index.VerifyCode(context.Background(), "test", "abc")
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func TestIndex_RemoveCode(t *testing.T) {
	if err := index.RemoveCode(context.Background(), "test"); err != nil {
		t.Error(err)
	}
}
