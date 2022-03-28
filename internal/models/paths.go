package models

import (
	"echo-starter/internal/wellknown"
)

type Paths struct {
	Home     string
	About    string
	Login    string
	Logout   string
	Deep     string
	Profiles string
}

func NewPaths() *Paths {
	return &Paths{
		Home:     wellknown.HomePath,
		About:    wellknown.AboutPath,
		Login:    wellknown.LoginPath,
		Logout:   wellknown.LogoutPath,
		Deep:     "/deep/a/b",
		Profiles: wellknown.ProfilesPath,
	}
}
