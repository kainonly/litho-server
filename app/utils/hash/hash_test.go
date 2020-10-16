package hash

import (
	"strconv"
	"testing"
)

var checkHash string

func TestMake(t *testing.T) {
	hash, err := Make(`pass`, Option{})
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
	checkHash = hash
}

func TestCheck(t *testing.T) {
	result, err := Verify(
		`pass`,
		checkHash,
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func BenchmarkMakeAndCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pass := "pass" + strconv.Itoa(i)
		hash, err := Make(pass, Option{})
		if err != nil {
			b.Error(err)
		}
		result, err := Verify(pass, hash)
		if err != nil {
			b.Error(err)
		}
		if result == false {
			b.Error("false")
		}
		b.Log(pass, ":", result)
	}

}
