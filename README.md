# Wazero Global

A [wazero](https://pkg.go.dev/github.com/tetratelabs/wazero) host module, ABI and guest SDK providing globals for WASI modules.

## Host Module

[![Go Reference](https://godoc.org/github.com/pantopic/wazero-global/host?status.svg)](https://godoc.org/github.com/pantopic/wazero-global/host)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-global/host)](https://goreportcard.com/report/github.com/pantopic/wazero-global/host)
[![Go Coverage](https://github.com/pantopic/wazero-global/wiki/host/coverage.svg)](https://raw.githack.com/wiki/pantopic/wazero-global/host/coverage.html)

First register the host module with the runtime

```go
import (
    "github.com/tetratelabs/wazero"
    "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

    "github.com/pantopic/wazero-global/host"
)

func main() {
    ctx := context.Background()
    r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())
    wasi_snapshot_preview1.MustInstantiate(ctx, r)

    module := wazero_global.New()
    module.Register(ctx, r)

    // ...
}
```

## Guest SDK (Go)

[![Go Reference](https://godoc.org/github.com/pantopic/wazero-global/sdk-go?status.svg)](https://godoc.org/github.com/pantopic/wazero-global/sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-global/sdk-go)](https://goreportcard.com/report/github.com/pantopic/wazero-global/sdk-go)

Then you can import the guest SDK into your WASI module to retrieve globals from the host.

```go
package main

import (
    "github.com/pantopic/wazero-global/sdk-go"
)

var (
	testBool     global.Bool
	testUint64   global.Uint64
	testDuration global.Duration
)

func main() {
	testBool     = global.NewBool(`TEST_BOOL`, true)
	testUint64   = global.NewUint64(`TEST_UINT64`, 42)
	testDuration = global.NewDuration(`TEST_DURATION`, time.Minute)
}

//export test
func test() {
    println(testBool())     // true
    println(testUint64())   // 42
    println(testDuration()) // 1m
}
```

The host application can then set and delete globals (as uint64)

```go
func main() {
    // ...

    module.Set(`TEST_BOOL`, 0)
    module.Set(`TEST_UINT64`, 43)
    module.Set(`TEST_DURATION`, uint64(time.Second))
    
    // ...
    
    module.Del(`TEST_BOOL`)
    module.Del(`TEST_UINT64`)
    module.Del(`TEST_DURATION`)
}
```

## Roadmap

This project is in alpha. Breaking API changes should be expected until Beta.

- `v0.0.x` - Alpha
  - [ ] Stabilize API
- `v0.x.x` - Beta
  - [ ] Finalize API
  - [ ] Test in production
- `v1.x.x` - General Availability
  - [ ] Proven long term stability in production
