package token

import (
	"log"
	"os"
	"testing"
	"van-api/bootstrap"
)

var token []byte
var err error

func TestMain(m *testing.M) {
	os.Chdir("../..")
	cfg, err := bootstrap.LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}
	Key = []byte(cfg.App.Key)
	Options = cfg.Token
	os.Exit(m.Run())
}

func TestMake(t *testing.T) {
	token, err = Make("system", map[string]interface{}{
		"username": "kain",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(token))
}

func TestVerify(t *testing.T) {
	claims, err := Verify("system", token, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claims)
}
