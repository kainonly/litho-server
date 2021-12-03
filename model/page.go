package model

type Page struct {
	Parent   string     `bson:"parent" json:"parent"`
	Fragment string     `bson:"fragment" json:"fragment"`
	Name     string     `bson:"name" json:"name"`
	Nav      bool       `bson:"nav" json:"nav"`
	Icon     string     `bson:"icon" json:"icon"`
	Sort     uint8      `bson:"sort" json:"sort"`
	Router   string     `bson:"router" json:"router"`
	Option   PageOption `bson:"option,omitempty" json:"option,omitempty"`
}

type PageOption struct {
	Schema     string       `bson:"schema,omitempty" json:"schema,omitempty"`
	Fetch      bool         `bson:"fetch,omitempty" json:"fetch,omitempty"`
	Fields     []ViewFields `bson:"fields,omitempty" json:"fields,omitempty"`
	Validation string       `bson:"validation,omitempty" json:"validation,omitempty"`
}

type ViewFields struct {
	Key     string `bson:"key" json:"key"`
	Label   string `bson:"label" json:"label"`
	Display bool   `bson:"display" json:"display"`
}
