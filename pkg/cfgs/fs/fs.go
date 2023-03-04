package cfgs

import (
	"github.com/amirvalhalla/fspool/pkg/cfgs"
	"time"
)

/*
* FSConfiguration is a configuration for each reader or writer instance of you will get from existing fs pool
* Perm: will define permission of fs
* memoryRent: will get specific amount of memory for reading from file or writing into file to speed up the write or read process (unit is byte)
* readerLimit: limit of getting reader instances from each filesystem instance (not fspool)
* flushType: flush type define how to flush into file
* flushDuration: flushing into disk for each instance by timer
* flushSize: flushing into disk for each instance by size (unit is byte)
 */
type FSConfiguration struct {
	Perm          cfgs.FSPerm
	MemoryRent    uint64
	ReaderLimit   uint32
	FlushType     cfgs.FlushType
	FlushDuration time.Duration //depends on FlushType
	FlushSize     uint64        //depends on FlushType
}
