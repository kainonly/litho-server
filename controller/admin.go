package controller

import "lab-api/service"

type Admin struct {
	admin *service.Admin
}

func NewAdmin(admin *service.Admin) *Index {
	return &Index{admin}
}
