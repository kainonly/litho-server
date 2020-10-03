package model

type Acl struct {
	ID         uint64
	AccessKey  string `gorm:"size:200;unique;not null;comment:api access control key"`
	Name       I18n   `gorm:"type:json;not null;comment:access control name"`
	ReadList   string `gorm:"comment:list of readable api"`
	WriteList  string `gorm:"comment:list of writable api"`
	Status     bool   `gorm:"type:tinyint(1) unsigned;not null;default:1"`
	CreateTime uint64 `gorm:"not null;default:0;autoCreateTime"`
	UpdateTime uint64 `gorm:"not null;default:0;autoUpdateTime"`
}
