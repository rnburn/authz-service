load("@rules_python//python:defs.bzl", "py_binary")
load("@integration_python_deps//:requirements.bzl", "requirement")

package(default_visibility = ["//visibility:public"])

py_binary(
    name = "server",
    srcs = [
        "server.py",
    ],
    deps = [
        requirement("Flask"),
    ],
)
