package devops

import "api/common"

type Service struct {
	*InjectService
}

type InjectService struct {
	common.App
}

func NewService(i *InjectService) *Service {
	return &Service{
		InjectService: i,
	}
}
