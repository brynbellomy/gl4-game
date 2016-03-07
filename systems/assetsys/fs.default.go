package assetsys

import (
	"os"
	"path"
)

type DefaultFilesystem struct {
	assetRoot string
}

func NewDefaultFilesystem(assetRoot string) *DefaultFilesystem {
	return &DefaultFilesystem{assetRoot}
}

func (fs *DefaultFilesystem) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(path.Join(fs.assetRoot, name), flag, perm)
}
