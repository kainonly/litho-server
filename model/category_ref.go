package model

type CategoryRef struct {
	ID         string `gorm:"column:id;type:bigint;primaryKey" json:"id"`
	Type       int16  `gorm:"column:type;type:smallint;not null;index" json:"type"`
	CategoryID string `gorm:"column:category_id;type:bigint;not null;index" json:"category_id"`
	RefID      string `gorm:"column:ref_id;type:bigint;not null;index" json:"ref_id"`
}

func (CategoryRef) TableName() string {
	return "category_ref"
}
