package model

type Admin struct {
	ID         uint64
	Username   string
	Password   string
	Role       string
	Call       string
	Email      string
	Phone      string
	Avatar     string
	Status     bool
	CreateTime uint64
	UpdateTime uint64
}
