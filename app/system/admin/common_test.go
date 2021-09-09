package admin

import (
	"lab-api/common"
	"log"
	"os"
	"testing"
)

var s *Service

func TestMain(m *testing.M) {
	os.Chdir("../../../")
	set, err := common.LoadSettings()
	if err != nil {
		log.Fatalln(err)
	}
	db, err := common.InitializeDatabase(set)
	if err != nil {
		log.Fatalln(err)
	}
	redis, err := common.InitializeRedis(set)
	if err != nil {
		log.Fatalln(err)
	}
	cipher, err := common.InitializeCipher(set)
	if err != nil {
		log.Fatalln(err)
	}
	d := common.Dependency{
		Set:    set,
		Db:     db,
		Redis:  redis,
		Cookie: common.InitializeCookie(set),
		Authx:  common.InitializeAuthx(set),
		Cipher: cipher,
	}
	s = NewService(&d)
	os.Exit(m.Run())
}
