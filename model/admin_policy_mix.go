package model

type AdminPolicyMix struct {
	AdminId uint64 `json:"-"`
	AclKey  string `json:"-"`
	Policy  *bool  `json:"-"`
}
