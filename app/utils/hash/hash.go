package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/argon2"
	"regexp"
	"strconv"
	"strings"
)

var (
	DefaultTime    = uint32(4)
	DefaultMemory  = uint32(64 * 1024)
	DefaultThreads = uint8(1)
)

var (
	ErrInvalidHash         = errors.New(`the encoded hash is not in the correct format`)
	ErrIncompatibleVersion = errors.New(`incompatible version of argon2`)
)

type Option struct {
	Time    uint32
	Memory  uint32
	Threads uint8
}

func Make(password string, option Option) (hashedPassword string, err error) {
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
	hash := argon2.IDKey([]byte(password), salt, option.Time, option.Memory, option.Threads, 32)
	var build strings.Builder
	build.WriteString("$argon2id$v=" + strconv.Itoa(argon2.Version))
	build.WriteString("$m=" + strconv.Itoa(int(option.Memory)))
	build.WriteString(",t=" + strconv.Itoa(int(option.Time)))
	build.WriteString(",p=" + strconv.Itoa(int(option.Threads)))
	build.WriteString("$" + base64.RawStdEncoding.EncodeToString(salt))
	build.WriteString("$" + base64.RawStdEncoding.EncodeToString(hash))
	hashedPassword = build.String()
	return
}

func Check(password string, hashedPassword string) (result bool, err error) {
	args := regexp.
		MustCompile(`^\$(\w+)\$v=(\d+)\$m=(\d+),t=(\d+),p=(\d+)\$(\S+)\$(\S+)`).
		FindStringSubmatch(hashedPassword)
	if len(args) != 8 {
		return false, ErrInvalidHash
	}
	if args[1] != `argon2id` {
		return false, ErrInvalidHash
	}
	version, err := strconv.Atoi(args[2])
	if err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, ErrIncompatibleVersion
	}
	memory, err := strconv.ParseUint(args[3], 10, 32)
	if err != nil {
		return false, err
	}
	time, err := strconv.ParseUint(args[4], 10, 32)
	if err != nil {
		return false, err
	}
	threads, err := strconv.Atoi(args[5])
	if err != nil {
		return false, err
	}
	option := Option{
		Memory:  uint32(memory),
		Time:    uint32(time),
		Threads: uint8(threads),
	}
	decodeSalt, err := base64.RawStdEncoding.DecodeString(args[6])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(args[7])
	if err != nil {
		return false, err
	}
	newHash := argon2.IDKey([]byte(password), decodeSalt, option.Time, option.Memory, option.Threads, 32)
	if subtle.ConstantTimeCompare(hash, newHash) == 1 {
		return true, nil
	}
	return false, nil
}
