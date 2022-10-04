package passlib_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/utils/passlib"
	"testing"
)

func TestPassword(t *testing.T) {
	hash, err := passlib.Hash("pass@VAN1234")
	assert.Nil(t, err)
	t.Log(hash)
	match, err := passlib.Verify("pass@VAN1234", "asdaqweqwexcxzcqweqw")
	assert.Error(t, err)
	assert.False(t, match)
	match, err = passlib.Verify("pass@VAN1235", hash)
	assert.Nil(t, err)
	assert.False(t, match)
	match, err = passlib.Verify("pass@VAN1234", hash)
	assert.Nil(t, err)
	assert.True(t, match)
}
