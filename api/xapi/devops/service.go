package devops

import (
	"laboratory/common"
)

type InjectService struct {
	*common.App
}

type Service struct {
	*InjectService
}

func NewService(i *InjectService) *Service {
	return &Service{
		InjectService: i,
	}
}
