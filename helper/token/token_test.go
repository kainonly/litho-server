package token

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"testing"
)

var tokenString string
var err error

func TestMain(m *testing.M) {
	Key = []byte("hello")
	os.Exit(m.Run())
}

func TestMake(t *testing.T) {
	tokenString, err = Make(jwt.MapClaims{
		"username": "kain",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tokenString)
}

func TestVerify(t *testing.T) {
	result, claims, err := Verify(tokenString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	t.Log(claims)
}
