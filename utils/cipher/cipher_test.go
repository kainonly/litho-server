package cipher_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/utils/cipher"
	"testing"
)

var x1 *cipher.Cipher
var x2 *cipher.Cipher

func TestUseCipher(t *testing.T) {
	var err error
	_, err = cipher.NewCipher("123456")
	assert.Error(t, err)
	// 必须是32位的密钥
	x1, err = cipher.NewCipher("6ixSiEXaqxsJTozbnxQ76CWdZXB2JazK")
	assert.Nil(t, err)
	x2, err = cipher.NewCipher("74rILbVooYLirHrQJcslHEAvKZI7PKF9")
	assert.Nil(t, err)
}

var text = "Gophers, gophers, gophers everywhere!"
var encryptedText string

func TestCipher_Encode(t *testing.T) {
	var err error
	encryptedText, err = x1.Encode([]byte(text))
	t.Log(encryptedText)
	assert.Nil(t, err)
}

func TestCipher_Decode(t *testing.T) {
	decryptedText, err := x1.Decode(encryptedText)
	assert.Nil(t, err)
	assert.Equal(t, text, string(decryptedText))
	// 不同密钥发生错误
	_, err = x2.Decode(encryptedText)
	assert.Error(t, err)
	// 非Base64发生错误
	_, err = x1.Decode("asdasdasd")
	assert.Error(t, err)
}
