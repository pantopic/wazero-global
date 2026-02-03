package main

import (
	"time"

	"github.com/pantopic/wazero-global/sdk-go"
)

var (
	testBool     global.Bool
	testUint64   global.Uint64
	testDuration global.Duration
)

func main() {
	testBool = global.NewBool(`TEST_BOOL`, true)
	testUint64 = global.NewUint64(`TEST_UINT64`, 42)
	testDuration = global.NewDuration(`TEST_DURATION`, time.Minute)
}

//export testGet
func testGet() {
	if !testBool() {
		panic(`Get bool failed`)
	}
	if testUint64() != 42 {
		panic(`Get uint64 failed`)
	}
	if testDuration() != time.Minute {
		panic(`Get duration failed`)
	}
}

// Fix for lint rule `unusedfunc`
var _ = testGet
