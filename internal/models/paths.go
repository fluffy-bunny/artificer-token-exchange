package models

import (
	"echo-starter/internal/wellknown"
)

type Paths struct {
	Home   string
	Login  string
	Logout string
	Deep   string
}

func NewPaths() *Paths {
	return &Paths{
		Home:   wellknown.HomePath,
		Login:  wellknown.LoginPath,
		Logout: wellknown.LogoutPath,
		Deep:   "/deep/a/b",
	}
}
