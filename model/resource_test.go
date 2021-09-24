package model

import "testing"

func TestGenerateResources(t *testing.T) {
	if err := GenerateResources(tx); err != nil {
		t.Error(err)
	}
}
