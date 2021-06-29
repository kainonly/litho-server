package controller

import (
	"github.com/kainonly/gin-planx"
)

type Acl struct {
	*planx.Crud
}

func NewAcl(planx *planx.Planx) *Acl {
	return &Acl{Crud: planx.Make()}
}
