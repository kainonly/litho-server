package hash

import "testing"

func TestMake(t *testing.T) {
	hash, err := Make(`pass`, Option{})
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
}

func TestCheck(t *testing.T) {
	result := Check(
		`pass`,
		`$argon2id$v=19$m=65536,t=4,p=1$HDPsadBAbRx8MV9opmyZ2A$UkhI1agUWgpGCbcofy6n17xdoSLIh0wu9HOXRZMWkhE`,
	)
	t.Log(result)
}
