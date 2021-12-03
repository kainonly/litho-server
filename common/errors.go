package common

import "errors"

var (
	LoginExpired = errors.New("login token has expired")
)
