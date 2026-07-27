const std = @import("std");

pub fn build(b: *std.Build) void {
    const optimize = b.standardOptimizeOption(.{ .preferred_optimize_mode = .ReleaseSmall });
    const target = b.resolveTargetQuery(.{ .cpu_arch = .wasm32, .os_tag = .wasi });
    const exe = b.addExecutable(.{
        .name = "test-zig",
        .root_module = b.createModule(.{
            .root_source_file = b.path("src/main.zig"),
            .target = target,
            .optimize = optimize,
        }),
    });
    exe.entry = .disabled;
    exe.rdynamic = true;
    exe.root_module.addImport("global", b.dependency("sdk", .{}).module("global"));
    b.installArtifact(exe);
}
