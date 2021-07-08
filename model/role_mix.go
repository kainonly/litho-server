package model

type RoleMix struct {
	ID         uint64 `json:"_id"`
	Key        string `json:"key"`
	Name       Object `json:"name"`
	Resource   string `json:"resource"`
	Acl        string `json:"acl"`
	Permission string `json:"permission"`
	Note       string `json:"note"`
	Status     bool   `json:"status"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}
