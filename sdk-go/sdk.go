package global

import (
	"time"
)

type Bool func() bool
type Uint64 func() uint64
type Duration func() time.Duration

func NewBool(n string, v bool) Bool {
	return func() bool {
		val = 0
		if v {
			val = 1
		}
		nameLen = uint32(len(n))
		copy(name[:nameLen], []byte(n))
		get()
		return val == 1
	}
}

func NewUint64(n string, v uint64) Uint64 {
	return func() uint64 {
		val = v
		nameLen = uint32(len(n))
		copy(name[:nameLen], []byte(n))
		get()
		return val
	}
}

func NewDuration(n string, v time.Duration) Duration {
	return func() time.Duration {
		val = uint64(v)
		nameLen = uint32(len(n))
		copy(name[:nameLen], []byte(n))
		get()
		return time.Duration(int64(val))
	}
}
