package hash

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
	"strconv"
	"strings"
)

var (
	DefaultTime    = uint32(4)
	DefaultMemory  = uint32(64 * 1024)
	DefaultThreads = uint8(1)
)

type Option struct {
	Time    uint32
	Memory  uint32
	Threads uint8
}

func Make(password []byte, option Option) (hash string, err error) {
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return
	}
	if option.Time == 0 {
		option.Time = DefaultTime
	}
	if option.Memory == 0 {
		option.Memory = DefaultMemory
	}
	if option.Threads == 0 {
		option.Threads = DefaultThreads
	}
	key := argon2.IDKey(password, salt, option.Time, option.Memory, option.Threads, 32)
	var build strings.Builder
	build.WriteString("$argon2id$v=" + strconv.Itoa(argon2.Version))
	build.WriteString("$m=" + strconv.Itoa(int(option.Memory)))
	build.WriteString(",t=" + strconv.Itoa(int(option.Time)))
	build.WriteString(",p=" + strconv.Itoa(int(option.Threads)))
	build.WriteString("$" + base64.RawStdEncoding.EncodeToString(salt))
	build.WriteString("$" + base64.RawStdEncoding.EncodeToString(key))
	hash = build.String()
	return
}
