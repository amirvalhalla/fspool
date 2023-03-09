package cfgs

import (
	"github.com/amirvalhalla/fspool/pkg/cfgs"
	fsConfig "github.com/amirvalhalla/fspool/pkg/cfgs/fs"
	"time"
)

/*
* FSPoolConfiguration is a configuration for fs pool
* Perm: will define permission of fspool
* Tip: if you define ROnly all instances which it will generate is ROnly but you can ovveride configuration of filesystem by fs configuration
* Tip: if you define WOnly all instances which it will generate is WOnly but you can ovveride configuration of filesystem by fs configuration
* Tip: if you define RW all instances which it will generate  just have 1 writer and unlimited readers that you can define readerLimit to restrict it
* memoryRent: will get specific amount of memory for reading from file or writing into file to speed up the write or read process (unit is byte)
* limit: limit of getting new filesystem instances
* readerLimit: limit of getting reader instances from each filesystem instance (not fspool)
* flushType: flush type define how to flush into file
* flushDuration: flushing into disk for each instance by timer
* flushSize: flushing into disk for each instance by size (unit is byte)
 */
type FSPoolConfiguration struct {
	Perm          cfgs.FSPerm    //required
	MemoryRent    uint64         //required
	Limit         uint32         //required
	ReaderLimit   uint32         //required
	FlushType     cfgs.FlushType //required
	FlushDuration time.Duration  //required (depends on FlushType)
	FlushSize     uint64         //required  (depends on FlushType)
}

func (c FSPoolConfiguration) MapToFsConfiguration() fsConfig.FSConfiguration {
	return fsConfig.FSConfiguration{
		Perm:          c.Perm,
		MemoryRent:    c.MemoryRent,
		FlushType:     c.FlushType,
		FlushDuration: c.FlushDuration,
		FlushSize:     c.FlushSize,
	}
}
