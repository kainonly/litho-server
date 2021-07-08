package controller

import (
	"github.com/gin-gonic/gin"
	bit "github.com/kainonly/gin-bit"
	"lab-api/model"
)

type Acl struct {
	*bit.Crud
}

func NewAcl(b *bit.Bit) *Acl {
	return &Acl{
		Crud: b.Crud(model.Acl{}),
	}
}

func (x *Acl) Add(c *gin.Context) interface{} {
	var body struct {
		Name  model.JSONObject `json:"name" binding:"required"`
		Key   string           `json:"key" binding:"required"`
		Write model.Array      `json:"write"`
		Read  model.Array      `json:"read"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	bit.Mixed(c,
		bit.SetBody(&body),
		bit.SetData(&model.Acl{
			Name:   body.Name,
			Key:    body.Key,
			Write:  body.Write,
			Read:   body.Read,
			Status: model.True(),
		}),
	)
	return x.Crud.Add(c)
}
