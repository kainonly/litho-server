package passport

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var passport = &Passport{
	Key: "6ixSiEXaqxsJTozbnxQ76CWdZXB2JazK",
	Option: Option{
		Iss: "weplanx",
		Sub: "system",
		Aud: []string{"api"},
		Nbf: 1,
		Exp: 720,
	},
}

func TestPassport(t *testing.T) {
	jti := uuid.New().String()
	tokenString, err := passport.Create(jti, map[string]interface{}{
		"userId": "xs1fp",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, tokenString)
	assert.Nil(t, err)
	t.Log(tokenString)
	var clamis jwt.MapClaims
	_, err = passport.Verify(tokenString)
	assert.Error(t, err)
	time.Sleep(time.Second)
	clamis, err = passport.Verify(tokenString)
	assert.Nil(t, err)
	assert.Equal(t, "weplanx", clamis["iss"])
	assert.Equal(t, "system", clamis["sub"])
	assert.Equal(t, []interface{}{"api"}, clamis["aud"])
	assert.Equal(t, map[string]interface{}{
		"userId": "xs1fp",
	}, clamis["context"])
	// 使用其他签名的Token进行验证
	_, err = passport.Verify(`eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ`)
	assert.Error(t, err)
}
