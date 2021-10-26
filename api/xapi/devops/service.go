package devops

import (
	"laboratory/common"
)

type Service struct {
	*InjectService
}

type InjectService struct {
	common.App
}
