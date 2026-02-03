package wazero_global

import (
	"context"
	"log"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// Name is the name of this host module.
const Name = "pantopic/wazero-global"

var (
	DefaultCtxKeyMeta    = `__wazero_global_meta`
	DefaultCtxKeyGlobals = `__wazero_global_map`
)

type meta struct {
	ptrVal     uint32
	ptrName    uint32
	ptrNameLen uint32
	ptrNameCap uint32
}

type hostModule struct {
	sync.RWMutex

	module        api.Module
	ctxKeyMeta    string
	ctxKeyGlobals string
}

type Option func(*hostModule)

func New(opts ...Option) *hostModule {
	p := &hostModule{
		ctxKeyMeta:    DefaultCtxKeyMeta,
		ctxKeyGlobals: DefaultCtxKeyGlobals,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (h *hostModule) Name() string {
	return Name
}
func (h *hostModule) Stop() {}

// Register instantiates the host module, making it available to all module instances in this runtime
func (h *hostModule) Register(ctx context.Context, r wazero.Runtime) (err error) {
	builder := r.NewHostModuleBuilder(Name)
	register := func(name string, fn func(ctx context.Context, m api.Module, stack []uint64)) {
		builder = builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(fn), nil, nil).Export(name)
	}
	for name, fn := range map[string]any{
		"__global_set": func(ctx context.Context, globals map[string]uint64, name string, val uint64) {
			globals[name] = val
		},
		"__global_get": func(ctx context.Context, globals map[string]uint64, name string) uint64 {
			return globals[name]
		},
	} {
		switch fn := fn.(type) {
		case func(ctx context.Context, globals map[string]uint64, name string, val uint64):
			register(name, func(ctx context.Context, m api.Module, stack []uint64) {
				meta := get[*meta](ctx, h.ctxKeyMeta)
				fn(ctx, h.globals(ctx), string(getName(m, meta)), getVal(m, meta))
			})
		default:
			log.Panicf("Method signature implementation missing: %#v", fn)
		}
	}
	h.module, err = builder.Instantiate(ctx)
	return
}

// InitContext retrieves the meta page from the wasm module
func (p *hostModule) InitContext(ctx context.Context, m api.Module) (context.Context, error) {
	stack, err := m.ExportedFunction(`__global`).Call(ctx)
	if err != nil {
		return ctx, err
	}
	meta := &meta{}
	ptr := uint32(stack[0])
	for i, v := range []*uint32{
		&meta.ptrVal,
		&meta.ptrNameCap,
		&meta.ptrNameLen,
		&meta.ptrName,
	} {
		*v = readUint32(m, ptr+uint32(4*i))
	}
	return context.WithValue(ctx, p.ctxKeyMeta, meta), nil
}

// ContextCopy populates dst context with the meta page from src context.
func (h *hostModule) ContextCopy(dst, src context.Context) context.Context {
	dst = context.WithValue(dst, h.ctxKeyMeta, get[*meta](src, h.ctxKeyMeta))
	dst = context.WithValue(dst, h.ctxKeyGlobals, make(map[string]uint64))
	return dst
}

func (p *hostModule) globals(ctx context.Context) map[string]uint64 {
	return get[map[string]uint64](ctx, p.ctxKeyGlobals)
}

func get[T any](ctx context.Context, key string) T {
	v := ctx.Value(key)
	if v == nil {
		log.Panicf("Context item missing %s", key)
	}
	return v.(T)
}

func readUint32(m api.Module, ptr uint32) (val uint32) {
	val, ok := m.Memory().ReadUint32Le(ptr)
	if !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
	return
}

func getVal(m api.Module, meta *meta) uint64 {
	return readUint64(m, meta.ptrVal)
}

func setVal(m api.Module, meta *meta, i uint64) {
	writeUint64(m, meta.ptrVal, i)
}

func getName(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrName, meta.ptrNameLen, meta.ptrNameCap)
}

func setName(m api.Module, meta *meta, i uint64) {
	writeUint64(m, meta.ptrName, i)
}

func setData(m api.Module, meta *meta, b []byte) {
	copy(read(m, meta.ptrName, uint32(len(b)), meta.ptrNameCap), b)
	writeUint32(m, meta.ptrNameLen, uint32(len(b)))
}

func read(m api.Module, ptrData, ptrLen, ptrMax uint32) (buf []byte) {
	buf, ok := m.Memory().Read(ptrData, readUint32(m, ptrMax))
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", ptrData, ptrLen)
	}
	return buf[:readUint32(m, ptrLen)]
}

func readUint64(m api.Module, ptr uint32) (val uint64) {
	val, ok := m.Memory().ReadUint64Le(ptr)
	if !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
	return
}

func writeUint32(m api.Module, ptr uint32, val uint32) {
	if ok := m.Memory().WriteUint32Le(ptr, val); !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
}

func writeUint64(m api.Module, ptr uint32, val uint64) {
	if ok := m.Memory().WriteUint64Le(ptr, val); !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
}
