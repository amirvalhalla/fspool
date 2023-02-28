// Package cfgs contains configurations of filesystem , filesystem pool and all others configurations
package cfgs

type FSPerm uint8

const (
	ROnly FSPerm = 0
	WOnly FSPerm = 1
	RW    FSPerm = 3
)
