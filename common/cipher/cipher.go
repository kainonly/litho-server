package cipher

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/chacha20poly1305"
)

type Cipher struct {
	AEAD cipher.AEAD
}

// NewCipher 创建加密
func NewCipher(key string) (x *Cipher, err error) {
	x = new(Cipher)
	if x.AEAD, err = chacha20poly1305.NewX([]byte(key)); err != nil {
		return
	}
	return
}

// Encode 加密
func (x *Cipher) Encode(data []byte) (ciphertext string, err error) {
	nonce := make([]byte, x.AEAD.NonceSize(), x.AEAD.NonceSize()+len(data)+x.AEAD.Overhead())
	rand.Read(nonce)
	encrypted := x.AEAD.Seal(nonce, nonce, data, nil)
	ciphertext = base64.StdEncoding.EncodeToString(encrypted)
	return
}

// Decode 解密
func (x *Cipher) Decode(ciphertext string) (data []byte, err error) {
	encrypted, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return
	}
	nonce, text := encrypted[:x.AEAD.NonceSize()], encrypted[x.AEAD.NonceSize():]
	return x.AEAD.Open(nil, nonce, text, nil)
}
