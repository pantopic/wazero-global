package wazero_global

import (
	"context"
	_ "embed"
	"os"
	"testing"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed test\.wasm
var testwasm []byte

func TestModule(t *testing.T) {
	var (
		ctx = context.Background()
	)
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	hostModule := New()
	hostModule.Register(ctx, r)

	compiled, err := r.CompileModule(ctx, testwasm)
	if err != nil {
		panic(err)
	}
	cfg := wazero.NewModuleConfig().WithStdout(os.Stdout).WithName(`mod1`)
	mod, err := r.InstantiateModule(ctx, compiled, cfg)
	if err != nil {
		t.Errorf(`%v`, err)
		return
	}
	ctx, err = hostModule.InitContext(ctx, mod)
	if err != nil {
		t.Fatalf(`%v`, err)
	}

	t.Run(`get`, func(t *testing.T) {
		_, err := mod.ExportedFunction(`testGet`).Call(ctx)
		if err != nil {
			panic(err.Error())
		}
	})
	t.Run(`set`, func(t *testing.T) {
		hostModule.Set(`TEST_BOOL`, 0)
		hostModule.Set(`TEST_UINT64`, 43)
		hostModule.Set(`TEST_DURATION`, uint64(time.Second))
		_, err := mod.ExportedFunction(`testOverride`).Call(ctx)
		if err != nil {
			panic(err.Error())
		}
	})
	t.Run(`del`, func(t *testing.T) {
		hostModule.Del(`TEST_BOOL`)
		hostModule.Del(`TEST_UINT64`)
		hostModule.Del(`TEST_DURATION`)
		_, err := mod.ExportedFunction(`testGet`).Call(ctx)
		if err != nil {
			panic(err.Error())
		}
	})

	hostModule.Stop()
}
