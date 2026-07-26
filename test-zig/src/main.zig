const std = @import("std");
const global = @import("global");

var test_bool: global.Bool = undefined;
var test_uint64: global.Uint64 = undefined;
var test_duration: global.Duration = undefined;

// Like tinygo's `-buildmode=wasi-legacy`, run initialization on _start without
// calling proc_exit so the instance remains callable after instantiation.
export fn _start() void {
    test_bool = global.newBool("TEST_BOOL", true);
    test_uint64 = global.newUint64("TEST_UINT64", 42);
    test_duration = global.newDuration("TEST_DURATION", std.time.ns_per_min);
}

export fn testGet() void {
    if (!test_bool.get())
        @panic("Get bool failed");
    if (test_uint64.get() != 42)
        @panic("Get uint64 failed");
    if (test_duration.get() != std.time.ns_per_min)
        @panic("Get duration failed");
}

export fn testOverride() void {
    if (test_bool.get())
        @panic("Override bool failed");
    if (test_uint64.get() != 43)
        @panic("Override uint64 failed");
    if (test_duration.get() != std.time.ns_per_s)
        @panic("Override duration failed");
}
