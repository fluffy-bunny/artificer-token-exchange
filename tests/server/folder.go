package server

import "syscall"

type (
	FolderChanger struct {
		originalPath  string
		changedToPath string
	}
)

func NewFolderChanger(target string) *FolderChanger {

	f := &FolderChanger{}
	f.changeTo(target)
	return f
}

func (s *FolderChanger) changeTo(target string) {
	s.changedToPath = target
	s.originalPath, _ = syscall.Getwd()
	syscall.Chdir(s.changedToPath)
}

func (s *FolderChanger) ChangeBack() {
	syscall.Chdir(s.originalPath)
}
