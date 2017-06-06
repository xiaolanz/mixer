load("@org_pubref_rules_protobuf//gogo:rules.bzl", "gogoslick_proto_library")

gogoslick_proto_library(
    name = "go_default_library",
    importmap = {
        "google/rpc/status.proto": "github.com/googleapis/googleapis/google/rpc",
        "google/protobuf/any.proto": "github.com/gogo/protobuf/types",
        "google/protobuf/struct.proto": "github.com/gogo/protobuf/types",
        "google/protobuf/timestamp.proto": "github.com/gogo/protobuf/types",
        "google/protobuf/duration.proto": "github.com/gogo/protobuf/types",
    },
    imports = [
        "external/com_github_google_protobuf/src",
    ],
    inputs = [
        "@com_github_google_protobuf//:well_known_protos",
    ],
    protos = [
        "service_controller.proto",
    ],
    verbose = 0,
    visibility = ["//adapter/serviceControl:__pkg__"],
    with_grpc = True,
    deps = [
        "@com_github_gogo_protobuf//types:go_default_library",
    ],
)
