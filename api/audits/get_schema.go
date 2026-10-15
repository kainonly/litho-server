package audits

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

type GetSchemaDto struct {
	Scopes string `query:"scopes" vd:"required"`
}

type GetSchemaResult struct {
	Schemas []Enum            `json:"schemas"`
	Acts    map[string][]Enum `json:"acts"`
}

var auditOnlyActions = map[string][]Enum{
	"users": {
		{Label: "修改账号信息", Value: "set_user"},
		{Label: "修改手机号", Value: "set_phone"},
		{Label: "修改密码", Value: "set_password"},
		{Label: "解除账号设置", Value: "unset"},
	},
}

func (x *Controller) GetSchema(ctx context.Context, c *app.RequestContext) {
	var dto GetSchemaDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, x.AuditsX.GetSchema(dto))
}

func (x *Service) GetSchema(dto GetSchemaDto) GetSchemaResult {
	result := GetSchemaResult{
		Schemas: make([]Enum, 0),
		Acts:    make(map[string][]Enum),
	}
	scopes := make(map[string]struct{})
	if dto.Scopes != "all" {
		for scope := range strings.SplitSeq(dto.Scopes, ",") {
			scope = strings.TrimSpace(scope)
			if scope != "" {
				scopes[scope] = struct{}{}
			}
		}
	}
	for _, resource := range x.V.Resources.List() {
		result.Schemas = append(result.Schemas, Enum{Label: resource.Name, Value: resource.Key})
		if dto.Scopes != "all" {
			if _, ok := scopes[resource.Key]; !ok {
				continue
			}
		}
		actions := make([]Enum, 0, len(resource.Actions))
		for _, action := range resource.Actions {
			actions = append(actions, Enum{Label: action.Name, Value: action.Path})
		}
		actions = append(actions, auditOnlyActions[resource.Key]...)
		result.Acts[resource.Key] = actions
	}
	return result
}
