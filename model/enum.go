package model

type Enum struct {
	// 枚举名称
	Label string `bson:"label"`

	// 枚举数值
	Value interface{} `bson:"value"`
}
