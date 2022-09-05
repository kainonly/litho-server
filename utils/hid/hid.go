package hid

import "github.com/speps/go-hashids/v2"

type HID struct {
	HashID *hashids.HashID
}

// NewHID 创建 ID 加密
func NewHID(key string, alphabet string) (x *HID, err error) {
	x = new(HID)
	if x.HashID, err = hashids.NewWithData(&hashids.HashIDData{
		Alphabet: alphabet,
		Salt:     key,
	}); err != nil {
		return
	}
	return
}

// Encode 加密
func (x *HID) Encode(value []int) (string, error) {
	return x.HashID.Encode(value)
}

// Decode 解密
func (x *HID) Decode(value string) ([]int, error) {
	return x.HashID.DecodeWithError(value)
}
