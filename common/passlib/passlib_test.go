package passlib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPassword(t *testing.T) {
	hash, err := Hash("pass@VAN1234")
	assert.Nil(t, err)
	t.Log(hash)
	err = Verify("pass@VAN1234", "asdaqweqwexcxzcqweqw")
	assert.Error(t, err)
	err = Verify("pass@VAN1235", hash)
	assert.Equal(t, ErrNotMatch, err)
	err = Verify("pass@VAN1234", hash)
	assert.Nil(t, err)
}
