//! Guest SDK (Zig) for the pantopic/wazero-global host module.
//!
//! Implements the same ABI as sdk-go: the host calls `__global` once to learn
//! the addresses of the value, name buffer and name length/capacity, then each
//! call to `__global_get` lets the host overwrite the value in place.

const std = @import("std");

var val: u64 = 0;
var name_cap: u32 = 64;
var name_len: u32 = 0;
var name: [64]u8 = undefined;
var meta: [4]u32 = undefined;

extern "pantopic/wazero-global" fn __global_get() void;

export fn __global() u32 {
    meta[0] = @intFromPtr(&val);
    meta[1] = @intFromPtr(&name_cap);
    meta[2] = @intFromPtr(&name_len);
    meta[3] = @intFromPtr(&name);
    return @intFromPtr(&meta);
}

fn fetch(n: []const u8, v: u64) u64 {
    val = v;
    name_len = @intCast(n.len);
    @memcpy(name[0..n.len], n);
    __global_get();
    return val;
}

pub const Bool = struct {
    name: []const u8,
    default: bool,

    pub fn get(self: Bool) bool {
        return fetch(self.name, @intFromBool(self.default)) == 1;
    }
};

pub const Uint64 = struct {
    name: []const u8,
    default: u64,

    pub fn get(self: Uint64) u64 {
        return fetch(self.name, self.default);
    }
};

/// Duration in nanoseconds, matching Go's time.Duration representation.
pub const Duration = struct {
    name: []const u8,
    default: i64,

    pub fn get(self: Duration) i64 {
        return @bitCast(fetch(self.name, @bitCast(self.default)));
    }
};

pub fn newBool(n: []const u8, v: bool) Bool {
    return .{ .name = n, .default = v };
}

pub fn newUint64(n: []const u8, v: u64) Uint64 {
    return .{ .name = n, .default = v };
}

pub fn newDuration(n: []const u8, v: i64) Duration {
    return .{ .name = n, .default = v };
}
