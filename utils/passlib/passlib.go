package passlib

import (
	"errors"
	"github.com/alexedwards/argon2id"
)

var (
	Params = &argon2id.Params{
		Memory:      65536,
		Iterations:  4,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}
	ErrNotMatch = errors.New("密码校验不匹配")
)

// Hash  创建密码散列
func Hash(password string) (string, error) {
	return argon2id.CreateHash(password, Params)
}

// Verify 验证密码是否和散列值匹配
func Verify(password string, hash string) error {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return err
	}
	if !match {
		return ErrNotMatch
	}
	return nil
}
