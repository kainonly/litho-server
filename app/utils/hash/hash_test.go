package hash

import "testing"

func TestMake(t *testing.T) {
	hash, err := Make([]byte("pass"), Option{})
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
}
