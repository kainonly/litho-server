package roles

import "api/common"

type Service struct {
	*InjectService
}

type InjectService struct {
	common.Inject
}
