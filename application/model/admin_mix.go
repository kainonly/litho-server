package model

type AdminMix struct {
	ID         uint64
	Username   string
	Password   string
	Role       string
	Resource   string
	Acl        string
	Permission string
	Call       string
	Email      string
	Phone      string
	Avatar     string
	Status     bool
	CreateTime uint64
	UpdateTime uint64
}
