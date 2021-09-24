package resource

import (
	"lab-api/common"
)

type Service struct {
	*InjectService
}

type InjectService struct {
	common.App
}
