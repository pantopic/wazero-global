package global

import (
	"unsafe"
)

var (
	val     uint64
	nameCap uint32 = 64
	nameLen uint32
	name    = make([]byte, nameCap)
	meta    = make([]uint32, 4)
)

//export __global
func __global() (res uint32) {
	for i, p := range []unsafe.Pointer{
		unsafe.Pointer(&val),
		unsafe.Pointer(&nameCap),
		unsafe.Pointer(&nameLen),
		unsafe.Pointer(&name[0]),
	} {
		meta[i] = uint32(uintptr(p))
	}
	return uint32(uintptr(unsafe.Pointer(&meta[0])))
}

//go:wasm-module pantopic/wazero-global
//export __global_get
func get()

// Fix for lint rule `unusedfunc`
var _ = __global
