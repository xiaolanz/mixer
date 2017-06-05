package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "reportMetrics.go",
    ],
    visibility = ["//cmd:__subpackages__"],
    deps = [
        "//adapter/serviceControl/config:go_default_library",
        "@com_github_googleapis_googleapis//:google/api/servicecontrol/v1",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
