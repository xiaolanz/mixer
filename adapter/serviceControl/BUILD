package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
    ],
    visibility = ["//cmd:__subpackages__"],
    deps = [
        "@com_github_googleapis_googleapis//:google/api/servicecontrol/v1",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

gogoslick_proto_library(
    name = "go_default_library",
    importmap = {
        "google/protobuf/duration.proto": "github.com/gogo/protobuf/types",
        "gogoproto/gogo.proto": "github.com/gogo/protobuf/gogoproto",
    },
    imports = [
        "external/com_github_gogo_protobuf",
        "external/com_github_google_protobuf/src",
    ],
    inputs = [
        "@com_github_google_protobuf//:well_known_protos",
        "@com_github_gogo_protobuf//gogoproto:go_default_library_protos",
    ],
    protos = [
        "config.proto",
    ],
    verbose = 0,
    visibility = ["//adapter/statsd:__pkg__"],
    deps = [
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//sortkeys:go_default_library",
        "@com_github_gogo_protobuf//types:go_default_library",
    ],
)
