package configs

import "time"

type FSPoolFlushType uint8

const (
	FlushBySize FSPoolFlushType = 0
	FlushByTime FSPoolFlushType = 1
)

/*
* FSPoolConfiguration is a configuration for fs pool
* memoryRent: will get specific amount of memory for reading from file or writing into file to speed up the write or read process (unit is byte)
* limit: limit of getting new filesystem instances
* readerLimit: limit of getting reader instances from each filesystem instance (not fspool)
* flushType: flush type define how to flush into file
* flushDuration: flushing into disk for each instance by timer
* flushSize: flushing into disk for each instance by size (unit is byte)
 */
type FSPoolConfiguration struct {
	memoryRent    uint64
	limit         uint32
	readerLimit   uint32
	flushType     FSPoolFlushType
	flushDuration time.Duration
	flushSize     uint64
}
