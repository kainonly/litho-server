package controller

import (
	"lab-api/service"
)

type Admin struct {
	*service.Services
}

func NewAdmin(s *service.Services) *Admin {
	return &Admin{s}
}
