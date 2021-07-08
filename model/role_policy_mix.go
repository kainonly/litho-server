package model

type RolePolicyMix struct {
	RoleKey string `json:"-"`
	AclKey  string `json:"-"`
	Policy  *bool  `json:"-"`
}
