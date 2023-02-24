package fs

import (
	cfgs2 "github.com/amirvalhalla/fspool/pkg/cfgs/fs"
	cfgs "github.com/amirvalhalla/fspool/pkg/cfgs/fspool"
	"testing"
)

func TestNewFilesystem(t *testing.T) {
	fsPoolConfig := cfgs.FSPoolConfiguration{}
	fsConfig := cfgs2.FSConfiguration{}
	NewFilesystem("/record/test.txt", fsPoolConfig, fsConfig)
}
