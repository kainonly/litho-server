package hid_test

import (
	"github.com/speps/go-hashids/v2"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/api/common/hid"
	"testing"
)

var x *hid.HID

func TestNewHID(t *testing.T) {
	var err error
	x, err = hid.NewHID("6ixSiEXaqxsJTozbnxQ76CWdZXB2JazK", hashids.DefaultAlphabet)
	assert.Nil(t, err)
}

var hash string
var value = []int{651}

func TestHID_Encode(t *testing.T) {
	var err error
	hash, err = x.Encode(value)
	assert.Nil(t, err)
}

func TestHID_Decode(t *testing.T) {
	data, err := x.Decode(hash)
	assert.Nil(t, err)
	assert.Equal(t, value, data)
}
