package passlib

import (
	"github.com/alexedwards/argon2id"
)

// Hash  创建密码散列
func Hash(password string) (string, error) {
	return argon2id.CreateHash(password, &argon2id.Params{
		Memory:      65536,
		Iterations:  4,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})
}

// Verify 验证密码是否和散列值匹配
func Verify(password string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
