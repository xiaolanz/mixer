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
        "@com_github_servicecontrol_client//servicecontrol:go_default_library",
        "@org_golang_google_x_net_context//:go_default_library",
        "@org_golang_google_x_oauth2//:go_default_library",
    ],
)
