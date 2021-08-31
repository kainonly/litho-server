package controller

import "github.com/kainonly/go-bit/authx"

type Index struct {
	*Dependency
	auth *authx.Auth
}

func NewIndex(d *Dependency, authx *authx.Authx) *Index {
	return &Index{
		Dependency: d,
		auth:       authx.Make("system"),
	}
}
