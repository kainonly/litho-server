package model

type Role struct {
	Common

	Resources []string `json:"resources"`
}
