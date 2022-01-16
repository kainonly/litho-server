package model

func Bool(v bool) *bool {
	return &v
}

type Value struct {
	// 名称
	Label string `bson:"label" json:"label"`

	// 数值
	Value interface{} `bson:"value" json:"value"`
}
