package model

import "testing"

func TestGenerateSchema(t *testing.T) {
	if err := GenerateSchema(tx); err != nil {
		t.Error(err)
	}
}
