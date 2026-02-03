package global

import (
	"time"
)

type Bool func() bool
type Uint64 func() uint64
type Duration func() time.Duration

func NewBool(n string, v bool) Bool {
	nameLen = uint32(len(n))
	copy(name[:nameLen], []byte(n))
	val = 0
	if v {
		val = 1
	}
	set()
	return func() bool {
		nameLen = uint32(len(n))
		copy(name[:nameLen], []byte(n))
		get()
		return val == 1
	}
}

func NewUint64(n string, v uint64) Uint64 {
	nameLen = uint32(len(n))
	copy(name[:nameLen], []byte(n))
	val = v
	set()
	return func() uint64 {
		nameLen = uint32(len(n))
		copy(name[:nameLen], []byte(n))
		get()
		return val
	}
}

func NewDuration(n string, v time.Duration) Duration {
	nameLen = uint32(len(n))
	copy(name[:nameLen], []byte(n))
	val = uint64(v)
	set()
	return func() time.Duration {
		nameLen = uint32(len(n))
		copy(name[:nameLen], []byte(n))
		get()
		return time.Duration(int64(val))
	}
}
