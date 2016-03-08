package assetsys

import "os"

type (
	System struct {
		fs IFilesystem
	}

	IFilesystem interface {
		OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
		Subdir(name string) (IFilesystem, error)
	}
)

func New(fs IFilesystem) *System {
	return &System{fs: fs}
}

func (s *System) Filesystem() IFilesystem {
	return s.fs
}

func (s *System) SetFilesystem(fs IFilesystem) {
	s.fs = fs
}
