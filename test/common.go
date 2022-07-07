package test

import (
	"server/common"
)

func Bed() (*common.Inject, error) {
	values, err := common.SetValues("config/config.yml")
	if err != nil {
		return nil, err
	}
	return Injectable(values)
}
