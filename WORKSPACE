workspace(
    name = "com_github_mjm_pi_tools",
    #    managed_directories = {
    #        "@npm_homebase": ["homebase/node_modules"],
    #    },
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "7c10271940c6bce577d51a075ae77728964db285dac0a46614a7934dc34303e6",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.26.0/rules_go-v0.26.0.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.26.0/rules_go-v0.26.0.tar.gz",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "62ca106be173579c0a167deb23358fdfe71ffa1e4cfdddf5582af26520f1c66f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

go_rules_dependencies()

go_register_toolchains(version = "1.16")

gazelle_dependencies()

# Download the rules_docker repository at release v0.14.4
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "4521794f0fba2e20f3bf15846ab5e01d5332e587e9ce81629c7f96c793bb7036",
    strip_prefix = "rules_docker-0.14.4",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.14.4/rules_docker-v0.14.4.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

load("@io_bazel_rules_docker//repositories:pip_repositories.bzl", io_bazel_rules_docker_pip_deps = "pip_deps")

io_bazel_rules_docker_pip_deps()

load("@io_bazel_rules_docker//container:container.bzl", "container_pull")

container_pull(
    name = "distroless_base",
    architecture = "arm64",
    digest = "sha256:a21a1f3eb18297325d39ff95534a694daad8cab80aa4230f02315ada36754eee",
    os = "linux",
    registry = "gcr.io",
    repository = "distroless/base",
)

container_pull(
    name = "ubuntu_bluetooth",
    architecture = "arm64",
    digest = "sha256:c95fcbc19ac85f52c228a4b0d9d2bca05d6c251f9659b314acaa2bda9eee6915",
    os = "linux",
    registry = "index.docker.io",
    repository = "mmoriarity/ubuntu-bluetooth",
)

container_pull(
    name = "rsync",
    architecture = "arm64",
    digest = "sha256:0bde8d3a682256c105ad2345bfa0e55e9c457225f4866b3cfaf6f4259c0c89d9",
    os = "linux",
    registry = "index.docker.io",
    repository = "mmoriarity/rsync",
)

container_pull(
    name = "backup",
    architecture = "arm64",
    digest = "sha256:df8e226d7174435515e0f024fba184e8477f804757a1ddf4ea2d1e60c08b6a92",
    os = "linux",
    registry = "index.docker.io",
    repository = "mmoriarity/backup",
)

container_pull(
    name = "node",
    architecture = "arm64",
    digest = "sha256:3647b6e57d6f4d8594526941b052d5823f131523b0fb0723b98bb1d423be80e6",  # 14.16.0-buster-slim
    os = "linux",
    registry = "index.docker.io",
    repository = "node",
)

http_archive(
    name = "rules_pkg",
    sha256 = "352c090cc3d3f9a6b4e676cf42a6047c16824959b438895a76c2989c6d7c246a",
    url = "https://github.com/bazelbuild/rules_pkg/releases/download/0.2.5/rules_pkg-0.2.5.tar.gz",
)

load("@rules_pkg//:deps.bzl", "rules_pkg_dependencies")

rules_pkg_dependencies()

#http_archive(
#    name = "build_bazel_rules_nodejs",
#    sha256 = "dd7ea7efda7655c218ca707f55c3e1b9c68055a70c31a98f264b3445bc8f4cb1",
#    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/3.2.3/rules_nodejs-3.2.3.tar.gz"],
#)
#
#load("@build_bazel_rules_nodejs//:index.bzl", "yarn_install")
#
#yarn_install(
#    name = "npm_homebase",
#    package_json = "//homebase:package.json",
#    yarn_lock = "//homebase:yarn.lock",
#)

http_archive(
    name = "rules_proto",
    sha256 = "602e7161d9195e50246177e7c55b2f39950a9cf7366f74ed5f22fd45750cd208",
    strip_prefix = "rules_proto-97d8af4dc474595af3900dd85cb3a29ad28cc313",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_proto/archive/97d8af4dc474595af3900dd85cb3a29ad28cc313.tar.gz",
        "https://github.com/bazelbuild/rules_proto/archive/97d8af4dc474595af3900dd85cb3a29ad28cc313.tar.gz",
    ],
)

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")

rules_proto_dependencies()

rules_proto_toolchains()

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
    name = "com_google_protobuf",
    commit = "09745575a923640154bcf307fba8aedff47f240a",
    remote = "https://github.com/protocolbuffers/protobuf",
    shallow_since = "1558721209 -0700",
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "build_bazel_rules_apple",
    sha256 = "734813e44eb5a2fcba5ffd45de9fe5d05325420a5aa1f6c97a3d88fe2c525b17",
    url = "https://github.com/bazelbuild/rules_apple/releases/download/0.21.1/rules_apple.0.21.1.tar.gz",
)

load(
    "@build_bazel_rules_apple//apple:repositories.bzl",
    "apple_rules_dependencies",
)

apple_rules_dependencies()

load(
    "@build_bazel_rules_swift//swift:repositories.bzl",
    "swift_rules_dependencies",
)

swift_rules_dependencies()

load(
    "@build_bazel_apple_support//lib:repositories.bzl",
    "apple_support_dependencies",
)

apple_support_dependencies()

http_archive(
    name = "com_github_relay_tools_relay_swift",
    sha256 = "1f0b68b5687267e592372b65b106d0d5ed03f6a21a02620d4841206ba9c3ace0",
    strip_prefix = "Relay.swift-e965401ec4734652ce71ddd19ee39a94e32c2e92",
    urls = ["https://github.com/relay-tools/Relay.swift/archive/e965401ec4734652ce71ddd19ee39a94e32c2e92.zip"],
)

#local_repository(
#    name = "com_github_relay_tools_relay_swift",
#    path = "../Relay.swift",
#)

http_file(
    name = "otelcol",
    sha256 = "1fed8f4c3cbacd2a3b4bc8ca03f64956054dd88f2d2551717dba6c3e8a005dc4",
    urls = ["https://github.com/open-telemetry/opentelemetry-collector/releases/download/v0.24.0/otelcol_linux_arm64"],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

go_repository(
    name = "com_github_cyberdelia_templates",
    importpath = "github.com/cyberdelia/templates",
    sum = "h1:/ovYnF02fwL0kvspmy9AuyKg1JhdTRUgPw4nUxd9oZM=",
    version = "v0.0.0-20141128023046-ca7fffd4298c",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_deepmap_oapi_codegen",
    importpath = "github.com/deepmap/oapi-codegen",
    sum = "h1:Wj44p9A0V0PJ+AUg0BWdyGcsS1LY18U+0rCuPQgK0+o=",
    version = "v1.3.6",
)

go_repository(
    name = "com_github_dgrijalva_jwt_go",
    importpath = "github.com/dgrijalva/jwt-go",
    sum = "h1:7qlOGliEKZXTDg6OTjfoBKDXWrumCAMpl/TFQ4/5kLM=",
    version = "v3.2.0+incompatible",
)

go_repository(
    name = "com_github_getkin_kin_openapi",
    importpath = "github.com/getkin/kin-openapi",
    sum = "h1:PbHHtYZpjKwZtGlIyELgA2DploRrsaXztoNNx9HjwNY=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_ghodss_yaml",
    importpath = "github.com/ghodss/yaml",
    sum = "h1:wQHKEahhL6wmXdzwWG11gIVCkOv05bNOh+Rxn0yngAk=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_go_chi_chi",
    importpath = "github.com/go-chi/chi",
    sum = "h1:maB6vn6FqCxrpz4FqWdh4+lwpyZIQS7YEAUcHlgXVRs=",
    version = "v4.0.2+incompatible",
)

go_repository(
    name = "com_github_golangci_lint_1",
    importpath = "github.com/golangci/lint-1",
    sum = "h1:utua3L2IbQJmauC5IXdEA547bcoU5dozgQAfc8Onsg4=",
    version = "v0.0.0-20181222135242-d2cdd8c08219",
)

go_repository(
    name = "com_github_influxdata_influxdb_client_go_v2",
    importpath = "github.com/influxdata/influxdb-client-go/v2",
    sum = "h1:vRla3taM+zkziP1NUGfN6Y6zJ9ZSSMg0fs/JhCGyX1s=",
    version = "v2.0.1",
)

go_repository(
    name = "com_github_influxdata_line_protocol",
    importpath = "github.com/influxdata/line-protocol",
    sum = "h1:W9WBk7wlPfJLvMCdtV4zPulc4uCPrlywQOmbFOhgQNU=",
    version = "v0.0.0-20200327222509-2487e7298839",
)

go_repository(
    name = "com_github_kr_pretty",
    importpath = "github.com/kr/pretty",
    sum = "h1:Fmg33tUaq4/8ym9TJN1x7sLJnHVwhP33CNkpYV/7rwI=",
    version = "v0.2.1",
)

go_repository(
    name = "com_github_kr_pty",
    importpath = "github.com/kr/pty",
    sum = "h1:AkaSdXYQOWeaO3neb8EM634ahkXXe3jYbVh/F9lq+GI=",
    version = "v1.1.8",
)

go_repository(
    name = "com_github_kr_text",
    importpath = "github.com/kr/text",
    sum = "h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_labstack_echo_v4",
    importpath = "github.com/labstack/echo/v4",
    sum = "h1:z0BZoArY4FqdpUEl+wlHp4hnr/oSR6MTmQmv8OHSoww=",
    version = "v4.1.11",
)

go_repository(
    name = "com_github_labstack_gommon",
    importpath = "github.com/labstack/gommon",
    sum = "h1:JEeO0bvc78PKdyHxloTKiF8BD5iGrH8T6MSeGvSgob0=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_matryer_moq",
    importpath = "github.com/matryer/moq",
    sum = "h1:HvFwW+cm9bCbZ/+vuGNq7CRWXql8c0y8nGeYpqmpvmk=",
    version = "v0.0.0-20190312154309-6cfb0558e1bd",
)

go_repository(
    name = "com_github_mattn_go_colorable",
    importpath = "github.com/mattn/go-colorable",
    sum = "h1:6Su7aK7lXmJ/U79bYtBjLNaha4Fs1Rg9plHpcH+vvnE=",
    version = "v0.1.6",
)

go_repository(
    name = "com_github_mattn_go_isatty",
    importpath = "github.com/mattn/go-isatty",
    sum = "h1:wuysRhFDzyxgEmMf5xjvJ2M9dZoWAXNNr5LSBS7uHXY=",
    version = "v0.0.12",
)

go_repository(
    name = "com_github_pkg_errors",
    importpath = "github.com/pkg/errors",
    sum = "h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=",
    version = "v0.9.1",
)

go_repository(
    name = "com_github_pmezard_go_difflib",
    importpath = "github.com/pmezard/go-difflib",
    sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_stretchr_objx",
    importpath = "github.com/stretchr/objx",
    sum = "h1:Hbg2NidpLE8veEBkEZTL3CvlkUIVzuU9jDplZO54c48=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    sum = "h1:nwc3DEeHmmLAfoZucVR881uASk0Mfjw8xYJ99tb5CcY=",
    version = "v1.7.0",
)

go_repository(
    name = "com_github_valyala_bytebufferpool",
    importpath = "github.com/valyala/bytebufferpool",
    sum = "h1:GqA5TC/0021Y/b9FG4Oi9Mr3q7XYx6KllzawFIhcdPw=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_valyala_fasttemplate",
    importpath = "github.com/valyala/fasttemplate",
    sum = "h1:RZqt0yGBsps8NGvLSGW804QQqCUYYLsaOjTVHy1Ocw4=",
    version = "v1.1.0",
)

go_repository(
    name = "in_gopkg_check_v1",
    importpath = "gopkg.in/check.v1",
    sum = "h1:BLraFXnmrev5lT+xlilqcH8XK9/i0At2xKjWk4p6zsU=",
    version = "v1.0.0-20200227125254-8fa46927fb4f",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    sum = "h1:clyUAQHOM3G0M3f5vQj7LuJrETvjVot3Z5el9nffUtU=",
    version = "v2.3.0",
)

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    sum = "h1:duLSQW+DZ5MsXKX7kc4rXlq6/mmxz4G6ewJuBPlhRe0=",
    version = "v0.0.0-20200930160638-afb6bcd081ae",
)

go_repository(
    name = "org_golang_x_net",
    importpath = "golang.org/x/net",
    sum = "h1:003p0dJM77cxMSyCPFphvZf/Y5/NXf5fzg6ufd1/Oew=",
    version = "v0.0.0-20210119194325-5f4716e94777",
)

go_repository(
    name = "org_golang_x_sync",
    importpath = "golang.org/x/sync",
    sum = "h1:5KslGYwFpkhGh+Q16bwMP3cOontH8FOep7tGV86Y7SQ=",
    version = "v0.0.0-20210220032951-036812b2e83c",
)

go_repository(
    name = "org_golang_x_sys",
    importpath = "golang.org/x/sys",
    sum = "h1:ggRgirZABFolTmi3sn6Ivd9SipZwLedQ5wR0aAKnFxU=",
    version = "v0.0.0-20210305230114-8fe3ee5dd75b",
)

go_repository(
    name = "org_golang_x_text",
    importpath = "golang.org/x/text",
    sum = "h1:i6eZZ+zk0SOf0xgBpEpPD18qWcJda6q1sxt3S0kzyUQ=",
    version = "v0.3.5",
)

go_repository(
    name = "org_golang_x_tools",
    importpath = "golang.org/x/tools",
    sum = "h1:po9/4sTYwZU9lPhi1tOrb4hCv3qrhiQ77LZfGa2OjwY=",
    version = "v0.1.0",
)

go_repository(
    name = "org_golang_x_xerrors",
    importpath = "golang.org/x/xerrors",
    sum = "h1:go1bK/D/BFZV2I8cIQd1NKEZ+0owSTG1fDTci4IqFcE=",
    version = "v0.0.0-20200804184101-5ec99f83aff1",
)

go_repository(
    name = "com_github_alecthomas_template",
    importpath = "github.com/alecthomas/template",
    sum = "h1:JYp7IbQjafoB+tBA3gMyHYHrpOtNuDiK/uB5uXxq5wM=",
    version = "v0.0.0-20190718012654-fb15b899a751",
)

go_repository(
    name = "com_github_alecthomas_units",
    importpath = "github.com/alecthomas/units",
    sum = "h1:UQZhZ2O0vMHr2cI+DC1Mbh0TJxzA3RcLoMsFw+aXw7E=",
    version = "v0.0.0-20190924025748-f65c72e2690d",
)

go_repository(
    name = "com_github_beorn7_perks",
    importpath = "github.com/beorn7/perks",
    sum = "h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_cespare_xxhash_v2",
    importpath = "github.com/cespare/xxhash/v2",
    sum = "h1:6MnRN8NT7+YBpUIWxHtefFZOKTAPgGjpQSxqLNn0+qY=",
    version = "v2.1.1",
)

go_repository(
    name = "com_github_go_kit_kit",
    importpath = "github.com/go-kit/kit",
    sum = "h1:dXFJfIHVvUcpSgDOV+Ne6t7jXri8Tfv2uOLHUZ2XNuo=",
    version = "v0.10.0",
)

go_repository(
    name = "com_github_go_logfmt_logfmt",
    importpath = "github.com/go-logfmt/logfmt",
    sum = "h1:TrB8swr/68K7m9CcGut2g3UOihhbcbiMAYiuTXdEih4=",
    version = "v0.5.0",
)

go_repository(
    name = "com_github_go_stack_stack",
    importpath = "github.com/go-stack/stack",
    sum = "h1:5SgMzNM5HxrEjV0ww2lTmX6E2Izsfxas4+YHWRs3Lsk=",
    version = "v1.8.0",
)

go_repository(
    name = "com_github_gogo_protobuf",
    importpath = "github.com/gogo/protobuf",
    sum = "h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=",
    version = "v1.3.2",
)

go_repository(
    name = "com_github_golang_protobuf",
    importpath = "github.com/golang/protobuf",
    sum = "h1:JjCZWpVbqXDqFVmTfYWEVTMIYrL/NPdPSCHPJ0T/raM=",
    version = "v1.4.3",
)

go_repository(
    name = "com_github_google_go_cmp",
    importpath = "github.com/google/go-cmp",
    sum = "h1:Khx7svrCpmxxtHBq5j2mp/xVjsi8hQMfNLvJFAlrGgU=",
    version = "v0.5.5",
)

go_repository(
    name = "com_github_google_gofuzz",
    importpath = "github.com/google/gofuzz",
    sum = "h1:Hsa8mG0dQ46ij8Sl2AYJDUv1oA9/d6Vk+3LG99Oe02g=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_json_iterator_go",
    importpath = "github.com/json-iterator/go",
    sum = "h1:Kz6Cvnvv2wGdaG/V8yMvfkmNiXq9Ya2KUv4rouJJr68=",
    version = "v1.1.10",
)

go_repository(
    name = "com_github_julienschmidt_httprouter",
    importpath = "github.com/julienschmidt/httprouter",
    sum = "h1:U0609e9tgbseu3rBINet9P48AI/D3oJs4dN7jwJOQ1U=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_konsorten_go_windows_terminal_sequences",
    importpath = "github.com/konsorten/go-windows-terminal-sequences",
    sum = "h1:CE8S1cTafDpPvMhIxNJKvHsGVBgn1xWYf1NbHQhywc8=",
    version = "v1.0.3",
)

go_repository(
    name = "com_github_kr_logfmt",
    importpath = "github.com/kr/logfmt",
    sum = "h1:T+h1c/A9Gawja4Y9mFVWj2vyii2bbUNDw3kt9VxK2EY=",
    version = "v0.0.0-20140226030751-b84e30acd515",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
    sum = "h1:I0XW9+e1XWDxdcEniV4rQAIOPUGDq67JSCiRCgGCZLI=",
    version = "v1.0.2-0.20181231171920-c182affec369",
)

go_repository(
    name = "com_github_modern_go_concurrent",
    importpath = "github.com/modern-go/concurrent",
    sum = "h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=",
    version = "v0.0.0-20180306012644-bacd9c7ef1dd",
)

go_repository(
    name = "com_github_modern_go_reflect2",
    importpath = "github.com/modern-go/reflect2",
    sum = "h1:9f412s+6RmYXLWZSEzVVgPGK7C2PphHj5RJrvfx9AWI=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_mwitkow_go_conntrack",
    importpath = "github.com/mwitkow/go-conntrack",
    sum = "h1:KUppIJq7/+SVif2QVs3tOP0zanoHgBEVAwHxUSIzRqU=",
    version = "v0.0.0-20190716064945-2f068394615f",
)

go_repository(
    name = "com_github_prometheus_client_golang",
    importpath = "github.com/prometheus/client_golang",
    sum = "h1:Rrch9mh17XcxvEu9D9DEpb4isxjGBtcevQjKvxPRQIU=",
    version = "v1.9.0",
)

go_repository(
    name = "com_github_prometheus_client_model",
    importpath = "github.com/prometheus/client_model",
    sum = "h1:uq5h0d+GuxiXLJLNABMgp2qUWDPiLvgCzz2dUR+/W/M=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_prometheus_common",
    importpath = "github.com/prometheus/common",
    sum = "h1:4fgOnadei3EZvgRwxJ7RMpG1k1pOZth5Pc13tyspaKM=",
    version = "v0.15.0",
)

go_repository(
    name = "com_github_prometheus_procfs",
    importpath = "github.com/prometheus/procfs",
    sum = "h1:wH4vA7pcjKuZzjF7lM8awk4fnuJO6idemZXoKnULUx4=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    importpath = "github.com/sirupsen/logrus",
    sum = "h1:UBcNElsrwanuuMsnGSlYmtmgbb23qDR5dG+6X6Oo89I=",
    version = "v1.6.0",
)

go_repository(
    name = "in_gopkg_alecthomas_kingpin_v2",
    importpath = "gopkg.in/alecthomas/kingpin.v2",
    sum = "h1:jMFz6MfLP0/4fUyZle81rXUoxOBFi19VUFKVDOQfozc=",
    version = "v2.2.6",
)

go_repository(
    name = "org_golang_google_protobuf",
    importpath = "google.golang.org/protobuf",
    sum = "h1:Ejskq+SyPohKW+1uil0JJMtmHCgJPJ/qWTxr8qp+R4c=",
    version = "v1.25.0",
)

go_repository(
    name = "co_honnef_go_tools",
    importpath = "honnef.co/go/tools",
    sum = "h1:UoveltGrhghAA7ePc+e+QYDHXrBps2PqFZiHkGR/xK8=",
    version = "v0.0.1-2020.1.4",
)

go_repository(
    name = "com_github_apache_arrow_go_arrow",
    importpath = "github.com/apache/arrow/go/arrow",
    sum = "h1:5ultmol0yeX75oh1hY78uAFn3dupBQ/QUNxERCkiaUQ=",
    version = "v0.0.0-20200601151325-b2287a20f230",
)

go_repository(
    name = "com_github_aws_aws_sdk_go",
    importpath = "github.com/aws/aws-sdk-go",
    sum = "h1:r0puXncSaAfRt7Btml2swUo74Kao+vKhO3VLjwDjK54=",
    version = "v1.35.3",
)

go_repository(
    name = "com_github_azure_go_ansiterm",
    importpath = "github.com/Azure/go-ansiterm",
    sum = "h1:w+iIsaOQNcT7OZ575w+acHgRric5iCyQh+xv+KJ4HB8=",
    version = "v0.0.0-20170929234023-d6e3b3328b78",
)

go_repository(
    name = "com_github_bitly_go_hostpool",
    importpath = "github.com/bitly/go-hostpool",
    sum = "h1:mXoPYz/Ul5HYEDvkta6I8/rnYM5gSdSV2tJ6XbZuEtY=",
    version = "v0.0.0-20171023180738-a3a6125de932",
)

go_repository(
    name = "com_github_bkaradzic_go_lz4",
    importpath = "github.com/bkaradzic/go-lz4",
    sum = "h1:RXc4wYsyz985CkXXeX04y4VnZFGG8Rd43pRaHsOXAKk=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_bmizerany_assert",
    importpath = "github.com/bmizerany/assert",
    sum = "h1:DDGfHa7BWjL4YnC6+E63dPcxHo2sUxDIu8g3QgEJdRY=",
    version = "v0.0.0-20160611221934-b7ed37b82869",
)

go_repository(
    name = "com_github_burntsushi_toml",
    importpath = "github.com/BurntSushi/toml",
    sum = "h1:WXkYYl6Yr3qBf1K79EBnL4mak0OimBfB0XUf9Vl28OQ=",
    version = "v0.3.1",
)

go_repository(
    name = "com_github_burntsushi_xgb",
    importpath = "github.com/BurntSushi/xgb",
    sum = "h1:1BDTz0u9nC3//pOCMdNH+CiXJVYJh5UQNCOBG7jbELc=",
    version = "v0.0.0-20160522181843-27f122750802",
)

go_repository(
    name = "com_github_cenkalti_backoff_v4",
    importpath = "github.com/cenkalti/backoff/v4",
    sum = "h1:JIufpQLbh4DkbQoii76ItQIUFzevQSqOLZca4eamEDs=",
    version = "v4.0.2",
)

go_repository(
    name = "com_github_census_instrumentation_opencensus_proto",
    importpath = "github.com/census-instrumentation/opencensus-proto",
    sum = "h1:glEXhBS5PSLLv4IXzLA5yPRVX4bilULVyxxbrfOtDAk=",
    version = "v0.2.1",
)

go_repository(
    name = "com_github_chzyer_logex",
    importpath = "github.com/chzyer/logex",
    sum = "h1:Swpa1K6QvQznwJRcfTfQJmTE72DqScAa40E+fbHEXEE=",
    version = "v1.1.10",
)

go_repository(
    name = "com_github_chzyer_readline",
    importpath = "github.com/chzyer/readline",
    sum = "h1:fY5BOSpyZCqRo5OhCuC+XN+r/bBCmeuuJtjz+bCNIf8=",
    version = "v0.0.0-20180603132655-2972be24d48e",
)

go_repository(
    name = "com_github_chzyer_test",
    importpath = "github.com/chzyer/test",
    sum = "h1:q763qf9huN11kDQavWsoZXJNW3xEE4JJyHa5Q25/sd8=",
    version = "v0.0.0-20180213035817-a1ea475d72b1",
)

go_repository(
    name = "com_github_clickhouse_clickhouse_go",
    importpath = "github.com/ClickHouse/clickhouse-go",
    sum = "h1:HvD2NhKPLSeO3Ots6YV0ePgs4l3wO0bLqa9Uk1yeMOs=",
    version = "v1.3.12",
)

go_repository(
    name = "com_github_client9_misspell",
    importpath = "github.com/client9/misspell",
    sum = "h1:ta993UF76GwbvJcIo3Y68y/M3WxlpEHPWIGDkJYwzJI=",
    version = "v0.3.4",
)

go_repository(
    name = "com_github_cloudflare_golz4",
    importpath = "github.com/cloudflare/golz4",
    sum = "h1:F1EaeKL/ta07PY/k9Os/UFtwERei2/XzGemhpGnBKNg=",
    version = "v0.0.0-20150217214814-ef862a3cdc58",
)

go_repository(
    name = "com_github_cncf_udpa_go",
    importpath = "github.com/cncf/udpa/go",
    sum = "h1:cqQfy1jclcSy/FwLjemeg3SR1yaINm74aQyupQ0Bl8M=",
    version = "v0.0.0-20201120205902-5459f2c99403",
)

go_repository(
    name = "com_github_cockroachdb_apd",
    importpath = "github.com/cockroachdb/apd",
    sum = "h1:3LFP3629v+1aKXU5Q37mxmRxX/pIu1nijXydLShEq5I=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_cockroachdb_cockroach_go",
    importpath = "github.com/cockroachdb/cockroach-go",
    sum = "h1:eApuUG8W2EtBVwxqLlY2wgoqDYOg3WvIHGvW4fUbbow=",
    version = "v0.0.0-20190925194419-606b3d062051",
)

go_repository(
    name = "com_github_containerd_containerd",
    importpath = "github.com/containerd/containerd",
    sum = "h1:LoIzb5y9x5l8VKAlyrbusNPXqBY0+kviRloxFUMFwKc=",
    version = "v1.3.3",
)

go_repository(
    name = "com_github_coreos_go_systemd",
    importpath = "github.com/coreos/go-systemd",
    sum = "h1:JOrtw2xFKzlg+cbHpyrpLDmnN1HqhBfnX7WDiW7eG2c=",
    version = "v0.0.0-20190719114852-fd7a80b32e1f",
)

go_repository(
    name = "com_github_creack_pty",
    importpath = "github.com/creack/pty",
    sum = "h1:uDmaGzcdjhF4i/plgjmEsriH11Y0o7RKapEf/LDaM3w=",
    version = "v1.1.9",
)

go_repository(
    name = "com_github_cznic_mathutil",
    importpath = "github.com/cznic/mathutil",
    sum = "h1:XNT/Zf5l++1Pyg08/HV04ppB0gKxAqtZQBRYiYrUuYk=",
    version = "v0.0.0-20180504122225-ca4c9f2c1369",
)

go_repository(
    name = "com_github_denisenkom_go_mssqldb",
    importpath = "github.com/denisenkom/go-mssqldb",
    sum = "h1:NfhRXXFDPxcF5Cwo06DzeIaE7uuJtAUhsDwH3LNsjos=",
    version = "v0.0.0-20200620013148-b91950f658ec",
)

go_repository(
    name = "com_github_dhui_dktest",
    importpath = "github.com/dhui/dktest",
    sum = "h1:nZSDcnkpbotzT/nEHNsO+JCKY8i1Qoki1AYOpeLRb6M=",
    version = "v0.3.2",
)

go_repository(
    name = "com_github_docker_distribution",
    importpath = "github.com/docker/distribution",
    sum = "h1:a5mlkVzth6W5A4fOsS3D2EO5BUmsJpcB+cRlLU7cSug=",
    version = "v2.7.1+incompatible",
)

go_repository(
    name = "com_github_docker_docker",
    importpath = "github.com/docker/docker",
    sum = "h1:x+LFTD06kE2vF3WSFCG3D4XdqE2FD61524u/cXQLgws=",
    version = "v17.12.0-ce-rc1.0.20200330121334-7f8b4b621b5d+incompatible",
)

go_repository(
    name = "com_github_docker_go_connections",
    importpath = "github.com/docker/go-connections",
    sum = "h1:El9xVISelRB7BuFusrZozjnkIM5YnzCViNKohAFqRJQ=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_docker_go_units",
    importpath = "github.com/docker/go-units",
    sum = "h1:3uh0PgVws3nIA0Q+MwDC8yjEPf9zjRfZZWXZYDct3Tw=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_edsrzf_mmap_go",
    importpath = "github.com/edsrzf/mmap-go",
    sum = "h1:CEBF7HpRnUCSJgGUb5h1Gm7e3VkmVDrR8lvWVLtrOFw=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_envoyproxy_go_control_plane",
    importpath = "github.com/envoyproxy/go-control-plane",
    sum = "h1:EmNYJhPYy0pOFjCx2PrgtaBXmee0iUX9hLlxE1xHOJE=",
    version = "v0.9.9-0.20201210154907-fd9021fe5dad",
)

go_repository(
    name = "com_github_envoyproxy_protoc_gen_validate",
    importpath = "github.com/envoyproxy/protoc-gen-validate",
    sum = "h1:EQciDnbrYxy13PgWoY8AqoxGiPrpgBZ1R8UNe3ddc+A=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_fsnotify_fsnotify",
    importpath = "github.com/fsnotify/fsnotify",
    sum = "h1:hsms1Qyu0jgnwNXIxa+/V/PDsU6CfLf6CNO8H7IWoS4=",
    version = "v1.4.9",
)

go_repository(
    name = "com_github_fsouza_fake_gcs_server",
    importpath = "github.com/fsouza/fake-gcs-server",
    sum = "h1:OeH75kBZcZa3ZE+zz/mFdJ2btt9FgqfjI7gIh9+5fvk=",
    version = "v1.17.0",
)

go_repository(
    name = "com_github_go_gl_glfw",
    importpath = "github.com/go-gl/glfw",
    sum = "h1:QbL/5oDUmRBzO9/Z7Seo6zf912W/a6Sr4Eu0G/3Jho0=",
    version = "v0.0.0-20190409004039-e6da0acd62b1",
)

go_repository(
    name = "com_github_go_gl_glfw_v3_3_glfw",
    importpath = "github.com/go-gl/glfw/v3.3/glfw",
    sum = "h1:WtGNWLvXpe6ZudgnXrq0barxBImvnnJoMEhXAzcbM0I=",
    version = "v0.0.0-20200222043503-6f7a984d4dc4",
)

go_repository(
    name = "com_github_go_sql_driver_mysql",
    importpath = "github.com/go-sql-driver/mysql",
    sum = "h1:ozyZYNQW3x3HtqT1jira07DN2PArx2v7/mN66gGcHOs=",
    version = "v1.5.0",
)

go_repository(
    name = "com_github_gobuffalo_here",
    importpath = "github.com/gobuffalo/here",
    sum = "h1:hYrd0a6gDmWxBM4TnrGw8mQg24iSVoIkHEk7FodQcBI=",
    version = "v0.6.0",
)

go_repository(
    name = "com_github_gocql_gocql",
    importpath = "github.com/gocql/gocql",
    sum = "h1:vF83LI8tAakwEwvWZtrIEx7pOySacl2TOxx6eXk4ePo=",
    version = "v0.0.0-20190301043612-f6df8288f9b4",
)

go_repository(
    name = "com_github_golang_glog",
    importpath = "github.com/golang/glog",
    sum = "h1:VKtxabqXZkF25pY9ekfRL6a582T4P37/31XEstQ5p58=",
    version = "v0.0.0-20160126235308-23def4e6c14b",
)

go_repository(
    name = "com_github_golang_groupcache",
    importpath = "github.com/golang/groupcache",
    sum = "h1:1r7pUrabqp18hOBcwBwiTsbnFeTZHV9eER/QT5JVZxY=",
    version = "v0.0.0-20200121045136-8c9f03a8e57e",
)

go_repository(
    name = "com_github_golang_migrate_migrate_v4",
    importpath = "github.com/golang-migrate/migrate/v4",
    sum = "h1:5S7HMjiq9u50X3+WXpzXPbUj1qUFuZRm8NCsX989Tn4=",
    version = "v4.13.0",
)

go_repository(
    name = "com_github_golang_mock",
    importpath = "github.com/golang/mock",
    sum = "h1:l75CXGRSwbaYNpl/Z2X1XIIAMSCquvXgpVZDhwEIJsc=",
    version = "v1.4.4",
)

go_repository(
    name = "com_github_golang_snappy",
    importpath = "github.com/golang/snappy",
    sum = "h1:Qgr9rKW7uDUkrbSmQeiDsGa8SjGyCOGtuasMWwvp2P4=",
    version = "v0.0.1",
)

go_repository(
    name = "com_github_golang_sql_civil",
    importpath = "github.com/golang-sql/civil",
    sum = "h1:lXe2qZdvpiX5WZkZR4hgp4KJVfY3nMkvmwbVkpv1rVY=",
    version = "v0.0.0-20190719163853-cb61b32ac6fe",
)

go_repository(
    name = "com_github_google_btree",
    importpath = "github.com/google/btree",
    sum = "h1:0udJVsspx3VBr5FwtLhQQtuAsVc79tTq0ocGIPAU6qo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_google_flatbuffers",
    importpath = "github.com/google/flatbuffers",
    sum = "h1:O7CEyB8Cb3/DmtxODGtLHcEvpr81Jm5qLg/hsHnxA2A=",
    version = "v1.11.0",
)

go_repository(
    name = "com_github_google_go_github",
    importpath = "github.com/google/go-github",
    sum = "h1:N0LgJ1j65A7kfXrZnUDaYCs/Sf4rEjNlfyDHW9dolSY=",
    version = "v17.0.0+incompatible",
)

go_repository(
    name = "com_github_google_go_querystring",
    importpath = "github.com/google/go-querystring",
    sum = "h1:Xkwi/a1rcvNg1PPYe5vI8GbeBY/jrVuDX5ASuANWTrk=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_google_martian",
    importpath = "github.com/google/martian",
    sum = "h1:/CP5g8u/VJHijgedC/Legn3BAbAaWPgecwXBIDzw5no=",
    version = "v2.1.0+incompatible",
)

go_repository(
    name = "com_github_google_martian_v3",
    importpath = "github.com/google/martian/v3",
    sum = "h1:wCKgOCHuUEVfsaQLpPSJb7VdYCdTVZQAuOdYm1yc/60=",
    version = "v3.1.0",
)

go_repository(
    name = "com_github_google_pprof",
    importpath = "github.com/google/pprof",
    sum = "h1:41CTEDOoUXp+FxbPYuEhth5dE/s+NT1cRuhSoqhBQ1E=",
    version = "v0.0.0-20210122040257-d980be63207e",
)

go_repository(
    name = "com_github_google_renameio",
    importpath = "github.com/google/renameio",
    sum = "h1:GOZbcHa3HfsPKPlmyPyN2KEohoMXOhdMbHrvbpl2QaA=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_google_uuid",
    importpath = "github.com/google/uuid",
    sum = "h1:EVhdT+1Kseyi1/pUmXKaFxYsDNy9RQYkMWRH68J/W7Y=",
    version = "v1.1.2",
)

go_repository(
    name = "com_github_googleapis_gax_go_v2",
    importpath = "github.com/googleapis/gax-go/v2",
    sum = "h1:sjZBwGj9Jlw33ImPtvFviGYvseOtDM7hkSKB7+Tv3SM=",
    version = "v2.0.5",
)

go_repository(
    name = "com_github_gorilla_context",
    importpath = "github.com/gorilla/context",
    sum = "h1:AWwleXJkX/nhcU9bZSnZoi3h/qGYqQAGhq6zZe/aQW8=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_gorilla_handlers",
    importpath = "github.com/gorilla/handlers",
    sum = "h1:0QniY0USkHQ1RGCLfKxeNHK9bkDHGRYGNDFBCS+YARg=",
    version = "v1.4.2",
)

go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    sum = "h1:VuZ8uybHlWmqV03+zRzdwKL4tUnIp1MAQtp1mIFE1bc=",
    version = "v1.7.4",
)

go_repository(
    name = "com_github_hailocab_go_hostpool",
    importpath = "github.com/hailocab/go-hostpool",
    sum = "h1:5upAirOpQc1Q53c0bnx2ufif5kANL7bfZWcc6VJWJd8=",
    version = "v0.0.0-20160125115350-e80d13ce29ed",
)

go_repository(
    name = "com_github_hashicorp_errwrap",
    importpath = "github.com/hashicorp/errwrap",
    sum = "h1:hLrqtEDnRye3+sgx6z4qVLNuviH3MR5aQ0ykNJa/UYA=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_hashicorp_go_multierror",
    importpath = "github.com/hashicorp/go-multierror",
    sum = "h1:B9UzwGQJehnUY1yNrnwREHc3fGbC2xefo8g4TbElacI=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_hashicorp_golang_lru",
    importpath = "github.com/hashicorp/golang-lru",
    sum = "h1:YDjusn29QI/Das2iO9M0BHnIbxPeyuCHsjMW+lJfyTc=",
    version = "v0.5.4",
)

go_repository(
    name = "com_github_hpcloud_tail",
    importpath = "github.com/hpcloud/tail",
    sum = "h1:8as8OQ+RF1QrsHvWWsKBtBKINhD9QaD1iozA1wrO4aA=",
    version = "v1.0.1-0.20170814160653-37f427138745",
)

go_repository(
    name = "com_github_ianlancetaylor_demangle",
    importpath = "github.com/ianlancetaylor/demangle",
    sum = "h1:mV02weKRL81bEnm8A0HT1/CAelMQDBuQIfLw8n+d6xI=",
    version = "v0.0.0-20200824232613-28f6c0f3b639",
)

go_repository(
    name = "com_github_jackc_chunkreader",
    importpath = "github.com/jackc/chunkreader",
    sum = "h1:4s39bBR8ByfqH+DKm8rQA3E1LHZWB9XWcrz8fqaZbe0=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_jackc_chunkreader_v2",
    importpath = "github.com/jackc/chunkreader/v2",
    sum = "h1:i+RDz65UE+mmpjTfyz0MoVTnzeYxroil2G82ki7MGG8=",
    version = "v2.0.1",
)

go_repository(
    name = "com_github_jackc_pgconn",
    importpath = "github.com/jackc/pgconn",
    sum = "h1:9UIGICxEAW70RQDGilGwsCG63NCcm5amjuBQCFzrmsw=",
    version = "v1.3.2",
)

go_repository(
    name = "com_github_jackc_pgio",
    importpath = "github.com/jackc/pgio",
    sum = "h1:g12B9UwVnzGhueNavwioyEEpAmqMe1E/BN9ES+8ovkE=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_jackc_pgmock",
    importpath = "github.com/jackc/pgmock",
    sum = "h1:JVX6jT/XfzNqIjye4717ITLaNwV9mWbJx0dLCpcRzdA=",
    version = "v0.0.0-20190831213851-13a1b77aafa2",
)

go_repository(
    name = "com_github_jackc_pgpassfile",
    importpath = "github.com/jackc/pgpassfile",
    sum = "h1:/6Hmqy13Ss2zCq62VdNG8tM1wchn8zjSGOBJ6icpsIM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_jackc_pgproto3",
    importpath = "github.com/jackc/pgproto3",
    sum = "h1:FYYE4yRw+AgI8wXIinMlNjBbp/UitDJwfj5LqqewP1A=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_jackc_pgproto3_v2",
    importpath = "github.com/jackc/pgproto3/v2",
    sum = "h1:Rdjp4NFjwHnEslx2b66FfCI2S0LhO4itac3hXz6WX9M=",
    version = "v2.0.1",
)

go_repository(
    name = "com_github_jackc_pgtype",
    importpath = "github.com/jackc/pgtype",
    sum = "h1:xOamcCJ9MFJTxR5bvw3ZXmiP8evQMohdt2VJ57C0W8Q=",
    version = "v0.0.0-20190828014616-a8802b16cc59",
)

go_repository(
    name = "com_github_jackc_pgx_v4",
    importpath = "github.com/jackc/pgx/v4",
    sum = "h1:ZQM8qLT/E/CGD6XX0E6q9FAwxJYmWpJufzmLMaFuzgQ=",
    version = "v4.0.0-pre1.0.20190824185557-6972a5742186",
)

go_repository(
    name = "com_github_jackc_puddle",
    importpath = "github.com/jackc/puddle",
    sum = "h1:KLBBPU++1T3DHtm1B1QaIHy80Vhu0wNMErIFCNgAL8Y=",
    version = "v0.0.0-20190608224051-11cab39313c9",
)

go_repository(
    name = "com_github_jmespath_go_jmespath",
    importpath = "github.com/jmespath/go-jmespath",
    sum = "h1:BEgLn5cpjn8UN1mAw4NjwDrS35OdebyEtFe+9YPoQUg=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_jmoiron_sqlx",
    importpath = "github.com/jmoiron/sqlx",
    sum = "h1:41Ip0zITnmWNR/vHV+S4m+VoUivnWY5E4OJfLZjCJMA=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_jstemmer_go_junit_report",
    importpath = "github.com/jstemmer/go-junit-report",
    sum = "h1:6QPYqodiu3GuPL+7mfx+NwDdp2eTkp9IfEUpgAwUN0o=",
    version = "v0.9.1",
)

go_repository(
    name = "com_github_kardianos_osext",
    importpath = "github.com/kardianos/osext",
    sum = "h1:iQTw/8FWTuc7uiaSepXwyf3o52HaUYcV+Tu66S3F5GA=",
    version = "v0.0.0-20190222173326-2bc1f35cddc0",
)

go_repository(
    name = "com_github_kisielk_errcheck",
    importpath = "github.com/kisielk/errcheck",
    sum = "h1:e8esj/e4R+SAOwFwN+n3zr0nYeCyeweozKfO23MvHzY=",
    version = "v1.5.0",
)

go_repository(
    name = "com_github_kisielk_gotool",
    importpath = "github.com/kisielk/gotool",
    sum = "h1:AV2c/EiW3KqPNT9ZKl07ehoAGi4C5/01Cfbblndcapg=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_lib_pq",
    importpath = "github.com/lib/pq",
    sum = "h1:9xohqzkUwzR4Ga4ivdTcawVS89YSDVxXMa3xJX3cGzg=",
    version = "v1.8.0",
)

go_repository(
    name = "com_github_markbates_pkger",
    importpath = "github.com/markbates/pkger",
    sum = "h1:3MPelV53RnGSW07izx5xGxl4e/sdRD6zqseIk0rMASY=",
    version = "v0.15.1",
)

go_repository(
    name = "com_github_mattn_go_sqlite3",
    importpath = "github.com/mattn/go-sqlite3",
    sum = "h1:jbhqpg7tQe4SupckyijYiy0mJJ/pRyHvXf7JdWK860o=",
    version = "v1.10.0",
)

go_repository(
    name = "com_github_microsoft_go_winio",
    importpath = "github.com/Microsoft/go-winio",
    sum = "h1:93nQ7k53GjoMQ07HVP8g6Zj1fQZDDj7Xy2VkNNtvX8o=",
    version = "v0.4.15-0.20200113171025-3fe6c5262873",
)

go_repository(
    name = "com_github_morikuni_aec",
    importpath = "github.com/morikuni/aec",
    sum = "h1:nP9CBfwrvYnBRgY6qfDQkygYDmYwOilePFkwzv4dU8A=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_mutecomm_go_sqlcipher_v4",
    importpath = "github.com/mutecomm/go-sqlcipher/v4",
    sum = "h1:sV1tWCWGAVlPhNGT95Q+z/txFxuhAYWwHD1afF5bMZg=",
    version = "v4.4.0",
)

go_repository(
    name = "com_github_nakagami_firebirdsql",
    importpath = "github.com/nakagami/firebirdsql",
    sum = "h1:P48LjvUQpTReR3TQRbxSeSBsMXzfK0uol7eRcr7VBYQ=",
    version = "v0.0.0-20190310045651-3c02a58cfed8",
)

go_repository(
    name = "com_github_neo4j_neo4j_go_driver",
    importpath = "github.com/neo4j/neo4j-go-driver",
    sum = "h1:fhFP5RliM2HW/8XdcO5QngSfFli9GcRIpMXvypTQt6E=",
    version = "v1.8.1-0.20200803113522-b626aa943eba",
)

go_repository(
    name = "com_github_onsi_ginkgo",
    importpath = "github.com/onsi/ginkgo",
    sum = "h1:Iw5WCbBcaAAd0fpRb1c9r5YCylv4XDoCSigm1zLevwU=",
    version = "v1.12.0",
)

go_repository(
    name = "com_github_onsi_gomega",
    importpath = "github.com/onsi/gomega",
    sum = "h1:R1uwffexN6Pr340GtYRIdZmAiN4J+iw6WG4wog1DUXg=",
    version = "v1.9.0",
)

go_repository(
    name = "com_github_opencontainers_go_digest",
    importpath = "github.com/opencontainers/go-digest",
    sum = "h1:apOUWs51W5PlhuyGyz9FCeeBIOUDA/6nW8Oi/yOhh5U=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_opencontainers_image_spec",
    importpath = "github.com/opencontainers/image-spec",
    sum = "h1:JMemWkRwHx4Zj+fVxWoMCFm/8sYGGrUVojFA6h/TRcI=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_pierrec_lz4",
    importpath = "github.com/pierrec/lz4",
    sum = "h1:xOYu2+sKj87pJz7V+I7260354UlcRyAZUGhMCToTzVw=",
    version = "v2.2.5+incompatible",
)

go_repository(
    name = "com_github_pkg_browser",
    importpath = "github.com/pkg/browser",
    sum = "h1:49lOXmGaUpV9Fz3gd7TFZY106KVlPVa5jcYD1gaQf98=",
    version = "v0.0.0-20180916011732-0a3d74bf9ce4",
)

go_repository(
    name = "com_github_remyoudompheng_bigfft",
    importpath = "github.com/remyoudompheng/bigfft",
    sum = "h1:HQagqIiBmr8YXawX/le3+O26N+vPPC1PtjaF3mwnook=",
    version = "v0.0.0-20190728182440-6a916e37a237",
)

go_repository(
    name = "com_github_rogpeppe_go_internal",
    importpath = "github.com/rogpeppe/go-internal",
    sum = "h1:RR9dF3JtopPvtkroDZuVD7qquD0bnHlKSqaQhgwt8yk=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_rs_xid",
    importpath = "github.com/rs/xid",
    sum = "h1:mhH9Nq+C1fY2l1XIpgxIiUOfNpRBYH1kKcr+qfKgjRc=",
    version = "v1.2.1",
)

go_repository(
    name = "com_github_rs_zerolog",
    importpath = "github.com/rs/zerolog",
    sum = "h1:uPRuwkWF4J6fGsJ2R0Gn2jB1EQiav9k3S6CSdygQJXY=",
    version = "v1.15.0",
)

go_repository(
    name = "com_github_satori_go_uuid",
    importpath = "github.com/satori/go.uuid",
    sum = "h1:0uYX9dsZ2yD7q2RtLRtPSdGDWzjeM3TbMJP9utgA0ww=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_shopspring_decimal",
    importpath = "github.com/shopspring/decimal",
    sum = "h1:pntxY8Ary0t43dCZ5dqY4YTJCObLY1kIXl0uzMv+7DE=",
    version = "v0.0.0-20180709203117-cd690d0c9e24",
)

go_repository(
    name = "com_github_snowflakedb_glog",
    importpath = "github.com/snowflakedb/glog",
    sum = "h1:CGR1hXCOeoZ1aJhCs8qdKJuEu3xoZnxsLcYoh5Bnr+4=",
    version = "v0.0.0-20180824191149-f5055e6f21ce",
)

go_repository(
    name = "com_github_snowflakedb_gosnowflake",
    importpath = "github.com/snowflakedb/gosnowflake",
    sum = "h1:/Ep0cXv4/3o+iXQvh+6CDjHCRPk2AM42l/AMR9PM94Q=",
    version = "v1.3.5",
)

go_repository(
    name = "com_github_spf13_pflag",
    importpath = "github.com/spf13/pflag",
    sum = "h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=",
    version = "v1.0.5",
)

go_repository(
    name = "com_github_tidwall_pretty",
    importpath = "github.com/tidwall/pretty",
    sum = "h1:BP2bjP495BBPaBcS5rmqviTfrOkN5rO5ceKAMRZCRFc=",
    version = "v0.0.0-20180105212114-65a9db5fad51",
)

go_repository(
    name = "com_github_xanzy_go_gitlab",
    importpath = "github.com/xanzy/go-gitlab",
    sum = "h1:rWtwKTgEnXyNUGrOArN7yyc3THRkpYcKXIXia9abywQ=",
    version = "v0.15.0",
)

go_repository(
    name = "com_github_xdg_scram",
    importpath = "github.com/xdg/scram",
    sum = "h1:u40Z8hqBAAQyv+vATcGgV0YCnDjqSL7/q/JyPhhJSPk=",
    version = "v0.0.0-20180814205039-7eeb5667e42c",
)

go_repository(
    name = "com_github_xdg_stringprep",
    importpath = "github.com/xdg/stringprep",
    sum = "h1:d9X0esnoa3dFsV0FG35rAT0RIhYFlPq7MiP+DW89La0=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_yuin_goldmark",
    importpath = "github.com/yuin/goldmark",
    sum = "h1:ruQGxdhGHe7FWOJPT0mKs5+pD2Xs1Bm/kdGlHO04FmM=",
    version = "v1.2.1",
)

go_repository(
    name = "com_github_zenazn_goji",
    importpath = "github.com/zenazn/goji",
    sum = "h1:RSQQAbXGArQ0dIDEq+PI6WqN6if+5KHu6x2Cx/GXLTQ=",
    version = "v0.9.0",
)

go_repository(
    name = "com_gitlab_nyarla_go_crypt",
    importpath = "gitlab.com/nyarla/go-crypt",
    sum = "h1:7gd+rd8P3bqcn/96gOZa3F5dpJr/vEiDQYlNb/y2uNs=",
    version = "v0.0.0-20160106005555-d9a5dc2b789b",
)

go_repository(
    name = "com_google_cloud_go",
    importpath = "cloud.google.com/go",
    sum = "h1:oKpsiyKMfVpwR3zSAkQixGzlVE5ovitBuO0qSmCf0bI=",
    version = "v0.78.0",
)

go_repository(
    name = "com_google_cloud_go_bigquery",
    importpath = "cloud.google.com/go/bigquery",
    sum = "h1:PQcPefKFdaIzjQFbiyOgAqyx8q5djaE7x9Sqe712DPA=",
    version = "v1.8.0",
)

go_repository(
    name = "com_google_cloud_go_datastore",
    importpath = "cloud.google.com/go/datastore",
    sum = "h1:/May9ojXjRkPBNVrq+oWLqmWCkr4OU5uRY29bu0mRyQ=",
    version = "v1.1.0",
)

go_repository(
    name = "com_google_cloud_go_pubsub",
    importpath = "cloud.google.com/go/pubsub",
    sum = "h1:ukjixP1wl0LpnZ6LWtZJ0mX5tBmjp1f8Sqer8Z2OMUU=",
    version = "v1.3.1",
)

go_repository(
    name = "com_google_cloud_go_spanner",
    importpath = "cloud.google.com/go/spanner",
    sum = "h1:WXuGWhUp5i7MeUMzMrJlodqJvSGtU0Cdw6BdHGgCgVo=",
    version = "v1.9.0",
)

go_repository(
    name = "com_google_cloud_go_storage",
    importpath = "cloud.google.com/go/storage",
    sum = "h1:STgFzyU5/8miMl0//zKh2aQeTyeaUH3WN9bSUiJ09bA=",
    version = "v1.10.0",
)

go_repository(
    name = "com_shuralyov_dmitri_gpu_mtl",
    importpath = "dmitri.shuralyov.com/gpu/mtl",
    sum = "h1:VpgP7xuJadIUuKccphEpTJnWhS2jkQyMt6Y7pJCD7fY=",
    version = "v0.0.0-20190408044501-666a987793e9",
)

go_repository(
    name = "in_gopkg_errgo_v2",
    importpath = "gopkg.in/errgo.v2",
    sum = "h1:0vLT13EuvQ0hNvakwLuFZ/jYrLp5F3kcWHXdRggjCE8=",
    version = "v2.1.0",
)

go_repository(
    name = "in_gopkg_fsnotify_v1",
    importpath = "gopkg.in/fsnotify.v1",
    sum = "h1:xOHLXZwVvI9hhs+cLKq5+I5onOuwQLhQwiu63xxlHs4=",
    version = "v1.4.7",
)

go_repository(
    name = "in_gopkg_inconshreveable_log15_v2",
    importpath = "gopkg.in/inconshreveable/log15.v2",
    sum = "h1:RlWgLqCMMIYYEVcAR5MDsuHlVkaIPDAF+5Dehzg8L5A=",
    version = "v2.0.0-20180818164646-67afb5ed74ec",
)

go_repository(
    name = "in_gopkg_inf_v0",
    importpath = "gopkg.in/inf.v0",
    sum = "h1:73M5CoZyi3ZLMOyDlQh031Cx6N9NDJ2Vvfl76EDAgDc=",
    version = "v0.9.1",
)

go_repository(
    name = "in_gopkg_tomb_v1",
    importpath = "gopkg.in/tomb.v1",
    sum = "h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=",
    version = "v1.0.0-20141024135613-dd632973f1e7",
)

go_repository(
    name = "io_opencensus_go",
    importpath = "go.opencensus.io",
    sum = "h1:dntmOdLpSpHlVqbW5Eay97DelsZHe+55D+xC6i0dDS0=",
    version = "v0.22.5",
)

go_repository(
    name = "io_rsc_binaryregexp",
    importpath = "rsc.io/binaryregexp",
    sum = "h1:HfqmD5MEmC0zvwBuF187nq9mdnXjXsSivRiXN7SmRkE=",
    version = "v0.2.0",
)

go_repository(
    name = "io_rsc_quote_v3",
    importpath = "rsc.io/quote/v3",
    sum = "h1:9JKUTTIUgS6kzR9mK1YuGKv6Nl+DijDNIc0ghT58FaY=",
    version = "v3.1.0",
)

go_repository(
    name = "io_rsc_sampler",
    importpath = "rsc.io/sampler",
    sum = "h1:7uVkIFmeBqHfdjD+gZwtXXI+RODJ2Wc4O7MPEh/QiW4=",
    version = "v1.3.0",
)

go_repository(
    name = "org_golang_google_api",
    importpath = "google.golang.org/api",
    sum = "h1:uWrpz12dpVPn7cojP82mk02XDgTJLDPc2KbVTxrWb4A=",
    version = "v0.40.0",
)

go_repository(
    name = "org_golang_google_appengine",
    importpath = "google.golang.org/appengine",
    sum = "h1:FZR1q0exgwxzPzp/aF+VccGrSfxfPpkBqjIIEq3ru6c=",
    version = "v1.6.7",
)

go_repository(
    name = "org_golang_google_genproto",
    importpath = "google.golang.org/genproto",
    sum = "h1:hcskBH5qZCOa7WpTUFUFvoebnSFZBYpjykLtjIp9DVk=",
    version = "v0.0.0-20210303154014-9728d6b83eeb",
)

go_repository(
    name = "org_golang_google_grpc",
    importpath = "google.golang.org/grpc",
    sum = "h1:o1bcQ6imQMIOpdrO3SWf2z5RV72WbDwdXuK0MDlc8As=",
    version = "v1.36.0",
)

go_repository(
    name = "org_golang_x_exp",
    importpath = "golang.org/x/exp",
    sum = "h1:QE6XYQK6naiK1EPAe1g/ILLxN5RBoH5xkJk3CqlMI/Y=",
    version = "v0.0.0-20200224162631-6cc2880d07d6",
)

go_repository(
    name = "org_golang_x_image",
    importpath = "golang.org/x/image",
    sum = "h1:+qEpEAPhDZ1o0x3tHzZTQDArnOixOzGD9HUJfcg0mb4=",
    version = "v0.0.0-20190802002840-cff245a6509b",
)

go_repository(
    name = "org_golang_x_lint",
    importpath = "golang.org/x/lint",
    sum = "h1:2M3HP5CCK1Si9FQhwnzYhXdG6DXeebvUHFpre8QvbyI=",
    version = "v0.0.0-20201208152925-83fdc39ff7b5",
)

go_repository(
    name = "org_golang_x_mobile",
    importpath = "golang.org/x/mobile",
    sum = "h1:4+4C/Iv2U4fMZBiMCc98MG1In4gJY5YRhtpDNeDeHWs=",
    version = "v0.0.0-20190719004257-d2bd2a29d028",
)

go_repository(
    name = "org_golang_x_mod",
    importpath = "golang.org/x/mod",
    sum = "h1:Kvvh58BN8Y9/lBi7hTekvtMpm07eUZ0ck5pRHpsMWrY=",
    version = "v0.4.1",
)

go_repository(
    name = "org_golang_x_oauth2",
    importpath = "golang.org/x/oauth2",
    sum = "h1:alLDrZkL34Y2bnGHfvC1CYBRBXCXgx8AC2vY4MRtYX4=",
    version = "v0.0.0-20210220000619-9bb904979d93",
)

go_repository(
    name = "org_golang_x_time",
    importpath = "golang.org/x/time",
    sum = "h1:EHBhcS0mlXEAVwNyO2dLfjToGsyY4j24pTs2ScHnX7s=",
    version = "v0.0.0-20200630173020-3af7569d3a1e",
)

go_repository(
    name = "org_modernc_b",
    importpath = "modernc.org/b",
    sum = "h1:vpvqeyp17ddcQWF29Czawql4lDdABCDRbXRAS4+aF2o=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_db",
    importpath = "modernc.org/db",
    sum = "h1:2c6NdCfaLnshSvY7OU09cyAY0gYXUZj4lmg5ItHyucg=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_file",
    importpath = "modernc.org/file",
    sum = "h1:9/PdvjVxd5+LcWUQIfapAWRGOkDLK90rloa8s/au06A=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_fileutil",
    importpath = "modernc.org/fileutil",
    sum = "h1:Z1AFLZwl6BO8A5NldQg/xTSjGLetp+1Ubvl4alfGx8w=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_golex",
    importpath = "modernc.org/golex",
    sum = "h1:wWpDlbK8ejRfSyi0frMyhilD3JBvtcx2AdGDnU+JtsE=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_internal",
    importpath = "modernc.org/internal",
    sum = "h1:XMDsFDcBDsibbBnHB2xzljZ+B1yrOVLEFkKL2u15Glw=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_lldb",
    importpath = "modernc.org/lldb",
    sum = "h1:6vjDJxQEfhlOLwl4bhpwIz00uyFK4EmSYcbwqwbynsc=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_mathutil",
    importpath = "modernc.org/mathutil",
    sum = "h1:93vKjrJopTPrtTNpZ8XIovER7iCIH1QU7wNbOQXC60I=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_ql",
    importpath = "modernc.org/ql",
    sum = "h1:bIQ/trWNVjQPlinI6jdOQsi195SIturGo3mp5hsDqVU=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_sortutil",
    importpath = "modernc.org/sortutil",
    sum = "h1:oP3U4uM+NT/qBQcbg/K2iqAX0Nx7B1b6YZtq3Gk/PjM=",
    version = "v1.1.0",
)

go_repository(
    name = "org_modernc_strutil",
    importpath = "modernc.org/strutil",
    sum = "h1:+1/yCzZxY2pZwwrsbH+4T7BQMoLQ9QiBshRC9eicYsc=",
    version = "v1.1.0",
)

go_repository(
    name = "org_modernc_zappy",
    importpath = "modernc.org/zappy",
    sum = "h1:dPVaP+3ueIUv4guk8PuZ2wiUGcJ1WUVvIheeSSTD0yk=",
    version = "v1.0.0",
)

go_repository(
    name = "org_mongodb_go_mongo_driver",
    importpath = "go.mongodb.org/mongo-driver",
    sum = "h1:aeOqSrhl9eDRAap/3T5pCfMBEBxZ0vuXBP+RMtp2KX8=",
    version = "v1.1.0",
)

go_repository(
    name = "org_uber_go_atomic",
    importpath = "go.uber.org/atomic",
    sum = "h1:OI5t8sDa1Or+q8AeE+yKeB/SDYioSHAgcVljj9JIETY=",
    version = "v1.5.0",
)

go_repository(
    name = "org_uber_go_multierr",
    importpath = "go.uber.org/multierr",
    sum = "h1:sFPn2GLc3poCkfrpIXGhBD2X0CMIo4Q/zSULXrj/+uc=",
    version = "v1.3.0",
)

go_repository(
    name = "org_uber_go_zap",
    importpath = "go.uber.org/zap",
    sum = "h1:nR6NoDBgAf67s68NhaXbsojM+2gxp3S1hWkHDl27pVU=",
    version = "v1.13.0",
)

go_repository(
    name = "tools_gotest_v3",
    importpath = "gotest.tools/v3",
    sum = "h1:kG1BFyqVHuQoVQiR1bWGnfz/fmHvvuiSPIV7rvl360E=",
    version = "v3.0.2",
)

go_repository(
    name = "com_github_lann_builder",
    importpath = "github.com/lann/builder",
    sum = "h1:SOEGU9fKiNWd/HOJuq6+3iTQz8KNCLtVX6idSoTLdUw=",
    version = "v0.0.0-20180802200727-47ae307949d0",
)

go_repository(
    name = "com_github_lann_ps",
    importpath = "github.com/lann/ps",
    sum = "h1:P6pPBnrTSX3DEVR4fDembhRWSsG5rVo6hYhAB/ADZrk=",
    version = "v0.0.0-20150810152359-62de8c46ede0",
)

go_repository(
    name = "com_github_masterminds_squirrel",
    importpath = "github.com/Masterminds/squirrel",
    sum = "h1:he5i/EXixZxrBUWcxzDYMiju9WZ3ld/l7QBNuo/eN3w=",
    version = "v1.4.0",
)

go_repository(
    name = "com_github_desertbit_timer",
    importpath = "github.com/desertbit/timer",
    sum = "h1:U5y3Y5UE0w7amNe7Z5G/twsBW0KEalRQXZzf8ufSh9I=",
    version = "v0.0.0-20180107155436-c41aec40b27f",
)

go_repository(
    name = "com_github_gorilla_websocket",
    importpath = "github.com/gorilla/websocket",
    sum = "h1:+/TMaTYc4QFitKJxsQ7Yye35DkWvkdLcvGKqM+x0Ufc=",
    version = "v1.4.2",
)

go_repository(
    name = "com_github_rs_cors",
    importpath = "github.com/rs/cors",
    sum = "h1:+88SsELBHx5r+hZ8TCkggzSstaWNbDvThkVK8H6f9ik=",
    version = "v1.7.0",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest",
    importpath = "github.com/Azure/go-autorest/autorest",
    sum = "h1:iWJqGEvip7mjibEqC/srXNdo+4wLEPiwlP/7dZLtoPc=",
    version = "v0.11.4",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest_adal",
    importpath = "github.com/Azure/go-autorest/autorest/adal",
    sum = "h1:Aze/GQeAN1RRbGmnUJvUj+tFGBzFdIg3293/A9rbxC4=",
    version = "v0.9.2",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest_date",
    importpath = "github.com/Azure/go-autorest/autorest/date",
    sum = "h1:7gUk1U5M/CQbp9WoqinNzJar+8KY+LPI6wiWrP/myHw=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest_mocks",
    importpath = "github.com/Azure/go-autorest/autorest/mocks",
    sum = "h1:K0laFcLE6VLTOwNgSxaGbUcLPuGXlNkbVvq4cW4nIHk=",
    version = "v0.4.1",
)

go_repository(
    name = "com_github_azure_go_autorest_logger",
    importpath = "github.com/Azure/go-autorest/logger",
    sum = "h1:e4RVHVZKC5p6UANLJHkM4OfR1UKZPj8Wt8Pcx+3oqrE=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_azure_go_autorest_tracing",
    importpath = "github.com/Azure/go-autorest/tracing",
    sum = "h1:TYi4+3m5t6K48TGI9AUdb+IzbnSxvnvUMfuitfgcfuo=",
    version = "v0.6.0",
)

go_repository(
    name = "com_github_docker_spdystream",
    importpath = "github.com/docker/spdystream",
    sum = "h1:cenwrSVm+Z7QLSV/BsnenAOcDXdX4cMv4wP0B/5QbPg=",
    version = "v0.0.0-20160310174837-449fdfce4d96",
)

go_repository(
    name = "com_github_docopt_docopt_go",
    importpath = "github.com/docopt/docopt-go",
    sum = "h1:bWDMxwH3px2JBh6AyO7hdCn/PkvCZXii8TGj7sbtEbQ=",
    version = "v0.0.0-20180111231733-ee0de3bc6815",
)

go_repository(
    name = "com_github_elazarl_goproxy",
    importpath = "github.com/elazarl/goproxy",
    sum = "h1:yUdfgN0XgIJw7foRItutHYUIhlcKzcSf5vDpdhQAKTc=",
    version = "v0.0.0-20180725130230-947c36da3153",
)

go_repository(
    name = "com_github_emicklei_go_restful",
    importpath = "github.com/emicklei/go-restful",
    sum = "h1:spTtZBk5DYEvbxMVutUuTyh1Ao2r4iyvLdACqsl/Ljk=",
    version = "v2.9.5+incompatible",
)

go_repository(
    name = "com_github_evanphx_json_patch",
    importpath = "github.com/evanphx/json-patch",
    sum = "h1:fUDGZCv/7iAN7u0puUVhvKCcsR6vRfwrJatElLBEf0I=",
    version = "v4.2.0+incompatible",
)

go_repository(
    name = "com_github_go_logr_logr",
    importpath = "github.com/go-logr/logr",
    sum = "h1:QvGt2nLcHH0WK9orKa+ppBPAxREcH364nPUedEpK0TY=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_go_openapi_jsonpointer",
    importpath = "github.com/go-openapi/jsonpointer",
    sum = "h1:gihV7YNZK1iK6Tgwwsxo2rJbD1GTbdm72325Bq8FI3w=",
    version = "v0.19.3",
)

go_repository(
    name = "com_github_go_openapi_jsonreference",
    importpath = "github.com/go-openapi/jsonreference",
    sum = "h1:5cxNfTy0UVC3X8JL5ymxzyoUZmo8iZb+jeTWn7tUa8o=",
    version = "v0.19.3",
)

go_repository(
    name = "com_github_go_openapi_spec",
    importpath = "github.com/go-openapi/spec",
    sum = "h1:0XRyw8kguri6Yw4SxhsQA/atC88yqrk0+G4YhI2wabc=",
    version = "v0.19.3",
)

go_repository(
    name = "com_github_go_openapi_swag",
    importpath = "github.com/go-openapi/swag",
    sum = "h1:lTz6Ys4CmqqCQmZPBlbQENR1/GucA2bzYTE12Pw4tFY=",
    version = "v0.19.5",
)

go_repository(
    name = "com_github_googleapis_gnostic",
    importpath = "github.com/googleapis/gnostic",
    sum = "h1:DLJCy1n/vrD4HPjOvYcT8aYQXpPIzoRZONaYwyycI+I=",
    version = "v0.4.1",
)

go_repository(
    name = "com_github_gregjones_httpcache",
    importpath = "github.com/gregjones/httpcache",
    sum = "h1:pdN6V1QBWetyv/0+wjACpqVH+eVULgEjkurDLq3goeM=",
    version = "v0.0.0-20180305231024-9cad4c3443a7",
)

go_repository(
    name = "com_github_imdario_mergo",
    importpath = "github.com/imdario/mergo",
    sum = "h1:CGgOkSJeqMRmt0D9XLWExdT4m4F1vd3FV3VPt+0VxkQ=",
    version = "v0.3.8",
)

go_repository(
    name = "com_github_mailru_easyjson",
    importpath = "github.com/mailru/easyjson",
    sum = "h1:aizVhC/NAAcKWb+5QsU1iNOZb4Yws5UO2I+aIprQITM=",
    version = "v0.7.0",
)

go_repository(
    name = "com_github_munnerz_goautoneg",
    importpath = "github.com/munnerz/goautoneg",
    sum = "h1:7PxY7LVfSZm7PEeBTyK1rj1gABdCO2mbri6GKO1cMDs=",
    version = "v0.0.0-20120707110453-a547fc61f48d",
)

go_repository(
    name = "com_github_mxk_go_flowrate",
    importpath = "github.com/mxk/go-flowrate",
    sum = "h1:y5//uYreIhSUg3J1GEMiLbxo1LJaP8RfCpH6pymGZus=",
    version = "v0.0.0-20140419014527-cca7078d478f",
)

go_repository(
    name = "com_github_nytimes_gziphandler",
    importpath = "github.com/NYTimes/gziphandler",
    sum = "h1:iLrQrdwjDd52kHDA5op2UBJFjmOb9g+7scBan4RN8F0=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_peterbourgon_diskv",
    importpath = "github.com/peterbourgon/diskv",
    sum = "h1:UBdAOUP5p4RWqPBg048CAvpKN+vxiaj6gdUUzhl4XmI=",
    version = "v2.0.1+incompatible",
)

go_repository(
    name = "com_github_puerkitobio_purell",
    importpath = "github.com/PuerkitoBio/purell",
    sum = "h1:WEQqlqaGbrPkxLJWfBwQmfEAE1Z7ONdDLqrN38tNFfI=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_puerkitobio_urlesc",
    importpath = "github.com/PuerkitoBio/urlesc",
    sum = "h1:d+Bc7a5rLufV/sSk/8dngufqelfh6jnri85riMAaF/M=",
    version = "v0.0.0-20170810143723-de5bf2ad4578",
)

go_repository(
    name = "com_github_spf13_afero",
    importpath = "github.com/spf13/afero",
    sum = "h1:5jhuqJyZCZf2JRofRvN/nIFgIWNzPa3/Vz8mYylgbWc=",
    version = "v1.2.2",
)

go_repository(
    name = "cc_mvdan_interfacer",
    importpath = "mvdan.cc/interfacer",
    sum = "h1:WX1yoOaKQfddO/mLzdV4wptyWgoH/6hwLs7QHTixo0I=",
    version = "v0.0.0-20180901003855-c20040233aed",
)

go_repository(
    name = "cc_mvdan_lint",
    importpath = "mvdan.cc/lint",
    sum = "h1:DxJ5nJdkhDlLok9K6qO+5290kphDJbHOQO1DFFFTeBo=",
    version = "v0.0.0-20170908181259-adc824a0674b",
)

go_repository(
    name = "cc_mvdan_unparam",
    importpath = "mvdan.cc/unparam",
    sum = "h1:ujHaqw4YdQP41+lVP79aaouzIU8ZIRMT9/S3YUru+JI=",
    version = "v0.0.0-20190917161559-b83a221c10a2",
)

go_repository(
    name = "com_github_armon_consul_api",
    importpath = "github.com/armon/consul-api",
    sum = "h1:G1bPvciwNyF7IUmKXNt9Ak3m6u9DE1rF+RmtIkBpVdA=",
    version = "v0.0.0-20180202201655-eb2c6b5be1b6",
)

go_repository(
    name = "com_github_benbjohnson_clock",
    importpath = "github.com/benbjohnson/clock",
    sum = "h1:vkLuvpK4fmtSCuo60+yC63p7y0BmQ8gm5ZXGuBCJyXg=",
    version = "v1.0.3",
)

go_repository(
    name = "com_github_bombsimon_wsl",
    importpath = "github.com/bombsimon/wsl",
    sum = "h1:9gTOkIwVtoDZywvX802SDHokeX4kW1cKnV8ZTVAPkRs=",
    version = "v1.2.5",
)

go_repository(
    name = "com_github_cespare_xxhash",
    importpath = "github.com/cespare/xxhash",
    sum = "h1:a6HrQnmkObjyL+Gs60czilIUGqrzKutQD6XZog3p+ko=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_coreos_bbolt",
    importpath = "github.com/coreos/bbolt",
    sum = "h1:wZwiHHUieZCquLkDL0B8UhzreNWsPHooDAG3q34zk0s=",
    version = "v1.3.2",
)

go_repository(
    name = "com_github_coreos_etcd",
    importpath = "github.com/coreos/etcd",
    sum = "h1:jFneRYjIvLMLhDLCzuTuU4rSJUjRplcJQ7pD7MnhC04=",
    version = "v3.3.10+incompatible",
)

go_repository(
    name = "com_github_coreos_go_etcd",
    importpath = "github.com/coreos/go-etcd",
    sum = "h1:bXhRBIXoTm9BYHS3gE0TtQuyNZyeEMux2sDi4oo5YOo=",
    version = "v2.0.0+incompatible",
)

go_repository(
    name = "com_github_coreos_go_semver",
    importpath = "github.com/coreos/go-semver",
    sum = "h1:wkHLiw0WNATZnSG7epLsujiMCgPAc9xhjJ4tgnAxmfM=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_coreos_pkg",
    importpath = "github.com/coreos/pkg",
    sum = "h1:lBNOc5arjvs8E5mO2tbpBpLoyyu8B6e44T7hJy6potg=",
    version = "v0.0.0-20180928190104-399ea9e2e55f",
)

go_repository(
    name = "com_github_cpuguy83_go_md2man",
    importpath = "github.com/cpuguy83/go-md2man",
    sum = "h1:BSKMNlYxDvnunlTymqtgONjNnaRV1sTpcovwwjF22jk=",
    version = "v1.0.10",
)

go_repository(
    name = "com_github_datadog_sketches_go",
    importpath = "github.com/DataDog/sketches-go",
    sum = "h1:RtG+76WKgZuz6FIaGsjoPePmadDBkuD/KC6+ZWu78b8=",
    version = "v0.0.1",
)

go_repository(
    name = "com_github_dgryski_go_sip13",
    importpath = "github.com/dgryski/go-sip13",
    sum = "h1:RMLoZVzv4GliuWafOuPuQDKSm1SJph7uCRnnS61JAn4=",
    version = "v0.0.0-20181026042036-e10d5fee7954",
)

go_repository(
    name = "com_github_fatih_color",
    importpath = "github.com/fatih/color",
    sum = "h1:8xPHl4/q1VyqGIPif1F+1V3Y3lSmrq01EabUW3CoW5s=",
    version = "v1.9.0",
)

go_repository(
    name = "com_github_go_critic_go_critic",
    importpath = "github.com/go-critic/go-critic",
    sum = "h1:GYXWx7Vr3+zv833u+8IoXbNnQY0AdXsxAgI0kX7xcwA=",
    version = "v0.3.5-0.20190904082202-d79a9f0c64db",
)

go_repository(
    name = "com_github_go_lintpack_lintpack",
    importpath = "github.com/go-lintpack/lintpack",
    sum = "h1:DI5mA3+eKdWeJ40nU4d6Wc26qmdG8RCi/btYq0TuRN0=",
    version = "v0.5.2",
)

go_repository(
    name = "com_github_go_ole_go_ole",
    importpath = "github.com/go-ole/go-ole",
    sum = "h1:2lOsA72HgjxAuMlKpFiCbHTvu44PIVkZ5hqm3RSdI/E=",
    version = "v1.2.1",
)

go_repository(
    name = "com_github_go_toolsmith_astcast",
    importpath = "github.com/go-toolsmith/astcast",
    sum = "h1:JojxlmI6STnFVG9yOImLeGREv8W2ocNUM+iOhR6jE7g=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_go_toolsmith_astcopy",
    importpath = "github.com/go-toolsmith/astcopy",
    sum = "h1:OMgl1b1MEpjFQ1m5ztEO06rz5CUd3oBv9RF7+DyvdG8=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_go_toolsmith_astequal",
    importpath = "github.com/go-toolsmith/astequal",
    sum = "h1:4zxD8j3JRFNyLN46lodQuqz3xdKSrur7U/sr0SDS/gQ=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_go_toolsmith_astfmt",
    importpath = "github.com/go-toolsmith/astfmt",
    sum = "h1:A0vDDXt+vsvLEdbMFJAUBI/uTbRw1ffOPnxsILnFL6k=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_go_toolsmith_astinfo",
    importpath = "github.com/go-toolsmith/astinfo",
    sum = "h1:wP6mXeB2V/d1P1K7bZ5vDUO3YqEzcvOREOxZPEu3gVI=",
    version = "v0.0.0-20180906194353-9809ff7efb21",
)

go_repository(
    name = "com_github_go_toolsmith_astp",
    importpath = "github.com/go-toolsmith/astp",
    sum = "h1:alXE75TXgcmupDsMK1fRAy0YUzLzqPVvBKoyWV+KPXg=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_go_toolsmith_pkgload",
    importpath = "github.com/go-toolsmith/pkgload",
    sum = "h1:4DFWWMXVfbcN5So1sBNW9+yeiMqLFGl1wFLTL5R0Tgg=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_go_toolsmith_strparse",
    importpath = "github.com/go-toolsmith/strparse",
    sum = "h1:Vcw78DnpCAKlM20kSbAyO4mPfJn/lyYA4BJUDxe2Jb4=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_go_toolsmith_typep",
    importpath = "github.com/go-toolsmith/typep",
    sum = "h1:zKymWyA1TRYvqYrYDrfEMZULyrhcnGY3x7LDKU2XQaA=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_gobwas_glob",
    importpath = "github.com/gobwas/glob",
    sum = "h1:A4xDbljILXROh+kObIiy5kIaPYD8e96x1tgBhUI5J+Y=",
    version = "v0.2.3",
)

go_repository(
    name = "com_github_gofrs_flock",
    importpath = "github.com/gofrs/flock",
    sum = "h1:ekuhfTjngPhisSjOJ0QWKpPQE8/rbknHaes6WVJj5Hw=",
    version = "v0.0.0-20190320160742-5135e617513b",
)

go_repository(
    name = "com_github_golangci_check",
    importpath = "github.com/golangci/check",
    sum = "h1:23T5iq8rbUYlhpt5DB4XJkc6BU31uODLD1o1gKvZmD0=",
    version = "v0.0.0-20180506172741-cfe4005ccda2",
)

go_repository(
    name = "com_github_golangci_dupl",
    importpath = "github.com/golangci/dupl",
    sum = "h1:w8hkcTqaFpzKqonE9uMCefW1WDie15eSP/4MssdenaM=",
    version = "v0.0.0-20180902072040-3e9179ac440a",
)

go_repository(
    name = "com_github_golangci_errcheck",
    importpath = "github.com/golangci/errcheck",
    sum = "h1:YYWNAGTKWhKpcLLt7aSj/odlKrSrelQwlovBpDuf19w=",
    version = "v0.0.0-20181223084120-ef45e06d44b6",
)

go_repository(
    name = "com_github_golangci_go_misc",
    importpath = "github.com/golangci/go-misc",
    sum = "h1:9kfjN3AdxcbsZBf8NjltjWihK2QfBBBZuv91cMFfDHw=",
    version = "v0.0.0-20180628070357-927a3d87b613",
)

go_repository(
    name = "com_github_golangci_goconst",
    importpath = "github.com/golangci/goconst",
    sum = "h1:pe9JHs3cHHDQgOFXJJdYkK6fLz2PWyYtP4hthoCMvs8=",
    version = "v0.0.0-20180610141641-041c5f2b40f3",
)

go_repository(
    name = "com_github_golangci_gocyclo",
    importpath = "github.com/golangci/gocyclo",
    sum = "h1:pXTK/gkVNs7Zyy7WKgLXmpQ5bHTrq5GDsp8R9Qs67g0=",
    version = "v0.0.0-20180528144436-0a533e8fa43d",
)

go_repository(
    name = "com_github_golangci_gofmt",
    importpath = "github.com/golangci/gofmt",
    sum = "h1:iR3fYXUjHCR97qWS8ch1y9zPNsgXThGwjKPrYfqMPks=",
    version = "v0.0.0-20190930125516-244bba706f1a",
)

go_repository(
    name = "com_github_golangci_golangci_lint",
    importpath = "github.com/golangci/golangci-lint",
    sum = "h1:HxAxpR8Z0M8omihvQdsD3PF0qPjlqYqp2vMJzstoKeI=",
    version = "v1.21.0",
)

go_repository(
    name = "com_github_golangci_ineffassign",
    importpath = "github.com/golangci/ineffassign",
    sum = "h1:gLLhTLMk2/SutryVJ6D4VZCU3CUqr8YloG7FPIBWFpI=",
    version = "v0.0.0-20190609212857-42439a7714cc",
)

go_repository(
    name = "com_github_golangci_maligned",
    importpath = "github.com/golangci/maligned",
    sum = "h1:kNY3/svz5T29MYHubXix4aDDuE3RWHkPvopM/EDv/MA=",
    version = "v0.0.0-20180506175553-b1d89398deca",
)

go_repository(
    name = "com_github_golangci_misspell",
    importpath = "github.com/golangci/misspell",
    sum = "h1:EL/O5HGrF7Jaq0yNhBLucz9hTuRzj2LdwGBOaENgxIk=",
    version = "v0.0.0-20180809174111-950f5d19e770",
)

go_repository(
    name = "com_github_golangci_prealloc",
    importpath = "github.com/golangci/prealloc",
    sum = "h1:leSNB7iYzLYSSx3J/s5sVf4Drkc68W2wm4Ixh/mr0us=",
    version = "v0.0.0-20180630174525-215b22d4de21",
)

go_repository(
    name = "com_github_golangci_revgrep",
    importpath = "github.com/golangci/revgrep",
    sum = "h1:XQKc8IYQOeRwVs36tDrEmTgDgP88d5iEURwpmtiAlOM=",
    version = "v0.0.0-20180812185044-276a5c0a1039",
)

go_repository(
    name = "com_github_golangci_unconvert",
    importpath = "github.com/golangci/unconvert",
    sum = "h1:zwtduBRr5SSWhqsYNgcuWO2kFlpdOZbP0+yRjmvPGys=",
    version = "v0.0.0-20180507085042-28b1c447d1f4",
)

go_repository(
    name = "com_github_gostaticanalysis_analysisutil",
    importpath = "github.com/gostaticanalysis/analysisutil",
    sum = "h1:iwp+5/UAyzQSFgQ4uR2sni99sJ8Eo9DEacKWM5pekIg=",
    version = "v0.0.3",
)

go_repository(
    name = "com_github_grpc_ecosystem_go_grpc_middleware",
    importpath = "github.com/grpc-ecosystem/go-grpc-middleware",
    sum = "h1:MqvH60+R2JhSdvVgGxmExOndrkRQtGW7w4+gcrymN64=",
    version = "v1.2.1-0.20200228141219-3ce3d519df39",
)

go_repository(
    name = "com_github_grpc_ecosystem_go_grpc_prometheus",
    importpath = "github.com/grpc-ecosystem/go-grpc-prometheus",
    sum = "h1:Ovs26xHkKqVztRpIrF/92BcuyuQ/YW4NSIpoGtfXNho=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_grpc_ecosystem_grpc_gateway",
    importpath = "github.com/grpc-ecosystem/grpc-gateway",
    sum = "h1:UImYN5qQ8tuGpGE16ZmjvcTtTw24zw1QAp/SlnNrZhI=",
    version = "v1.9.5",
)

go_repository(
    name = "com_github_hashicorp_hcl",
    importpath = "github.com/hashicorp/hcl",
    sum = "h1:8B4HqvMUtYSjsGkYjiQGStc9pXffY2J+Z2SPQAj+wMY=",
    version = "v1.0.1-0.20201016140508-a07e7d50bbee",
)

go_repository(
    name = "com_github_inconshreveable_mousetrap",
    importpath = "github.com/inconshreveable/mousetrap",
    sum = "h1:Z8tu5sraLXCXIcARxBp/8cbvlwVa7Z1NHg9XEKhtSvM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_jonboulle_clockwork",
    importpath = "github.com/jonboulle/clockwork",
    sum = "h1:UOGuzwb1PwsrDAObMuhUnj0p5ULPj8V/xJ7Kx9qUBdQ=",
    version = "v0.2.2",
)

go_repository(
    name = "com_github_klauspost_compress",
    importpath = "github.com/klauspost/compress",
    sum = "h1:MiK62aErc3gIiVEtyzKfeOHgW7atJb5g/KNX5m3c2nQ=",
    version = "v1.11.2",
)

go_repository(
    name = "com_github_klauspost_cpuid",
    importpath = "github.com/klauspost/cpuid",
    sum = "h1:NMpwD2G9JSFOE1/TJjGSo5zG7Yb2bTe7eq1jH+irmeE=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_logrusorgru_aurora",
    importpath = "github.com/logrusorgru/aurora",
    sum = "h1:9MlwzLdW7QSDrhDjFlsEYmxpFyIoXmYRon3dt0io31k=",
    version = "v0.0.0-20181002194514-a7b3b318ed4e",
)

go_repository(
    name = "com_github_magiconair_properties",
    importpath = "github.com/magiconair/properties",
    sum = "h1:LLgXmsheXeRoUOBOjtwPQCWIYqM/LU1ayDtDePerRcY=",
    version = "v1.8.0",
)

go_repository(
    name = "com_github_matoous_godox",
    importpath = "github.com/matoous/godox",
    sum = "h1:RHba4YImhrUVQDHUCe2BNSOz4tVy2yGyXhvYDvxGgeE=",
    version = "v0.0.0-20190911065817-5d6d842e92eb",
)

go_repository(
    name = "com_github_mattn_goveralls",
    importpath = "github.com/mattn/goveralls",
    sum = "h1:7eJB6EqsPhRVxvwEXGnqdO2sJI0PTsrWoTMXEk9/OQc=",
    version = "v0.0.2",
)

go_repository(
    name = "com_github_mitchellh_go_homedir",
    importpath = "github.com/mitchellh/go-homedir",
    sum = "h1:lukF9ziXFxDFPkA1vsr5zpc1XuPDn/wFntq5mG+4E0Y=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_mitchellh_go_ps",
    importpath = "github.com/mitchellh/go-ps",
    sum = "h1:9+ke9YJ9KGWw5ANXK6ozjoK47uI3uNbXv4YVINBnGm8=",
    version = "v0.0.0-20190716172923-621e5597135b",
)

go_repository(
    name = "com_github_mitchellh_mapstructure",
    importpath = "github.com/mitchellh/mapstructure",
    sum = "h1:SzB1nHZ2Xi+17FP0zVQBHIZqvwRN9408fJO8h+eeNA8=",
    version = "v1.3.3",
)

go_repository(
    name = "com_github_mozilla_tls_observatory",
    importpath = "github.com/mozilla/tls-observatory",
    sum = "h1:Av0AX0PnAlPZ3AY2rQUobGFaZfE4KHVRdKWIEPvsCWY=",
    version = "v0.0.0-20190404164649-a3c1b6cfecfd",
)

go_repository(
    name = "com_github_nbutton23_zxcvbn_go",
    importpath = "github.com/nbutton23/zxcvbn-go",
    sum = "h1:AREM5mwr4u1ORQBMvzfzBgpsctsbQikCVpvC+tX285E=",
    version = "v0.0.0-20180912185939-ae427f1e4c1d",
)

go_repository(
    name = "com_github_oklog_ulid",
    importpath = "github.com/oklog/ulid",
    sum = "h1:EGfNDEx6MqHz8B3uNV6QAib1UR2Lm97sHi3ocA6ESJ4=",
    version = "v1.3.1",
)

go_repository(
    name = "com_github_oneofone_xxhash",
    importpath = "github.com/OneOfOne/xxhash",
    sum = "h1:KMrpdQIwFcEqXDklaen+P1axHaj9BSKzvpUUfnHldSE=",
    version = "v1.2.2",
)

go_repository(
    name = "com_github_openpeedeep_depguard",
    importpath = "github.com/OpenPeeDeeP/depguard",
    sum = "h1:VlW4R6jmBIv3/u1JNlawEvJMM4J+dPORPaZasQee8Us=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_opentracing_opentracing_go",
    importpath = "github.com/opentracing/opentracing-go",
    sum = "h1:pWlfV3Bxv7k65HYwkikxat0+s3pV4bsqf19k25Ur8rU=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_pelletier_go_toml",
    importpath = "github.com/pelletier/go-toml",
    sum = "h1:T5zMGML61Wp+FlcbWjRDT7yAxhJNAiPPLOFECq181zc=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_prometheus_tsdb",
    importpath = "github.com/prometheus/tsdb",
    sum = "h1:YZcsG11NqnK4czYLrWd9mpEuAJIHVQLwdrleYfszMAA=",
    version = "v0.7.1",
)

go_repository(
    name = "com_github_quasilyte_go_consistent",
    importpath = "github.com/quasilyte/go-consistent",
    sum = "h1:JoUA0uz9U0FVFq5p4LjEq4C0VgQ0El320s3Ms0V4eww=",
    version = "v0.0.0-20190521200055-c6f3937de18c",
)

go_repository(
    name = "com_github_rogpeppe_fastuuid",
    importpath = "github.com/rogpeppe/fastuuid",
    sum = "h1:gu+uRPtBe88sKxUCEXRoeCvVG90TJmwhiqRpvdhQFng=",
    version = "v0.0.0-20150106093220-6724a57986af",
)

go_repository(
    name = "com_github_russross_blackfriday",
    importpath = "github.com/russross/blackfriday",
    sum = "h1:HyvC0ARfnZBqnXwABFeSZHpKvJHJJfPz81GNueLj0oo=",
    version = "v1.5.2",
)

go_repository(
    name = "com_github_securego_gosec",
    importpath = "github.com/securego/gosec",
    sum = "h1:qlgneduvrO2HHdK1CcP+P1vlTAkT3J64A+1qd+PobNg=",
    version = "v0.0.0-20191008095658-28c1128b7336",
)

go_repository(
    name = "com_github_shirou_gopsutil",
    importpath = "github.com/shirou/gopsutil",
    sum = "h1:msXs2frUV+O/JLva9EDLpuJ84PrFsdCTCQex8PUdtkQ=",
    version = "v2.20.9+incompatible",
)

go_repository(
    name = "com_github_shirou_w32",
    importpath = "github.com/shirou/w32",
    sum = "h1:udFKJ0aHUL60LboW/A+DfgoHVedieIzIXE8uylPue0U=",
    version = "v0.0.0-20160930032740-bb4de0191aa4",
)

go_repository(
    name = "com_github_shurcool_go",
    importpath = "github.com/shurcooL/go",
    sum = "h1:MZM7FHLqUHYI0Y/mQAt3d2aYa0SiNms/hFqC9qJYolM=",
    version = "v0.0.0-20180423040247-9e1955d9fb6e",
)

go_repository(
    name = "com_github_shurcool_go_goon",
    importpath = "github.com/shurcooL/go-goon",
    sum = "h1:llrF3Fs4018ePo4+G/HV/uQUqEI1HMDjCeOf2V6puPc=",
    version = "v0.0.0-20170922171312-37c2f522c041",
)

go_repository(
    name = "com_github_soheilhy_cmux",
    importpath = "github.com/soheilhy/cmux",
    sum = "h1:0HKaf1o97UwFjHH9o5XsHUOF+tqmdA7KEzXLpiyaw0E=",
    version = "v0.1.4",
)

go_repository(
    name = "com_github_sourcegraph_go_diff",
    importpath = "github.com/sourcegraph/go-diff",
    sum = "h1:gO6i5zugwzo1RVTvgvfwCOSVegNuvnNi6bAD1QCmkHs=",
    version = "v0.5.1",
)

go_repository(
    name = "com_github_spaolacci_murmur3",
    importpath = "github.com/spaolacci/murmur3",
    sum = "h1:qLC7fQah7D6K1B0ujays3HV9gkFtllcxhzImRR7ArPQ=",
    version = "v0.0.0-20180118202830-f09979ecbc72",
)

go_repository(
    name = "com_github_spf13_cast",
    importpath = "github.com/spf13/cast",
    sum = "h1:oget//CVOEoFewqQxwr0Ej5yjygnqGkvggSE/gB35Q8=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_spf13_cobra",
    importpath = "github.com/spf13/cobra",
    sum = "h1:6m/oheQuQ13N9ks4hubMG6BnvwOeaJrqSPLahSnczz8=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_spf13_jwalterweatherman",
    importpath = "github.com/spf13/jwalterweatherman",
    sum = "h1:XHEdyB+EcvlqZamSM4ZOMGlc93t6AcsBEu9Gc1vn7yk=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_spf13_viper",
    importpath = "github.com/spf13/viper",
    sum = "h1:yXHLWeravcrgGyFSyCgdYpXQ9dR9c/WED3pg1RhxqEU=",
    version = "v1.4.0",
)

go_repository(
    name = "com_github_stackexchange_wmi",
    importpath = "github.com/StackExchange/wmi",
    sum = "h1:fLjPD/aNc3UIOA6tDi6QXUemppXK3P9BI7mr2hd6gx8=",
    version = "v0.0.0-20180116203802-5d049714c4a6",
)

go_repository(
    name = "com_github_timakin_bodyclose",
    importpath = "github.com/timakin/bodyclose",
    sum = "h1:RumXZ56IrCj4CL+g1b9OL/oH0QnsF976bC8xQFYUD5Q=",
    version = "v0.0.0-20190930140734-f7f2e9bca95e",
)

go_repository(
    name = "com_github_tmc_grpc_websocket_proxy",
    importpath = "github.com/tmc/grpc-websocket-proxy",
    sum = "h1:LnC5Kc/wtumK+WB441p7ynQJzVuNRJiqddSIE3IlSEQ=",
    version = "v0.0.0-20190109142713-0ad062ec5ee5",
)

go_repository(
    name = "com_github_ugorji_go",
    importpath = "github.com/ugorji/go",
    sum = "h1:j4s+tAvLfL3bZyefP2SEWmhBzmuIlH/eqNuPdFPgngw=",
    version = "v1.1.4",
)

go_repository(
    name = "com_github_ugorji_go_codec",
    importpath = "github.com/ugorji/go/codec",
    sum = "h1:3SVOIvH7Ae1KRYyQWRjXWJEA9sS/c/pjvH++55Gr648=",
    version = "v0.0.0-20181204163529-d75b2dcb6bc8",
)

go_repository(
    name = "com_github_ultraware_funlen",
    importpath = "github.com/ultraware/funlen",
    sum = "h1:Av96YVBwwNSe4MLR7iI/BIa3VyI7/djnto/pK3Uxbdo=",
    version = "v0.0.2",
)

go_repository(
    name = "com_github_ultraware_whitespace",
    importpath = "github.com/ultraware/whitespace",
    sum = "h1:If7Va4cM03mpgrNH9k49/VOicWpGoG70XPBFFODYDsg=",
    version = "v0.0.4",
)

go_repository(
    name = "com_github_uudashr_gocognit",
    importpath = "github.com/uudashr/gocognit",
    sum = "h1:NST9SQhYHpFiyKdXn8UAiogdB1xkt5RnA1/rR/oBcGw=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_valyala_fasthttp",
    importpath = "github.com/valyala/fasthttp",
    sum = "h1:dzZJf2IuMiclVjdw0kkT+f9u4YdrapbNyGAN47E/qnk=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_valyala_quicktemplate",
    importpath = "github.com/valyala/quicktemplate",
    sum = "h1:BaO1nHTkspYzmAjPXj0QiDJxai96tlcZyKcI9dyEGvM=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_valyala_tcplisten",
    importpath = "github.com/valyala/tcplisten",
    sum = "h1:0R4NLDRDZX6JcmhJgXi5E4b8Wg84ihbmUKp/GvSPEzc=",
    version = "v0.0.0-20161114210144-ceec8f93295a",
)

go_repository(
    name = "com_github_xiang90_probing",
    importpath = "github.com/xiang90/probing",
    sum = "h1:eY9dn8+vbi4tKz5Qo6v2eYzo7kUS51QINcR5jNpbZS8=",
    version = "v0.0.0-20190116061207-43a291ad63a2",
)

go_repository(
    name = "com_github_xordataexchange_crypt",
    importpath = "github.com/xordataexchange/crypt",
    sum = "h1:ESFSdwYZvkeru3RtdrYueztKhOBCSAAzS4Gf+k0tEow=",
    version = "v0.0.3-0.20170626215501-b2862e3d0a77",
)

go_repository(
    name = "com_sourcegraph_sqs_pbtypes",
    importpath = "sourcegraph.com/sqs/pbtypes",
    sum = "h1:f7lAwqviDEGvON4kRv0o5V7FT/IQK+tbkF664XMbP3o=",
    version = "v1.0.0",
)

go_repository(
    name = "in_gopkg_resty_v1",
    importpath = "gopkg.in/resty.v1",
    sum = "h1:CuXP0Pjfw9rOuY6EP+UvtNvt5DSqHpIxILZKT/quCZI=",
    version = "v1.12.0",
)

go_repository(
    name = "in_gopkg_yaml_v3",
    importpath = "gopkg.in/yaml.v3",
    sum = "h1:dUUwHk2QECo/6vqA44rthZ8ie2QXMNeKRTHCNY2nXvo=",
    version = "v3.0.0-20200313102051-9f266ea9e77c",
)

go_repository(
    name = "io_etcd_go_bbolt",
    importpath = "go.etcd.io/bbolt",
    sum = "h1:MUGmc65QhB3pIlaQ5bB4LwqSj6GIonVJXpZiaKNyaKk=",
    version = "v1.3.3",
)

go_repository(
    name = "io_opentelemetry_go_otel",
    build_file_proto_mode = "disable_global",
    importpath = "go.opentelemetry.io/otel",
    sum = "h1:Lenfy7QHRXPZVsw/12CWpxX6d/JkrX8wrx2vO8G80Ng=",
    version = "v0.19.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_exporter_metric_prometheus",
    importpath = "go.opentelemetry.io/otel/exporter/metric/prometheus",
    sum = "h1:MkrRXJOaI9iHbv2k4/9//CgrrvEcpgNjyHHeviYejxc=",
    version = "v0.2.1",
)

go_repository(
    name = "io_opentelemetry_go_otel_sdk",
    importpath = "go.opentelemetry.io/otel/sdk",
    sum = "h1:13pQquZyGbIvGxBWcVzUqe8kg5VGbTBiKKKXpYCylRM=",
    version = "v0.19.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_exporters_metric_prometheus",
    importpath = "go.opentelemetry.io/otel/exporters/metric/prometheus",
    sum = "h1:DMHfiNaNzn0z/uG2goN1fEe+/LSXrd3nBVu/Ag8Ju7M=",
    version = "v0.19.0",
)

go_repository(
    name = "com_github_segmentio_ksuid",
    importpath = "github.com/segmentio/ksuid",
    sum = "h1:FoResxvleQwYiPAVKe1tMUlEirodZqlqglIuFsdDntY=",
    version = "v1.0.3",
)

go_repository(
    name = "io_opentelemetry_go_contrib",
    importpath = "go.opentelemetry.io/contrib",
    sum = "h1:x6Josyb/V+aDHg6IozzmZMaOhE+0Jb2NvEAM4/0Gftc=",
    version = "v0.19.0",
)

go_repository(
    name = "io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc",
    importpath = "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc",
    sum = "h1:zekwSWkeZPKiEQo3tl82RVryxARMXbazgG6pLPzKgn0=",
    version = "v0.19.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_exporters_stdout",
    importpath = "go.opentelemetry.io/otel/exporters/stdout",
    sum = "h1:A+XiGIPQbGoJoBOJfKAKnZyiUSjSWvL3XWETUvtom5k=",
    version = "v0.13.0",
)

go_repository(
    name = "com_github_felixge_httpsnoop",
    importpath = "github.com/felixge/httpsnoop",
    sum = "h1:lvB5Jl89CsZtGIWuTcDM1E/vkVs49/Ml7JJe07l8SPQ=",
    version = "v1.0.1",
)

go_repository(
    name = "io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp",
    importpath = "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
    sum = "h1:HOKafMKQkF8/+m57PrGDgV2OAbWKFKhbb1wbgLZ0+J4=",
    version = "v0.19.0",
)

go_repository(
    name = "com_github_apache_thrift",
    build_extra_args = [
        "-exclude=test",
        "-exclude=tutorial",
        "-exclude=lib/go/test",
    ],
    importpath = "github.com/apache/thrift",
    sum = "h1:5hryIiq9gtn+MiLVn0wP37kb/uTeRZgN08WoCsAhIhI=",
    version = "v0.13.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_exporters_trace_jaeger",
    importpath = "go.opentelemetry.io/otel/exporters/trace/jaeger",
    sum = "h1:qU7sGQoDrlAvNYryR2SJT4ULP/q+5tE38W+h31WEE/M=",
    version = "v0.19.0",
)

go_repository(
    name = "com_zombiezen_go_postgrestest",
    importpath = "zombiezen.com/go/postgrestest",
    sum = "h1:4QKJsCXlbdlRI6AVmqDx+Xn3sqEnZ8TY3nex4MTM5+Q=",
    version = "v1.0.0",
)

go_repository(
    name = "io_opentelemetry_go_contrib_instrumentation_runtime",
    importpath = "go.opentelemetry.io/contrib/instrumentation/runtime",
    sum = "h1:SKXMYj/UXxfO6EUTLSJ1gfvbhuu/sgpJF8B2FeJpppE=",
    version = "v0.19.0",
)

go_repository(
    name = "com_github_agnivade_levenshtein",
    importpath = "github.com/agnivade/levenshtein",
    sum = "h1:3oJU7J3FGFmyhn8KHjmVaZCN5hxTr7GxgRue+sxIXdQ=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_ajstarks_svgo",
    importpath = "github.com/ajstarks/svgo",
    sum = "h1:wVe6/Ea46ZMeNkQjjBW6xcqyQA/j5e0D6GytH95g0gQ=",
    version = "v0.0.0-20180226025133-644b8db467af",
)

go_repository(
    name = "com_github_andreyvit_diff",
    importpath = "github.com/andreyvit/diff",
    sum = "h1:bvNMNQO63//z+xNgfBlViaCIJKLlCJ6/fmUseuG0wVQ=",
    version = "v0.0.0-20170406064948-c7f18ee00883",
)

go_repository(
    name = "com_github_armon_circbuf",
    importpath = "github.com/armon/circbuf",
    sum = "h1:QEF07wC0T1rKkctt1RINW/+RMTVmiwxETico2l3gxJA=",
    version = "v0.0.0-20150827004946-bbbad097214e",
)

go_repository(
    name = "com_github_asaskevich_govalidator",
    importpath = "github.com/asaskevich/govalidator",
    sum = "h1:idn718Q4B6AGu/h5Sxe66HYVdqdGu2l9Iebqhi/AEoA=",
    version = "v0.0.0-20190424111038-f61b66f89f4a",
)

go_repository(
    name = "com_github_auth0_go_jwt_middleware",
    importpath = "github.com/auth0/go-jwt-middleware",
    sum = "h1:irR1cO6eek3n5uquIVaRAsQmZnlsfPuHNz31cXo4eyk=",
    version = "v0.0.0-20170425171159-5493cabe49f7",
)

go_repository(
    name = "com_github_azure_azure_sdk_for_go",
    importpath = "github.com/Azure/azure-sdk-for-go",
    sum = "h1:e82Yv2HNpS0kuyeCrV29OPKvEiqfs2/uJHic3/3iKdg=",
    version = "v44.0.0+incompatible",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest_to",
    importpath = "github.com/Azure/go-autorest/autorest/to",
    sum = "h1:oXVqrxakqqV1UZdSazDOPOLvOIz+XA683u8EctwboHk=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest_validation",
    importpath = "github.com/Azure/go-autorest/autorest/validation",
    sum = "h1:3I9AAI63HfcLtphd9g39ruUwRI+Ca+z/f36KHPFRUss=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_bgentry_speakeasy",
    importpath = "github.com/bgentry/speakeasy",
    sum = "h1:ByYyxL9InA1OWqxJqqp2A5pYHUrCiAL6K3J+LKSsQkY=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_bifurcation_mint",
    importpath = "github.com/bifurcation/mint",
    sum = "h1:fUjoj2bT6dG8LoEe+uNsKk8J+sLkDbQkJnB6Z1F02Bc=",
    version = "v0.0.0-20180715133206-93c51c6ce115",
)

go_repository(
    name = "com_github_blang_semver",
    importpath = "github.com/blang/semver",
    sum = "h1:7hqmJYuaEK3qwVjWubYiht3j93YI0WQBuysxHIfUriU=",
    version = "v3.1.0+incompatible",
)

go_repository(
    name = "com_github_boltdb_bolt",
    importpath = "github.com/boltdb/bolt",
    sum = "h1:JQmyP4ZBrce+ZQu0dY660FMfatumYDLun9hBCUVIkF4=",
    version = "v1.3.1",
)

go_repository(
    name = "com_github_caddyserver_caddy",
    importpath = "github.com/caddyserver/caddy",
    sum = "h1:i9gRhBgvc5ifchwWtSe7pDpsdS9+Q0Rw9oYQmYUTw1w=",
    version = "v1.0.3",
)

go_repository(
    name = "com_github_cenkalti_backoff",
    importpath = "github.com/cenkalti/backoff",
    sum = "h1:tNowT99t7UNflLxfYYSlKYsBpXdEet03Pg2g16Swow4=",
    version = "v2.2.1+incompatible",
)

go_repository(
    name = "com_github_chai2010_gettext_go",
    importpath = "github.com/chai2010/gettext-go",
    sum = "h1:7aWHqerlJ41y6FOsEUvknqgXnGmJyJSbjhAWq5pO4F8=",
    version = "v0.0.0-20160711120539-c6fed771bfd5",
)

go_repository(
    name = "com_github_checkpoint_restore_go_criu_v4",
    importpath = "github.com/checkpoint-restore/go-criu/v4",
    sum = "h1:WW2B2uxx9KWF6bGlHqhm8Okiafwwx7Y2kcpn8lCpjgo=",
    version = "v4.1.0",
)

go_repository(
    name = "com_github_cheekybits_genny",
    importpath = "github.com/cheekybits/genny",
    sum = "h1:a1zrFsLFac2xoM6zG1u72DWJwZG3ayttYLfmLbxVETk=",
    version = "v0.0.0-20170328200008-9127e812e1e9",
)

go_repository(
    name = "com_github_cilium_ebpf",
    importpath = "github.com/cilium/ebpf",
    sum = "h1:cHzBGGVew0ezFsq2grfy2RsB8hO/eNyBgOLHBCqfR1U=",
    version = "v0.0.0-20200702112145-1c8d4c9ef775",
)

go_repository(
    name = "com_github_clusterhq_flocker_go",
    importpath = "github.com/clusterhq/flocker-go",
    sum = "h1:eIHD9GNM3Hp7kcRW5mvcz7WTR3ETeoYYKwpgA04kaXE=",
    version = "v0.0.0-20160920122132-2b8b7259d313",
)

go_repository(
    name = "com_github_cockroachdb_datadriven",
    importpath = "github.com/cockroachdb/datadriven",
    sum = "h1:OaNxuTZr7kxeODyLWsRMC+OD03aFUH+mW6r2d+MWa5Y=",
    version = "v0.0.0-20190809214429-80d97fb3cbaa",
)

go_repository(
    name = "com_github_codegangsta_negroni",
    importpath = "github.com/codegangsta/negroni",
    sum = "h1:+aYywywx4bnKXWvoWtRfJ91vC59NbEhEY03sZjQhbVY=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_container_storage_interface_spec",
    importpath = "github.com/container-storage-interface/spec",
    sum = "h1:m2LYF3fo9IPapVt5FGRVw5bJPmlWqWIezB0jkQh03Zo=",
    version = "v1.2.0-rc1.0.20191021210849-a33ece0a8a9f",
)

go_repository(
    name = "com_github_containerd_cgroups",
    importpath = "github.com/containerd/cgroups",
    sum = "h1:tSNMc+rJDfmYntojat8lljbt1mgKNpTxUZJsSzJ9Y1s=",
    version = "v0.0.0-20190919134610-bf292b21730f",
)

go_repository(
    name = "com_github_containerd_console",
    importpath = "github.com/containerd/console",
    sum = "h1:fU3UuQapBs+zLJu82NhR11Rif1ny2zfMMAyPJzSN5tQ=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_containerd_continuity",
    importpath = "github.com/containerd/continuity",
    sum = "h1:8ahmSVELW1wghbjerVAyuEYD5+Dio66RYvSS0iGfL1M=",
    version = "v0.0.0-20200228182428-0f16d7a0959c",
)

go_repository(
    name = "com_github_containerd_fifo",
    importpath = "github.com/containerd/fifo",
    sum = "h1:PUD50EuOMkXVcpBIA/R95d56duJR9VxhwncsFbNnxW4=",
    version = "v0.0.0-20190226154929-a9fb20d87448",
)

go_repository(
    name = "com_github_containerd_go_runc",
    importpath = "github.com/containerd/go-runc",
    sum = "h1:esQOJREg8nw8aXj6uCN5dfW5cKUBiEJ/+nni1Q/D/sw=",
    version = "v0.0.0-20180907222934-5a6d9f37cfa3",
)

go_repository(
    name = "com_github_containerd_ttrpc",
    importpath = "github.com/containerd/ttrpc",
    sum = "h1:dlfGmNcE3jDAecLqwKPMNX6nk2qh1c1Vg1/YTzpOOF4=",
    version = "v0.0.0-20190828154514-0e0f228740de",
)

go_repository(
    name = "com_github_containerd_typeurl",
    importpath = "github.com/containerd/typeurl",
    sum = "h1:JNn81o/xG+8NEo3bC/vx9pbi/g2WI8mtP2/nXzu297Y=",
    version = "v0.0.0-20180627222232-a93fcdb778cd",
)

go_repository(
    name = "com_github_containernetworking_cni",
    importpath = "github.com/containernetworking/cni",
    sum = "h1:zLuuwrIjPmE9XxixRYl8KBiGUNaN5DpxLZJ39f3BFio=",
    version = "v0.7.2-0.20190612152420-dc953e2fd91f",
)

go_repository(
    name = "com_github_coredns_corefile_migration",
    importpath = "github.com/coredns/corefile-migration",
    sum = "h1:7HI4r5S5Fne749a+JDxUZppqBpYoZK8Q53ZVK9cn3aM=",
    version = "v1.0.10",
)

go_repository(
    name = "com_github_coreos_go_oidc",
    importpath = "github.com/coreos/go-oidc",
    sum = "h1:sdJrfw8akMnCuUlaZU3tE/uYXFgfqom8DBE9so9EBsM=",
    version = "v2.1.0+incompatible",
)

go_repository(
    name = "com_github_coreos_go_systemd_v22",
    importpath = "github.com/coreos/go-systemd/v22",
    sum = "h1:kq/SbG2BCKLkDKkjQf5OWwKWUKj1lgs3lFI4PxnR5lg=",
    version = "v22.1.0",
)

go_repository(
    name = "com_github_cpuguy83_go_md2man_v2",
    importpath = "github.com/cpuguy83/go-md2man/v2",
    sum = "h1:EoUDS0afbrsXAZ9YQ9jdu/mZ2sXgT1/2yyNng4PGlyM=",
    version = "v2.0.0",
)

go_repository(
    name = "com_github_cyphar_filepath_securejoin",
    importpath = "github.com/cyphar/filepath-securejoin",
    sum = "h1:dCqRswe3ZAwkQWdvFLwRqmJCpGP3DWb7bFogdqY3+QU=",
    version = "v0.2.3-0.20190205144030-7efe413b52e1",
)

go_repository(
    name = "com_github_daviddengcn_go_colortext",
    importpath = "github.com/daviddengcn/go-colortext",
    sum = "h1:uVsMphB1eRx7xB1njzL3fuMdWRN8HtVzoUOItHMwv5c=",
    version = "v0.0.0-20160507010035-511bcaf42ccd",
)

go_repository(
    name = "com_github_dnaeon_go_vcr",
    importpath = "github.com/dnaeon/go-vcr",
    sum = "h1:r8L/HqC0Hje5AXMu1ooW8oyQyOFv4GxqpL0nRP7SLLY=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_dustin_go_humanize",
    importpath = "github.com/dustin/go-humanize",
    sum = "h1:VSnTsYCnlFHaM2/igO1h6X3HA71jcobQuxemgkq4zYo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_euank_go_kmsg_parser",
    importpath = "github.com/euank/go-kmsg-parser",
    sum = "h1:cHD53+PLQuuQyLZeriD1V/esuG4MuU0Pjs5y6iknohY=",
    version = "v2.0.0+incompatible",
)

go_repository(
    name = "com_github_exponent_io_jsonpath",
    importpath = "github.com/exponent-io/jsonpath",
    sum = "h1:105gxyaGwCFad8crR9dcMQWvV9Hvulu6hwUh4tWPJnM=",
    version = "v0.0.0-20151013193312-d6023ce2651d",
)

go_repository(
    name = "com_github_fatih_camelcase",
    importpath = "github.com/fatih/camelcase",
    sum = "h1:hxNvNX/xYBp0ovncs8WyWZrOrpBNub/JfaMvbURyft8=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_flynn_go_shlex",
    importpath = "github.com/flynn/go-shlex",
    sum = "h1:BHsljHzVlRcyQhjrss6TZTdY2VfCqZPbv5k3iBFa2ZQ=",
    version = "v0.0.0-20150515145356-3f9db97f8568",
)

go_repository(
    name = "com_github_fogleman_gg",
    importpath = "github.com/fogleman/gg",
    sum = "h1:WXb3TSNmHp2vHoCroCIB1foO/yQ36swABL8aOVeDpgg=",
    version = "v1.2.1-0.20190220221249-0403632d5b90",
)

go_repository(
    name = "com_github_globalsign_mgo",
    importpath = "github.com/globalsign/mgo",
    sum = "h1:DujepqpGd1hyOd7aW59XpK7Qymp8iy83xq74fLr21is=",
    version = "v0.0.0-20181015135952-eeefdecb41b8",
)

go_repository(
    name = "com_github_go_acme_lego",
    importpath = "github.com/go-acme/lego",
    sum = "h1:5fNN9yRQfv8ymH3DSsxla+4aYeQt2IgfZqHKVnK8f0s=",
    version = "v2.5.0+incompatible",
)

go_repository(
    name = "com_github_go_bindata_go_bindata",
    importpath = "github.com/go-bindata/go-bindata",
    sum = "h1:tR4f0e4VTO7LK6B2YWyAoVEzG9ByG1wrXB4TL9+jiYg=",
    version = "v3.1.1+incompatible",
)

go_repository(
    name = "com_github_go_ini_ini",
    importpath = "github.com/go-ini/ini",
    sum = "h1:Mujh4R/dH6YL8bxuISne3xX2+qcQ9p0IxKAP6ExWoUo=",
    version = "v1.25.4",
)

go_repository(
    name = "com_github_go_openapi_analysis",
    importpath = "github.com/go-openapi/analysis",
    sum = "h1:8b2ZgKfKIUTVQpTb77MoRDIMEIwvDVw40o3aOXdfYzI=",
    version = "v0.19.5",
)

go_repository(
    name = "com_github_go_openapi_errors",
    importpath = "github.com/go-openapi/errors",
    sum = "h1:a2kIyV3w+OS3S97zxUndRVD46+FhGOUBDFY7nmu4CsY=",
    version = "v0.19.2",
)

go_repository(
    name = "com_github_go_openapi_loads",
    importpath = "github.com/go-openapi/loads",
    sum = "h1:5I4CCSqoWzT+82bBkNIvmLc0UOsoKKQ4Fz+3VxOB7SY=",
    version = "v0.19.4",
)

go_repository(
    name = "com_github_go_openapi_runtime",
    importpath = "github.com/go-openapi/runtime",
    sum = "h1:csnOgcgAiuGoM/Po7PEpKDoNulCcF3FGbSnbHfxgjMI=",
    version = "v0.19.4",
)

go_repository(
    name = "com_github_go_openapi_strfmt",
    importpath = "github.com/go-openapi/strfmt",
    sum = "h1:eRfyY5SkaNJCAwmmMcADjY31ow9+N7MCLW7oRkbsINA=",
    version = "v0.19.3",
)

go_repository(
    name = "com_github_go_openapi_validate",
    importpath = "github.com/go-openapi/validate",
    sum = "h1:QhCBKRYqZR+SKo4gl1lPhPahope8/RLt6EVgY8X80w0=",
    version = "v0.19.5",
)

go_repository(
    name = "com_github_go_ozzo_ozzo_validation",
    importpath = "github.com/go-ozzo/ozzo-validation",
    sum = "h1:sUy/in/P6askYr16XJgTKq/0SZhiWsdg4WZGaLsGQkM=",
    version = "v3.5.0+incompatible",
)

go_repository(
    name = "com_github_godbus_dbus_v5",
    importpath = "github.com/godbus/dbus/v5",
    sum = "h1:ZqHaoEF7TBzh4jzPmqVhE/5A1z9of6orkAe5uHoAeME=",
    version = "v5.0.3",
)

go_repository(
    name = "com_github_golang_freetype",
    importpath = "github.com/golang/freetype",
    sum = "h1:DACJavvAHhabrF08vX0COfcOBJRhZ8lUbR+ZWIs0Y5g=",
    version = "v0.0.0-20170609003504-e2365dfdc4a0",
)

go_repository(
    name = "com_github_golangplus_bytes",
    importpath = "github.com/golangplus/bytes",
    sum = "h1:7xqw01UYS+KCI25bMrPxwNYkSns2Db1ziQPpVq99FpE=",
    version = "v0.0.0-20160111154220-45c989fe5450",
)

go_repository(
    name = "com_github_golangplus_fmt",
    importpath = "github.com/golangplus/fmt",
    sum = "h1:f5gsjBiF9tRRVomCvrkGMMWI8W1f2OBFar2c5oakAP0=",
    version = "v0.0.0-20150411045040-2a5d6d7d2995",
)

go_repository(
    name = "com_github_golangplus_testing",
    importpath = "github.com/golangplus/testing",
    sum = "h1:KhcknUwkWHKZPbFy2P7jH5LKJ3La+0ZeknkkmrSgqb0=",
    version = "v0.0.0-20180327235837-af21d9c3145e",
)

go_repository(
    name = "com_github_google_cadvisor",
    importpath = "github.com/google/cadvisor",
    sum = "h1:t3txV4zNZZGTuwuA/Onm3HToPhg16GjigAHZHEVIz+c=",
    version = "v0.37.0",
)

go_repository(
    name = "com_github_googlecloudplatform_k8s_cloud_provider",
    importpath = "github.com/GoogleCloudPlatform/k8s-cloud-provider",
    sum = "h1:JhyuWIqYrstW7KHMjk/fTqU0xtMpBOHuiTA2FVc7L4E=",
    version = "v0.0.0-20200415212048-7901bc822317",
)

go_repository(
    name = "com_github_gophercloud_gophercloud",
    importpath = "github.com/gophercloud/gophercloud",
    sum = "h1:P/nh25+rzXouhytV2pUHBb65fnds26Ghl8/391+sT5o=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_gopherjs_gopherjs",
    importpath = "github.com/gopherjs/gopherjs",
    sum = "h1:EGx4pi6eqNxGaHF6qqu48+N2wcFQ5qg5FXgOdqsJ5d8=",
    version = "v0.0.0-20181017120253-0766667cb4d1",
)

go_repository(
    name = "com_github_hashicorp_go_syslog",
    importpath = "github.com/hashicorp/go-syslog",
    sum = "h1:KaodqZuhUoZereWVIYmpUgZysurB1kBLX2j0MwMrUAE=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_heketi_heketi",
    importpath = "github.com/heketi/heketi",
    sum = "h1:ysqc8k973k1lLJ4BOOHAkx14K2nt4cLjsIm+hwWDZDE=",
    version = "v9.0.1-0.20190917153846-c2e2a4ab7ab9+incompatible",
)

go_repository(
    name = "com_github_heketi_tests",
    importpath = "github.com/heketi/tests",
    sum = "h1:oJ/NLadJn5HoxvonA6VxG31lg0d6XOURNA09BTtM4fY=",
    version = "v0.0.0-20151005000721-f3775cbcefd6",
)

go_repository(
    name = "com_github_ishidawataru_sctp",
    importpath = "github.com/ishidawataru/sctp",
    sum = "h1:rw3IAne6CDuVFlZbPOkA7bhxlqawFh7RJJ+CejfMaxE=",
    version = "v0.0.0-20191218070446-00ab2ac2db07",
)

go_repository(
    name = "com_github_jeffashton_win_pdh",
    importpath = "github.com/JeffAshton/win_pdh",
    sum = "h1:UKkYhof1njT1/xq4SEg5z+VpTgjmNeHwPGRQl7takDI=",
    version = "v0.0.0-20161109143554-76bb4ee9f0ab",
)

go_repository(
    name = "com_github_jimstudt_http_authentication",
    importpath = "github.com/jimstudt/http-authentication",
    sum = "h1:BcF8coBl0QFVhe8vAMMlD+CV8EISiu9MGKLoj6ZEyJA=",
    version = "v0.0.0-20140401203705-3eca13d6893a",
)

go_repository(
    name = "com_github_jtolds_gls",
    importpath = "github.com/jtolds/gls",
    sum = "h1:xdiiI2gbIgH/gLH7ADydsJ1uDOEzR8yvV7C0MuV77Wo=",
    version = "v4.20.0+incompatible",
)

go_repository(
    name = "com_github_jung_kurt_gofpdf",
    importpath = "github.com/jung-kurt/gofpdf",
    sum = "h1:PJr+ZMXIecYc1Ey2zucXdR73SMBtgjPgwa31099IMv0=",
    version = "v1.0.3-0.20190309125859-24315acbbda5",
)

go_repository(
    name = "com_github_karrick_godirwalk",
    importpath = "github.com/karrick/godirwalk",
    sum = "h1:VbzFqwXwNbAZoA6W5odrLr+hKK197CcENcPh6E/gJ0M=",
    version = "v1.7.5",
)

go_repository(
    name = "com_github_kylelemons_godebug",
    importpath = "github.com/kylelemons/godebug",
    sum = "h1:MtvEpTB6LX3vkb4ax0b5D2DHbNAUsen0Gx5wZoq3lV4=",
    version = "v0.0.0-20170820004349-d65d576e9348",
)

go_repository(
    name = "com_github_libopenstorage_openstorage",
    importpath = "github.com/libopenstorage/openstorage",
    sum = "h1:GLPam7/0mpdP8ZZtKjbfcXJBTIA/T1O6CBErVEFEyIM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_liggitt_tabwriter",
    importpath = "github.com/liggitt/tabwriter",
    sum = "h1:9TO3cAIGXtEhnIaL+V+BEER86oLrvS+kWobKpbJuye0=",
    version = "v0.0.0-20181228230101-89fcab3d43de",
)

go_repository(
    name = "com_github_lithammer_dedent",
    importpath = "github.com/lithammer/dedent",
    sum = "h1:VNzHMVCBNG1j0fh3OrsFRkVUwStdDArbgBWoPAffktY=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_lpabon_godbc",
    importpath = "github.com/lpabon/godbc",
    sum = "h1:ilqjArN1UOENJJdM34I2YHKmF/B0gGq4VLoSGy9iAao=",
    version = "v0.1.1",
)

go_repository(
    name = "com_github_lucas_clemente_aes12",
    importpath = "github.com/lucas-clemente/aes12",
    sum = "h1:sSeNEkJrs+0F9TUau0CgWTTNEwF23HST3Eq0A+QIx+A=",
    version = "v0.0.0-20171027163421-cd47fb39b79f",
)

go_repository(
    name = "com_github_lucas_clemente_quic_clients",
    importpath = "github.com/lucas-clemente/quic-clients",
    sum = "h1:/P9n0nICT/GnQJkZovtBqridjxU0ao34m7DpMts79qY=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_lucas_clemente_quic_go",
    importpath = "github.com/lucas-clemente/quic-go",
    sum = "h1:iQtTSZVbd44k94Lu0U16lLBIG3lrnjDvQongjPd4B/s=",
    version = "v0.10.2",
)

go_repository(
    name = "com_github_lucas_clemente_quic_go_certificates",
    importpath = "github.com/lucas-clemente/quic-go-certificates",
    sum = "h1:zqEC1GJZFbGZA0tRyNZqRjep92K5fujFtFsu5ZW7Aug=",
    version = "v0.0.0-20160823095156-d2f86524cced",
)

go_repository(
    name = "com_github_makenowjust_heredoc",
    importpath = "github.com/MakeNowJust/heredoc",
    sum = "h1:sjQovDkwrZp8u+gxLtPgKGjk5hCxuy2hrRejBTA9xFU=",
    version = "v0.0.0-20170808103936-bb23615498cd",
)

go_repository(
    name = "com_github_marten_seemann_qtls",
    importpath = "github.com/marten-seemann/qtls",
    sum = "h1:0yWJ43C62LsZt08vuQJDK1uC1czUc3FJeCLPoNAI4vA=",
    version = "v0.2.3",
)

go_repository(
    name = "com_github_mattn_go_runewidth",
    importpath = "github.com/mattn/go-runewidth",
    sum = "h1:2BvfKmzob6Bmd4YsL0zygOqfdFnK7GR4QL06Do4/p7Y=",
    version = "v0.0.4",
)

go_repository(
    name = "com_github_mholt_certmagic",
    importpath = "github.com/mholt/certmagic",
    sum = "h1:xKE9kZ5C8gelJC3+BNM6LJs1x21rivK7yxfTZMAuY2s=",
    version = "v0.6.2-0.20190624175158-6a42ef9fe8c2",
)

go_repository(
    name = "com_github_microsoft_hcsshim",
    importpath = "github.com/Microsoft/hcsshim",
    sum = "h1:2T9t72RkTRjAcuFc+4vaGWnRx/anVngE1/VGN/HFEVk=",
    version = "v0.8.8-0.20200312192636-fd0797d766b1",
)

go_repository(
    name = "com_github_miekg_dns",
    importpath = "github.com/miekg/dns",
    sum = "h1:sJFOl9BgwbYAWOGEwr61FU28pqsBNdpRBnhGXtO06Oo=",
    version = "v1.1.31",
)

go_repository(
    name = "com_github_mindprince_gonvml",
    importpath = "github.com/mindprince/gonvml",
    sum = "h1:PS1dLCGtD8bb9RPKJrc8bS7qHL6JnW1CZvwzH9dPoUs=",
    version = "v0.0.0-20190828220739-9ebdce4bb989",
)

go_repository(
    name = "com_github_mistifyio_go_zfs",
    importpath = "github.com/mistifyio/go-zfs",
    sum = "h1:aKW/4cBs+yK6gpqU3K/oIwk9Q/XICqd3zOX/UFuvqmk=",
    version = "v2.1.2-0.20190413222219-f784269be439+incompatible",
)

go_repository(
    name = "com_github_mitchellh_go_wordwrap",
    importpath = "github.com/mitchellh/go-wordwrap",
    sum = "h1:6GlHJ/LTGMrIJbwgdqdl2eEH8o+Exx/0m8ir9Gns0u4=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_moby_ipvs",
    importpath = "github.com/moby/ipvs",
    sum = "h1:aoZ7fhLTXgDbzVrAnvV+XbKOU8kOET7B3+xULDF/1o0=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_moby_sys_mountinfo",
    importpath = "github.com/moby/sys/mountinfo",
    sum = "h1:KIrhRO14+AkwKvG/g2yIpNMOUVZ02xNhOw8KY1WsLOI=",
    version = "v0.1.3",
)

go_repository(
    name = "com_github_moby_term",
    importpath = "github.com/moby/term",
    sum = "h1:aY7OQNf2XqY/JQ6qREWamhI/81os/agb2BAGpcx5yWI=",
    version = "v0.0.0-20200312100748-672ec06f55cd",
)

go_repository(
    name = "com_github_mohae_deepcopy",
    importpath = "github.com/mohae/deepcopy",
    sum = "h1:e+l77LJOEqXTIQihQJVkA6ZxPOUmfPM5e4H7rcpgtSk=",
    version = "v0.0.0-20170603005431-491d3605edfb",
)

go_repository(
    name = "com_github_mrunalp_fileutils",
    importpath = "github.com/mrunalp/fileutils",
    sum = "h1:aZQToFSLH8ejFeSkTc3r3L4dPImcj7Ib/KgmkQqbGGg=",
    version = "v0.0.0-20200520151820-abd8a0e76976",
)

go_repository(
    name = "com_github_mvdan_xurls",
    importpath = "github.com/mvdan/xurls",
    sum = "h1:OpuDelGQ1R1ueQ6sSryzi6P+1RtBpfQHM8fJwlE45ww=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_naoina_go_stringutil",
    importpath = "github.com/naoina/go-stringutil",
    sum = "h1:rCUeRUHjBjGTSHl0VC00jUPLz8/F9dDzYI70Hzifhks=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_naoina_toml",
    importpath = "github.com/naoina/toml",
    sum = "h1:PT/lllxVVN0gzzSqSlHEmP8MJB4MY2U7STGxiouV4X8=",
    version = "v0.1.1",
)

go_repository(
    name = "com_github_olekukonko_tablewriter",
    importpath = "github.com/olekukonko/tablewriter",
    sum = "h1:58+kh9C6jJVXYjt8IE48G2eWl6BjwU5Gj0gqY84fy78=",
    version = "v0.0.0-20170122224234-a0225b3f23b5",
)

go_repository(
    name = "com_github_opencontainers_runc",
    importpath = "github.com/opencontainers/runc",
    sum = "h1:+IczUKCRzDzFDnw99O/PAqrcBBCoRp9xN3cB1SYSNS4=",
    version = "v1.0.0-rc92",
)

go_repository(
    name = "com_github_opencontainers_runtime_spec",
    importpath = "github.com/opencontainers/runtime-spec",
    sum = "h1:NhsM2gc769rVWDqJvapK37r+7+CBXI8xHhnfnt8uQsg=",
    version = "v1.0.3-0.20200728170252-4d89ac9fbff6",
)

go_repository(
    name = "com_github_opencontainers_selinux",
    importpath = "github.com/opencontainers/selinux",
    sum = "h1:+bIAS/Za3q5FTwWym4fTB0vObnfCf3G/NC7K6Jx62mY=",
    version = "v1.6.0",
)

go_repository(
    name = "com_github_pborman_uuid",
    importpath = "github.com/pborman/uuid",
    sum = "h1:J7Q5mO4ysT1dv8hyrUGHb9+ooztCXu1D8MY8DZYsu3g=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_pquerna_cachecontrol",
    importpath = "github.com/pquerna/cachecontrol",
    sum = "h1:J9b7z+QKAmPf4YLrFg6oQUotqHQeUNWwkvo7jZp1GLU=",
    version = "v0.0.0-20180517163645-1555304b9b35",
)

go_repository(
    name = "com_github_quobyte_api",
    importpath = "github.com/quobyte/api",
    sum = "h1:lPHLsuvtjFyk8WhC4uHoHRkScijIHcffTWBBP+YpzYo=",
    version = "v0.1.2",
)

go_repository(
    name = "com_github_robfig_cron",
    importpath = "github.com/robfig/cron",
    sum = "h1:jk4/Hud3TTdcrJgUOBgsqrZBarcxl6ADIjSC2iniwLY=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_rubiojr_go_vhd",
    importpath = "github.com/rubiojr/go-vhd",
    sum = "h1:if3/24+h9Sq6eDx8UUz1SO9cT9tizyIsATfB7b4D3tc=",
    version = "v0.0.0-20200706105327-02e210299021",
)

go_repository(
    name = "com_github_russross_blackfriday_v2",
    importpath = "github.com/russross/blackfriday/v2",
    sum = "h1:lPqVAte+HuHNfhJ/0LC98ESWRz8afy9tM/0RK8m9o+Q=",
    version = "v2.0.1",
)

go_repository(
    name = "com_github_seccomp_libseccomp_golang",
    importpath = "github.com/seccomp/libseccomp-golang",
    sum = "h1:vWI7jlfJVWN//T8jrt5JduYkSid+Sl/fRz33J1jQ83k=",
    version = "v0.9.2-0.20200314001724-bdab42bd5128",
)

go_repository(
    name = "com_github_sergi_go_diff",
    importpath = "github.com/sergi/go-diff",
    sum = "h1:Kpca3qRNrduNnOQeazBd0ysaKrUJiIuISHxogkT9RPQ=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_shurcool_sanitized_anchor_name",
    importpath = "github.com/shurcooL/sanitized_anchor_name",
    sum = "h1:PdmoCO6wvbs+7yrJyMORt4/BmY5IYyJwS/kOiWx8mHo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_smartystreets_assertions",
    importpath = "github.com/smartystreets/assertions",
    sum = "h1:zE9ykElWQ6/NYmHa3jpm/yHnI4xSofP+UP6SpjHcSeM=",
    version = "v0.0.0-20180927180507-b2de0cb4f26d",
)

go_repository(
    name = "com_github_smartystreets_goconvey",
    importpath = "github.com/smartystreets/goconvey",
    sum = "h1:fv0U8FUIMPNf1L9lnHLvLhgicrIVChEkdzIKYqbNC9s=",
    version = "v1.6.4",
)

go_repository(
    name = "com_github_storageos_go_api",
    importpath = "github.com/storageos/go-api",
    sum = "h1:n+WYaU0kQ6WIiuEyWSgbXqkBx16irO69kYCtwVYoO5s=",
    version = "v0.0.0-20180912212459-343b3eff91fc",
)

go_repository(
    name = "com_github_syndtr_gocapability",
    importpath = "github.com/syndtr/gocapability",
    sum = "h1:b6uOv7YOFK0TYG7HtkIgExQo+2RdLuwRft63jn2HWj8=",
    version = "v0.0.0-20180916011248-d98352740cb2",
)

go_repository(
    name = "com_github_thecodeteam_goscaleio",
    importpath = "github.com/thecodeteam/goscaleio",
    sum = "h1:SB5tO98lawC+UK8ds/U2jyfOCH7GTcFztcF5x9gbut4=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_urfave_cli",
    importpath = "github.com/urfave/cli",
    sum = "h1:+mkCCcOFKPnCmVYVcURKps1Xe+3zP90gSYGNfRkjoIY=",
    version = "v1.22.1",
)

go_repository(
    name = "com_github_urfave_negroni",
    importpath = "github.com/urfave/negroni",
    sum = "h1:kIimOitoypq34K7TG7DUaJ9kq/N4Ofuwi1sjz0KipXc=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_vektah_gqlparser",
    importpath = "github.com/vektah/gqlparser",
    sum = "h1:ZsyLGn7/7jDNI+y4SEhI4yAxRChlv15pUHMjijT+e68=",
    version = "v1.1.2",
)

go_repository(
    name = "com_github_vishvananda_netlink",
    importpath = "github.com/vishvananda/netlink",
    sum = "h1:1iyaYNBLmP6L0220aDnYQpo1QEV4t4hJ+xEEhhJH8j0=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_vishvananda_netns",
    importpath = "github.com/vishvananda/netns",
    sum = "h1:OviZH7qLw/7ZovXvuNyL3XQl8UFofeikI1NW1Gypu7k=",
    version = "v0.0.0-20191106174202-0a2b9b5464df",
)

go_repository(
    name = "com_github_vmware_govmomi",
    importpath = "github.com/vmware/govmomi",
    sum = "h1:f7QxSmP7meCtoAmiKZogvVbLInT+CZx6Px6K5rYsJZo=",
    version = "v0.18.0",
)

go_repository(
    name = "com_github_xlab_handysort",
    importpath = "github.com/xlab/handysort",
    sum = "h1:j2hhcujLRHAg872RWAV5yaUrEjHEObwDv3aImCaNLek=",
    version = "v0.0.0-20150421192137-fb3537ed64a1",
)

go_repository(
    name = "in_gopkg_cheggaaa_pb_v1",
    importpath = "gopkg.in/cheggaaa/pb.v1",
    sum = "h1:kJdccidYzt3CaHD1crCFTS1hxyhSi059NhOFUf03YFo=",
    version = "v1.0.27",
)

go_repository(
    name = "in_gopkg_gcfg_v1",
    importpath = "gopkg.in/gcfg.v1",
    sum = "h1:m8OOJ4ccYHnx2f4gQwpno8nAX5OGOh7RLaaz0pj3Ogs=",
    version = "v1.2.3",
)

go_repository(
    name = "in_gopkg_mcuadros_go_syslog_v2",
    importpath = "gopkg.in/mcuadros/go-syslog.v2",
    sum = "h1:60g8zx1BijSVSgLTzLCW9UC4/+i1Ih9jJ1DR5Tgp9vE=",
    version = "v2.2.1",
)

go_repository(
    name = "in_gopkg_natefinch_lumberjack_v2",
    importpath = "gopkg.in/natefinch/lumberjack.v2",
    sum = "h1:1Lc07Kr7qY4U2YPouBjpCLxpiyxIVoxqXgkXLknAOE8=",
    version = "v2.0.0",
)

go_repository(
    name = "in_gopkg_square_go_jose_v2",
    importpath = "gopkg.in/square/go-jose.v2",
    sum = "h1:H0TmLt7/KmzlrDOpa1F+zr0Tk90PbJYBfsVUmRLrf9Y=",
    version = "v2.4.1",
)

go_repository(
    name = "in_gopkg_warnings_v0",
    importpath = "gopkg.in/warnings.v0",
    sum = "h1:wFXVbFY8DY5/xOe1ECiWdKCzZlxgshcYVNkBHstARME=",
    version = "v0.1.2",
)

go_repository(
    name = "io_etcd_go_etcd",
    importpath = "go.etcd.io/etcd",
    sum = "h1:VcrIfasaLFkyjk6KNlXQSzO+B0fZcnECiDrKJsfxka0=",
    version = "v0.0.0-20191023171146-3cf2f69b5738",
)

go_repository(
    name = "io_rsc_pdf",
    importpath = "rsc.io/pdf",
    sum = "h1:k1MczvYDUvJBe93bYd7wrZLLUEcLZAuF824/I4e5Xr4=",
    version = "v0.1.1",
)

go_repository(
    name = "ml_vbom_util",
    importpath = "vbom.ml/util",
    sum = "h1:MksmcCZQWAQJCTA5T0jgI/0sJ51AVm4Z41MrmfczEoc=",
    version = "v0.0.0-20160121211510-db5cfe13f5cc",
)

go_repository(
    name = "org_bitbucket_bertimus9_systemstat",
    importpath = "bitbucket.org/bertimus9/systemstat",
    sum = "h1:N9r8OBSXAgEUfho3SQtZLY8zo6E1OdOMvelvP22aVFc=",
    version = "v0.0.0-20180207000608-0eeff89b0690",
)

go_repository(
    name = "org_gonum_v1_gonum",
    importpath = "gonum.org/v1/gonum",
    sum = "h1:4r+yNT0+8SWcOkXP+63H2zQbN+USnC73cjGUxnDF94Q=",
    version = "v0.6.2",
)

go_repository(
    name = "org_gonum_v1_netlib",
    importpath = "gonum.org/v1/netlib",
    sum = "h1:jRyg0XfpwWlhEV8mDfdNGBeSJM2fuyh9Yjrnd8kF2Ts=",
    version = "v0.0.0-20190331212654-76723241ea4e",
)

go_repository(
    name = "org_gonum_v1_plot",
    importpath = "gonum.org/v1/plot",
    sum = "h1:Qh4dB5D/WpoUUp3lSod7qgoyEHbDGPUWjIbnqdqqe1k=",
    version = "v0.0.0-20190515093506-e2840ee46a6b",
)

go_repository(
    name = "org_modernc_cc",
    importpath = "modernc.org/cc",
    sum = "h1:nPibNuDEx6tvYrUAtvDTTw98rx5juGsa5zuDnKwEEQQ=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_xc",
    importpath = "modernc.org/xc",
    sum = "h1:7ccXrupWZIS3twbUGrtKmHS2DXY6xegFua+6O3xgAFU=",
    version = "v1.0.0",
)

go_repository(
    name = "tools_gotest",
    importpath = "gotest.tools",
    sum = "h1:VsBPFP1AI068pPrMxtb/S8Zkgf9xEmTLJjfM+P5UIEo=",
    version = "v2.2.0+incompatible",
)

go_repository(
    name = "com_github_etherlabsio_healthcheck",
    importpath = "github.com/etherlabsio/healthcheck",
    sum = "h1:az9jaEKre+mwUWiS9Pl8h1FuOvdiFM7UqplmCmJtHUQ=",
    version = "v0.0.0-20191224061800-dd3d2fd8c3f6",
)

go_repository(
    name = "com_github_zserge_hid",
    importpath = "github.com/zserge/hid",
    sum = "h1:DAdExZJtAXR150WHt4wzTCzHH7lXkD4vKN9xLtF05VU=",
    version = "v0.0.0-20190124175232-e1626f1782f3",
)

go_repository(
    name = "com_github_dsnet_compress",
    importpath = "github.com/dsnet/compress",
    sum = "h1:PlZu0n3Tuv04TzpfPbrnI0HW/YwodEXDS+oPKahKF0Q=",
    version = "v0.0.1",
)

go_repository(
    name = "com_github_dsnet_golib",
    importpath = "github.com/dsnet/golib",
    sum = "h1:tFh1tRc4CA31yP6qDcu+Trax5wW5GuMxvkIba07qVLY=",
    version = "v0.0.0-20171103203638-1ea166775780",
)

go_repository(
    name = "com_github_magefile_mage",
    importpath = "github.com/magefile/mage",
    sum = "h1:3HiXzCUY12kh9bIuyXShaVe529fJfyqoVM42o/uom2g=",
    version = "v1.10.0",
)

go_repository(
    name = "com_github_mholt_archiver",
    importpath = "github.com/mholt/archiver",
    sum = "h1:1dCVxuqs0dJseYEhi5pl7MYPH9zDa1wBi7mF09cbNkU=",
    version = "v3.1.1+incompatible",
)

go_repository(
    name = "com_github_nwaples_rardecode",
    importpath = "github.com/nwaples/rardecode",
    sum = "h1:vSxaY8vQhOcVr4mm5e8XllHWTiM4JF507A0Katqw7MQ=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_ulikunitz_xz",
    importpath = "github.com/ulikunitz/xz",
    sum = "h1:ERv8V6GKqVi23rgu5cj9pVfVzJbOqAY2Ntl88O6c2nQ=",
    version = "v0.5.8",
)

go_repository(
    name = "com_github_wrouesnel_postgres_exporter",
    importpath = "github.com/wrouesnel/postgres_exporter",
    sum = "h1:DjiSYqYpq3dtuYkP18WhhhVtxaKJZk2+vJoOPuowKtU=",
    version = "v0.8.0",
)

go_repository(
    name = "com_github_xi2_xz",
    importpath = "github.com/xi2/xz",
    sum = "h1:nIPpBwaJSVYIxUFsDv3M8ofmx9yWTog9BfvIu0q41lo=",
    version = "v0.0.0-20171230120015-48954b6210f8",
)

go_repository(
    name = "com_github_mdlayher_unifi",
    importpath = "github.com/mdlayher/unifi",
    patch_args = ["-p1"],
    patches = [
        "//tools:unifi.patch",
    ],
    sum = "h1:Ojrc2OzOsMpz8m8sTLjptIJ2dg4ITH+4ZA1cPWykkMg=",
    version = "v0.0.0-20180604180305-be73af8d7922",
)

go_repository(
    name = "com_github_mdlayher_unifi_exporter",
    importpath = "github.com/mdlayher/unifi_exporter",
    patch_args = ["-p1"],
    # This is needed because the patch above adds the promhttp dependency, but patches in go_repository
    # rules run after Gazelle has already run for the repo, so it ends up missing from the build file.
    patch_cmds = ["""sed -i.bak -e '/com_github_prometheus_client_golang/ a\\
        "@com_github_prometheus_client_golang//prometheus/promhttp",
' -e s/private/public/ cmd/unifi_exporter/BUILD.bazel
"""],
    patches = [
        "//tools:unifi_exporter.patch",
    ],
    sum = "h1:Fva/ZqxceEhabBhEdSnxVZC/17X+pKw917rrV8DeVGM=",
    version = "v0.0.0-20180610203931-85455df7c491",
)

go_repository(
    name = "com_github_google_go_github_v33",
    importpath = "github.com/google/go-github/v33",
    sum = "h1:qAf9yP0qc54ufQxzwv+u9H0tiVOnPJxo0lI/JXqw3ZM=",
    version = "v33.0.0",
)

go_repository(
    name = "com_github_hako_durafmt",
    importpath = "github.com/hako/durafmt",
    sum = "h1:BpJ2o0OR5FV7vrkDYfXYVJQeMNWa8RhklZOpW2ITAIQ=",
    version = "v0.0.0-20200710122514-c0fb7b4da026",
)

go_repository(
    name = "com_github_mjm_graphql_go",
    importpath = "github.com/mjm/graphql-go",
    sum = "h1:uFZvCeM6RutHugCbyVQc5GHyTOVYn7KMb4OOR5niMQg=",
    version = "v0.0.0-20200213070811-48021cabfccf",
)

go_repository(
    name = "com_github_abdullin_seq",
    importpath = "github.com/abdullin/seq",
    sum = "h1:DBNMBMuMiWYu0b+8KMJuWmfCkcxl09JwdlqwDZZ6U14=",
    version = "v0.0.0-20160510034733-d5467c17e7af",
)

go_repository(
    name = "com_github_armon_go_metrics",
    importpath = "github.com/armon/go-metrics",
    sum = "h1:23nUvGE+8HYFc0AUXuYxgFws6IdyzOrSJJmKfPMJmi8=",
    version = "v0.3.5-0.20201104215618-6fd5a4ddf425",
)

go_repository(
    name = "com_github_armon_go_radix",
    importpath = "github.com/armon/go-radix",
    sum = "h1:F4z6KzEeeQIMeLFa97iZU6vupzoecKdU5TX24SNppXI=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest_azure_auth",
    importpath = "github.com/Azure/go-autorest/autorest/azure/auth",
    sum = "h1:bvUhZciHydpBxBmCheUgxxbSwJy7xcfjkUsjUcqSojc=",
    version = "v0.5.1",
)

go_repository(
    name = "com_github_azure_go_autorest_autorest_azure_cli",
    importpath = "github.com/Azure/go-autorest/autorest/azure/cli",
    sum = "h1:Ml+UCrnlKD+cJmSzrZ/RDcDw86NjkRUpnFh7V5JUhzU=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_circonus_labs_circonus_gometrics",
    importpath = "github.com/circonus-labs/circonus-gometrics",
    sum = "h1:C29Ae4G5GtYyYMm1aztcyj/J5ckgJm2zwdDajFbx1NY=",
    version = "v2.3.1+incompatible",
)

go_repository(
    name = "com_github_circonus_labs_circonusllhist",
    importpath = "github.com/circonus-labs/circonusllhist",
    sum = "h1:TJH+oke8D16535+jHExHj4nQvzlZrj7ug5D7I/orNUA=",
    version = "v0.1.3",
)

go_repository(
    name = "com_github_coredns_coredns",
    importpath = "github.com/coredns/coredns",
    sum = "h1:bAFHrSsBeTeRG5W3Nf2su3lUGw7Npw2UKeCJm/3A638=",
    version = "v1.1.2",
)

go_repository(
    name = "com_github_datadog_datadog_go",
    importpath = "github.com/DataDog/datadog-go",
    sum = "h1:qSG2N4FghB1He/r2mFrWKCaL7dXCilEuNEeAn20fdD4=",
    version = "v3.2.0+incompatible",
)

go_repository(
    name = "com_github_denverdino_aliyungo",
    importpath = "github.com/denverdino/aliyungo",
    sum = "h1:p6poVbjHDkKa+wtC8frBMwQtT3BmqGYBjzMwJ63tuR4=",
    version = "v0.0.0-20190125010748-a747050bb1ba",
)

go_repository(
    name = "com_github_digitalocean_godo",
    importpath = "github.com/digitalocean/godo",
    sum = "h1:uW1/FcvZE/hoixnJcnlmIUvTVNdZCLjRLzmDtRi1xXY=",
    version = "v1.10.0",
)

go_repository(
    name = "com_github_dimchansky_utfbom",
    importpath = "github.com/dimchansky/utfbom",
    sum = "h1:FcM3g+nofKgUteL8dm/UpdRXNC9KmADgTpLKsu0TRo4=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_elazarl_go_bindata_assetfs",
    importpath = "github.com/elazarl/go-bindata-assetfs",
    sum = "h1:AwZUiMWfYSmIiHdFJIubTSs8BFIFoMmUFbeuwBzHIPs=",
    version = "v1.0.1-0.20200509193318-234c15e7648f",
)

go_repository(
    name = "com_github_fatih_structs",
    importpath = "github.com/fatih/structs",
    sum = "h1:Q7juDM0QtcnhCpeyLGQKyg4TOIghuNXrkL32pHAUMxo=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_frankban_quicktest",
    importpath = "github.com/frankban/quicktest",
    sum = "h1:rCSCih1FnSWJEel/eub9wclBSqpF2F/PuvxUWGWnbO8=",
    version = "v1.4.0",
)

go_repository(
    name = "com_github_go_asn1_ber_asn1_ber",
    importpath = "github.com/go-asn1-ber/asn1-ber",
    sum = "h1:gvPdv/Hr++TRFCl0UbPFHC54P9N9jgsRPnmnr419Uck=",
    version = "v1.3.1",
)

go_repository(
    name = "com_github_go_ldap_ldap_v3",
    importpath = "github.com/go-ldap/ldap/v3",
    sum = "h1:RIgdpHXJpsUqUK5WXwKyVsESrGFqo5BRWPk3RR4/ogQ=",
    version = "v3.1.3",
)

go_repository(
    name = "com_github_go_test_deep",
    importpath = "github.com/go-test/deep",
    sum = "h1:ZrJSEWsXzPOxaZnFteGEfooLba+ju3FYIbOrS+rQd68=",
    version = "v1.0.3",
)

go_repository(
    name = "com_github_google_tcpproxy",
    importpath = "github.com/google/tcpproxy",
    sum = "h1:AtvtonGEH/fZK0XPNNBdB6swgy7Iudfx88wzyIpwqJ8=",
    version = "v0.0.0-20180808230851-dfa16c61dad2",
)

go_repository(
    name = "com_github_hashicorp_consul",
    importpath = "github.com/hashicorp/consul",
    sum = "h1:hp308KxAf3zWoGuwp2e+0UUhrm6qHjeBQk3jCZ+bjcY=",
    version = "v1.7.8",
)

go_repository(
    name = "com_github_hashicorp_consul_api",
    importpath = "github.com/hashicorp/consul/api",
    sum = "h1:BOEQaMWoGMhmQ29fC26bi0qb7/rId9JzZP2V0Xmx7m8=",
    version = "v1.8.1",
)

go_repository(
    name = "com_github_hashicorp_consul_sdk",
    importpath = "github.com/hashicorp/consul/sdk",
    sum = "h1:H6R9d008jDcHPQPAqPNuydAshJ4v5/8URdFnUvK/+sc=",
    version = "v0.7.0",
)

go_repository(
    name = "com_github_hashicorp_go_bexpr",
    importpath = "github.com/hashicorp/go-bexpr",
    sum = "h1:ijMXI4qERbzxbCnkxmfUtwMyjrrk3y+Vt0MxojNCbBs=",
    version = "v0.1.2",
)

go_repository(
    name = "com_github_hashicorp_go_checkpoint",
    importpath = "github.com/hashicorp/go-checkpoint",
    sum = "h1:XDCSythtg8aWSRSO29uwhgh7b127fWr+m5SemqjSUL8=",
    version = "v0.0.0-20171009173528-1545e56e46de",
)

go_repository(
    name = "com_github_hashicorp_go_cleanhttp",
    importpath = "github.com/hashicorp/go-cleanhttp",
    sum = "h1:dH3aiDG9Jvb5r5+bYHsikaOUIpcM0xvgMXVoDkXMzJM=",
    version = "v0.5.1",
)

go_repository(
    name = "com_github_hashicorp_go_connlimit",
    importpath = "github.com/hashicorp/go-connlimit",
    sum = "h1:oAojHGjFxUTTTA8c5XXnDqWJ2HLuWbDiBPTpWvNzvqM=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_hashicorp_go_discover",
    importpath = "github.com/hashicorp/go-discover",
    sum = "h1:7WFMVeuJQp6BkzuTv9O52pzwtEFVUJubKYN+zez8eTI=",
    version = "v0.0.0-20200812215701-c4b85f6ed31f",
)

go_repository(
    name = "com_github_hashicorp_go_hclog",
    importpath = "github.com/hashicorp/go-hclog",
    sum = "h1:nQcJDQwIAGnmoUWp8ubocEX40cCml/17YkF6csQLReU=",
    version = "v0.14.1",
)

go_repository(
    name = "com_github_hashicorp_go_immutable_radix",
    importpath = "github.com/hashicorp/go-immutable-radix",
    sum = "h1:8exGP7ego3OmkfksihtSouGMZ+hQrhxx+FVELeXpVPE=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_hashicorp_go_kms_wrapping_entropy",
    importpath = "github.com/hashicorp/go-kms-wrapping/entropy",
    sum = "h1:xuTi5ZwjimfpvpL09jDE71smCBRpnF5xfo871BSX4gs=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_hashicorp_go_memdb",
    importpath = "github.com/hashicorp/go-memdb",
    sum = "h1:xdXq34gBOMEloa9rlGStLxmfX/dyIK8htOv36dQUwHU=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_hashicorp_go_msgpack",
    importpath = "github.com/hashicorp/go-msgpack",
    sum = "h1:9byZdVjKTe5mce63pRVNP1L7UAmdHOTEMGehn6KvJWs=",
    version = "v1.1.5",
)

go_repository(
    name = "com_github_hashicorp_go_plugin",
    importpath = "github.com/hashicorp/go-plugin",
    sum = "h1:XzRWU4VSJBqGVxl6tWD+6ITHftMhvRDKKyp1WVSXAhY=",
    version = "v1.0.2-0.20191004171845-809113480b55",
)

go_repository(
    name = "com_github_hashicorp_go_raftchunking",
    importpath = "github.com/hashicorp/go-raftchunking",
    sum = "h1:moEnaG3gcwsWNyIBJoD5PCByE+Ewkqxh6N05CT+MbwA=",
    version = "v0.6.1",
)

go_repository(
    name = "com_github_hashicorp_go_retryablehttp",
    importpath = "github.com/hashicorp/go-retryablehttp",
    sum = "h1:1BZvpawXoJCWX6pNtow9+rpEj+3itIlutiqnntI6jOE=",
    version = "v0.5.4",
)

go_repository(
    name = "com_github_hashicorp_go_rootcerts",
    importpath = "github.com/hashicorp/go-rootcerts",
    sum = "h1:jzhAVGtqPKbwpyCPELlgNWhE1znq+qwJtW5Oi2viEzc=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_hashicorp_go_sockaddr",
    importpath = "github.com/hashicorp/go-sockaddr",
    sum = "h1:ztczhD1jLxIRjVejw8gFomI1BQZOe2WoVOu0SyteCQc=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_hashicorp_go_uuid",
    importpath = "github.com/hashicorp/go-uuid",
    sum = "h1:cfejS+Tpcp13yd5nYHWDI6qVCny6wyX2Mt5SGur2IGE=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_hashicorp_go_version",
    importpath = "github.com/hashicorp/go-version",
    sum = "h1:8gXJgyR86/nGv7IqrnmnKJVNxC9APdvAVSxc9mGxfKk=",
    version = "v1.2.1-0.20191009193637-2046c9d0f0b0",
)

go_repository(
    name = "com_github_hashicorp_hil",
    importpath = "github.com/hashicorp/hil",
    sum = "h1:uk280DXEbQiCOZgCOI3elFSeNxf8YIZiNsbr2pQLYD0=",
    version = "v0.0.0-20160711231837-1e86c6b523c5",
)

go_repository(
    name = "com_github_hashicorp_logutils",
    importpath = "github.com/hashicorp/logutils",
    sum = "h1:dLEQVugN8vlakKOUE3ihGLTZJRB4j+M2cdTm/ORI65Y=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_hashicorp_mdns",
    importpath = "github.com/hashicorp/mdns",
    sum = "h1:XFSOubp8KWB+Jd2PDyaX5xUd5bhSP/+pTDZVDMzZJM8=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_hashicorp_memberlist",
    importpath = "github.com/hashicorp/memberlist",
    sum = "h1:5+RffWKwqJ71YPu9mWsF7ZOscZmwfasdA8kbdC7AO2g=",
    version = "v0.2.2",
)

go_repository(
    name = "com_github_hashicorp_net_rpc_msgpackrpc",
    importpath = "github.com/hashicorp/net-rpc-msgpackrpc",
    sum = "h1:lc3c72qGlIMDqQpQH82Y4vaglRMMFdJbziYWriR4UcE=",
    version = "v0.0.0-20151116020338-a14192a58a69",
)

go_repository(
    name = "com_github_hashicorp_raft",
    importpath = "github.com/hashicorp/raft",
    sum = "h1:+H0NF++gFCFqQQdjRstZssPZcx9BiaMOGHxSxoNvBN8=",
    version = "v1.1.3-0.20200211192230-365023de17e6",
)

go_repository(
    name = "com_github_hashicorp_raft_autopilot",
    importpath = "github.com/hashicorp/raft-autopilot",
    sum = "h1:f8Dv2y1Vq8ttuH2+oh5l87Paj/BINpMm5TBrMLx+qGQ=",
    version = "v0.1.1",
)

go_repository(
    name = "com_github_hashicorp_raft_boltdb",
    importpath = "github.com/hashicorp/raft-boltdb",
    sum = "h1:xykPFhrBAS2J0VBzVa5e80b5ZtYuNQtgXjN40qBZlD4=",
    version = "v0.0.0-20171010151810-6e5ba93211ea",
)

go_repository(
    name = "com_github_hashicorp_serf",
    importpath = "github.com/hashicorp/serf",
    sum = "h1:EBWvyu9tcRszt3Bxp3KNssBMP1KuHWyO51lz9+786iM=",
    version = "v0.9.5",
)

go_repository(
    name = "com_github_hashicorp_vault_api",
    importpath = "github.com/hashicorp/vault/api",
    sum = "h1:2qdbEUXjHohC+OYHtVU5lujvPAHPKYR4IMs9rsiUTk8=",
    version = "v1.0.5-0.20190730042357-746c0b111519",
)

go_repository(
    name = "com_github_hashicorp_vault_sdk",
    build_file_proto_mode = "disable",
    importpath = "github.com/hashicorp/vault/sdk",
    sum = "h1:fLUoZ8cI/pqlVCk09r88cVoY7ggKEl1A4e6Mujr3RvU=",
    version = "v0.1.14-0.20190730042320-0dc007d98cc8",
)

go_repository(
    name = "com_github_hashicorp_vic",
    importpath = "github.com/hashicorp/vic",
    sum = "h1:O/pT5C1Q3mVXMyuqg7yuAWUg/jMZR1/0QTzTRdNR6Uw=",
    version = "v1.5.1-0.20190403131502-bbfe86ec9443",
)

go_repository(
    name = "com_github_hashicorp_yamux",
    importpath = "github.com/hashicorp/yamux",
    sum = "h1:kJCB4vdITiW1eC1vq2e6IsrXKrZit1bv/TDYFGMp4BQ=",
    version = "v0.0.0-20181012175058-2f1d1f20f75d",
)

go_repository(
    name = "com_github_jackc_fake",
    importpath = "github.com/jackc/fake",
    sum = "h1:vr3AYkKovP8uR8AvSGGUK1IDqRa5lAAvEkZG1LKaCRc=",
    version = "v0.0.0-20150926172116-812a484cc733",
)

go_repository(
    name = "com_github_jackc_pgx",
    importpath = "github.com/jackc/pgx",
    sum = "h1:Wa90/+qsITBAPkAZjiByeIGHFcj3Ztu+VzrrIpHjL90=",
    version = "v3.3.0+incompatible",
)

go_repository(
    name = "com_github_jarcoal_httpmock",
    importpath = "github.com/jarcoal/httpmock",
    sum = "h1:FjHUJJ7oBW4G/9j1KzlHaXL09LyMVM9rupS39lncbXk=",
    version = "v0.0.0-20180424175123-9c70cfe4a1da",
)

go_repository(
    name = "com_github_joyent_triton_go",
    importpath = "github.com/joyent/triton-go",
    sum = "h1:BvV6PYcRz0yGnWXNZrd5wginNT1GfFfPvvWpPbjfFL8=",
    version = "v0.0.0-20190112182421-51ffac552869",
)

go_repository(
    name = "com_github_linode_linodego",
    importpath = "github.com/linode/linodego",
    sum = "h1:4WZmMpSA2NRwlPZcc0+4Gyn7rr99Evk9bnr0B3gXRKE=",
    version = "v0.7.1",
)

go_repository(
    name = "com_github_mitchellh_cli",
    importpath = "github.com/mitchellh/cli",
    sum = "h1:tEElEatulEHDeedTxwckzyYMA5c86fbmNIUL1hBIiTg=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_mitchellh_copystructure",
    importpath = "github.com/mitchellh/copystructure",
    sum = "h1:Laisrj+bAB6b/yJwB5Bt3ITZhGJdqmxquMKeZ+mmkFQ=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_mitchellh_go_testing_interface",
    importpath = "github.com/mitchellh/go-testing-interface",
    sum = "h1:/x0XQ6h+3U3nAyk1yx+bHPURrKa9sVVvYbuqZ7pIAtI=",
    version = "v1.14.0",
)

go_repository(
    name = "com_github_mitchellh_hashstructure",
    importpath = "github.com/mitchellh/hashstructure",
    sum = "h1:ZkRJX1CyOoTkar7p/mLS5TZU4nJ1Rn/F8u9dGS02Q3Y=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_mitchellh_pointerstructure",
    importpath = "github.com/mitchellh/pointerstructure",
    sum = "h1:ATSdz4NWrmWPOF1CeCBU4sMCno2hgqdbSrRPFWQSVZI=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_mitchellh_reflectwalk",
    importpath = "github.com/mitchellh/reflectwalk",
    sum = "h1:FVzMWA5RllMAKIdUSC8mdWo3XtwoecrH79BY70sEEpE=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_nicolai86_scaleway_sdk",
    importpath = "github.com/nicolai86/scaleway-sdk",
    sum = "h1:BQ1HW7hr4IVovMwWg0E0PYcyW8CzqDcVmaew9cujU4s=",
    version = "v1.10.2-0.20180628010248-798f60e20bb2",
)

go_repository(
    name = "com_github_oklog_run",
    importpath = "github.com/oklog/run",
    sum = "h1:R9vmquWCeGmxTHUVnTQJrU4oPlgEn9+x48nwXSqkIKg=",
    version = "v1.0.1-0.20180308005104-6934b124db28",
)

go_repository(
    name = "com_github_packethost_packngo",
    importpath = "github.com/packethost/packngo",
    sum = "h1:vwpFWvAO8DeIZfFeqASzZfsxuWPno9ncAebBEP0N3uE=",
    version = "v0.1.1-0.20180711074735-b9cb5096f54c",
)

go_repository(
    name = "com_github_pascaldekloe_goe",
    importpath = "github.com/pascaldekloe/goe",
    sum = "h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_patrickmn_go_cache",
    importpath = "github.com/patrickmn/go-cache",
    sum = "h1:HRMgzkcYKYpi3C8ajMPV8OFXaaRUnok+kx1WdO15EQc=",
    version = "v2.1.0+incompatible",
)

go_repository(
    name = "com_github_posener_complete",
    importpath = "github.com/posener/complete",
    sum = "h1:NP0eAhjcjImqslEwo/1hq7gpajME0fTLTezBKDqfXqo=",
    version = "v1.2.3",
)

go_repository(
    name = "com_github_rboyer_safeio",
    importpath = "github.com/rboyer/safeio",
    sum = "h1:05xhhdRNAdS3apYm7JRjOqngf4xruaW959jmRxGDuSU=",
    version = "v0.2.1",
)

go_repository(
    name = "com_github_renier_xmlrpc",
    importpath = "github.com/renier/xmlrpc",
    sum = "h1:Wdi9nwnhFNAlseAOekn6B5G/+GMtks9UKbvRU/CMM/o=",
    version = "v0.0.0-20170708154548-ce4a1a486c03",
)

go_repository(
    name = "com_github_ryanuber_columnize",
    importpath = "github.com/ryanuber/columnize",
    sum = "h1:ZVvb3R06BXNcQ3nXT2eANJAwGeRu4tbUzLI7sIyFZ2s=",
    version = "v2.1.1-0.20170703205827-abc90934186a+incompatible",
)

go_repository(
    name = "com_github_ryanuber_go_glob",
    importpath = "github.com/ryanuber/go-glob",
    sum = "h1:iQh3xXAumdQ+4Ufa5b25cRpC5TYKlno6hsv6Cb3pkBk=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_sean_conswriter",
    importpath = "github.com/sean-/conswriter",
    sum = "h1:Tn2Iev07a4oOcAuFna8AJxDOF/M+6OkNbpEZLX30D6M=",
    version = "v0.0.0-20180208195008-f5ae3917a627",
)

go_repository(
    name = "com_github_sean_pager",
    importpath = "github.com/sean-/pager",
    sum = "h1:D07EBYJLI26GmLRKNtrs47p8vs/5QqpUX3VcwsAPkEo=",
    version = "v0.0.0-20180208200047-666be9bf53b5",
)

go_repository(
    name = "com_github_sean_seed",
    importpath = "github.com/sean-/seed",
    sum = "h1:nn5Wsu0esKSJiIVhscUtVbo7ada43DJhG55ua/hjS5I=",
    version = "v0.0.0-20170313163322-e2103e2c3529",
)

go_repository(
    name = "com_github_shirou_gopsutil_v3",
    importpath = "github.com/shirou/gopsutil/v3",
    sum = "h1:7zomV9HJv6UGk225YtvEa5+camNLpbua3MAz/GqiVJY=",
    version = "v3.20.10",
)

go_repository(
    name = "com_github_softlayer_softlayer_go",
    importpath = "github.com/softlayer/softlayer-go",
    sum = "h1:bVQRCxQvfjNUeRqaY/uT0tFuvuFY0ulgnczuR684Xic=",
    version = "v0.0.0-20180806151055-260589d94c7d",
)

go_repository(
    name = "com_github_tencentcloud_tencentcloud_sdk_go",
    importpath = "github.com/tencentcloud/tencentcloud-sdk-go",
    sum = "h1:8uRvJleFpqLsO77WaAh2UrasMOzd8MxXrNj20e7El+Q=",
    version = "v3.0.83+incompatible",
)

go_repository(
    name = "com_github_tent_http_link_go",
    importpath = "github.com/tent/http-link-go",
    sum = "h1:/Bsw4C+DEdqPjt8vAqaC9LAqpAQnaCQQqmolqq3S1T4=",
    version = "v0.0.0-20130702225549-ac974c61c2f9",
)

go_repository(
    name = "com_github_tv42_httpunix",
    importpath = "github.com/tv42/httpunix",
    sum = "h1:G3dpKMzFDjgEh2q1Z7zUUtKa8ViPtH+ocF0bE0g00O8=",
    version = "v0.0.0-20150427012821-b75d8614f926",
)

go_repository(
    name = "in_gopkg_airbrake_gobrake_v2",
    importpath = "gopkg.in/airbrake/gobrake.v2",
    sum = "h1:7z2uVWwn7oVeeugY1DtlPAy5H+KYgB1KeKTnqjNatLo=",
    version = "v2.0.9",
)

go_repository(
    name = "in_gopkg_gemnasium_logrus_airbrake_hook_v2",
    importpath = "gopkg.in/gemnasium/logrus-airbrake-hook.v2",
    sum = "h1:OAj3g0cR6Dx/R07QgQe8wkA9RNjB2u4i700xBkIT4e0=",
    version = "v2.1.2",
)

go_repository(
    name = "net_launchpad_gocheck",
    importpath = "launchpad.net/gocheck",
    sum = "h1:Izowp2XBH6Ya6rv+hqbceQyw/gSGoXfH/UPoTGduL54=",
    version = "v0.0.0-20140225173054-000000000087",
)

go_repository(
    name = "org_uber_go_goleak",
    importpath = "go.uber.org/goleak",
    sum = "h1:z+mqJhf6ss6BSfSM671tgKyZBFPTTJM+HLxnhPC3wu0=",
    version = "v1.1.10",
)

go_repository(
    name = "com_github_jessevdk_go_flags",
    importpath = "github.com/jessevdk/go-flags",
    sum = "h1:4IU2WS7AumrZ/40jfhf4QVDMsQwqA7VEHozFRrGARJA=",
    version = "v1.4.0",
)

go_repository(
    name = "net_howett_plist",
    importpath = "howett.net/plist",
    sum = "h1:QDxUo/w2COstK1wIBYpzQlHX/NqaQTcf9jyz347nI58=",
    version = "v0.0.0-20201203080718-1454fab16a06",
)

go_repository(
    name = "com_github_agext_levenshtein",
    importpath = "github.com/agext/levenshtein",
    sum = "h1:QmvMAjj2aEICytGiWzmxoE0x2KZvE0fvmqMOfy2tjT8=",
    version = "v1.2.1",
)

go_repository(
    name = "com_github_apparentlymart_go_cidr",
    importpath = "github.com/apparentlymart/go-cidr",
    sum = "h1:NmIwLZ/KdsjIUlhf+/Np40atNXm/+lZ5txfTJ/SpF+U=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_apparentlymart_go_dump",
    importpath = "github.com/apparentlymart/go-dump",
    sum = "h1:ZSTrOEhiM5J5RFxEaFvMZVEAM1KvT1YzbEOwB2EAGjA=",
    version = "v0.0.0-20180507223929-23540a00eaa3",
)

go_repository(
    name = "com_github_apparentlymart_go_textseg",
    importpath = "github.com/apparentlymart/go-textseg",
    sum = "h1:rRmlIsPEEhUTIKQb7T++Nz/A5Q6C9IuX2wFoYVvnCs0=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_apparentlymart_go_textseg_v12",
    importpath = "github.com/apparentlymart/go-textseg/v12",
    sum = "h1:bNEQyAGak9tojivJNkoqWErVCQbjdL7GzRt3F8NvfJ0=",
    version = "v12.0.0",
)

go_repository(
    name = "com_github_azure_go_autorest",
    importpath = "github.com/Azure/go-autorest",
    sum = "h1:V5VMDjClD3GiElqLWO7mz2MxNAK/vTfRHdAubSIPRgs=",
    version = "v14.2.0+incompatible",
)

go_repository(
    name = "com_github_bgentry_go_netrc",
    importpath = "github.com/bgentry/go-netrc",
    sum = "h1:xDfNPAt8lFiC1UJrqV3uuy861HCTo708pDMbjHHdCas=",
    version = "v0.0.0-20140422174119-9fd32a8b3d3d",
)

go_repository(
    name = "com_github_bitly_go_simplejson",
    importpath = "github.com/bitly/go-simplejson",
    sum = "h1:6IH+V8/tVMab511d5bn4M7EwGXZf9Hj6i2xSwkNEM+Y=",
    version = "v0.5.0",
)

go_repository(
    name = "com_github_bmatcuk_doublestar",
    importpath = "github.com/bmatcuk/doublestar",
    sum = "h1:2bNwBOmhyFEFcoB3tGvTD5xanq+4kyOZlB8wFYbMjkk=",
    version = "v1.1.5",
)

go_repository(
    name = "com_github_bshuster_repo_logrus_logstash_hook",
    importpath = "github.com/bshuster-repo/logrus-logstash-hook",
    sum = "h1:pgAtgj+A31JBVtEHu2uHuEx0n+2ukqUJnS2vVe5pQNA=",
    version = "v0.4.1",
)

go_repository(
    name = "com_github_bugsnag_bugsnag_go",
    importpath = "github.com/bugsnag/bugsnag-go",
    sum = "h1:rFt+Y/IK1aEZkEHchZRSq9OQbsSzIT/OrI8YFFmRIng=",
    version = "v0.0.0-20141110184014-b1d153021fcd",
)

go_repository(
    name = "com_github_bugsnag_osext",
    importpath = "github.com/bugsnag/osext",
    sum = "h1:otBG+dV+YK+Soembjv71DPz3uX/V/6MMlSyD9JBQ6kQ=",
    version = "v0.0.0-20130617224835-0dd3f918b21b",
)

go_repository(
    name = "com_github_bugsnag_panicwrap",
    importpath = "github.com/bugsnag/panicwrap",
    sum = "h1:nvj0OLI3YqYXer/kZD8Ri1aaunCxIEsOst1BVJswV0o=",
    version = "v0.0.0-20151223152923-e2c28503fcd0",
)

go_repository(
    name = "com_github_cheggaaa_pb",
    importpath = "github.com/cheggaaa/pb",
    sum = "h1:wIkZHkNfC7R6GI5w7l/PdAdzXzlrbcI3p8OAlnkTsnc=",
    version = "v1.0.27",
)

go_repository(
    name = "com_github_containerd_go_cni",
    importpath = "github.com/containerd/go-cni",
    sum = "h1:fV37ZKnYs79fSyI3mu/XZFJVezrVsXBLbfojcTPpdXM=",
    version = "v0.0.0-20190904155053-d20b7eebc7ee",
)

go_repository(
    name = "com_github_containernetworking_plugins",
    importpath = "github.com/containernetworking/plugins",
    sum = "h1:WT1OeJOkmYxVP5/e3P+/MVJ1ftii14o80u0N9aK2tCg=",
    version = "v0.7.3-0.20190501191748-2d6d46d308b2",
)

go_repository(
    name = "com_github_coreos_go_iptables",
    importpath = "github.com/coreos/go-iptables",
    sum = "h1:9n0vuQxZw5yuleXDNqGLSw3o5s6vZJJH36ZSU7626zw=",
    version = "v0.4.3-0.20190724151750-969b135e941d",
)

go_repository(
    name = "com_github_docker_cli",
    importpath = "github.com/docker/cli",
    sum = "h1:5/ZRsUbguX7xFNLlbxVQY/yhD3Psy+vylKZrNme5BJs=",
    version = "v0.0.0-20200303215952-eb310fca4956",
)

go_repository(
    name = "com_github_docker_docker_credential_helpers",
    importpath = "github.com/docker/docker-credential-helpers",
    sum = "h1:zI2p9+1NQYdnG6sMU26EX4aVGlqbInSQxQXLvzJ4RPQ=",
    version = "v0.6.3",
)

go_repository(
    name = "com_github_docker_go_events",
    importpath = "github.com/docker/go-events",
    sum = "h1:+pKlWGMw7gf6bQ+oDZB4KHQFypsfjYlq/C4rfL7D3g8=",
    version = "v0.0.0-20190806004212-e31b211e4f1c",
)

go_repository(
    name = "com_github_docker_go_metrics",
    importpath = "github.com/docker/go-metrics",
    sum = "h1:yWHOI+vFjEsAakUTSrtqc/SAHrhSkmn48pqjidZX3QA=",
    version = "v0.0.0-20180209012529-399ea8c73916",
)

go_repository(
    name = "com_github_docker_libnetwork",
    importpath = "github.com/docker/libnetwork",
    sum = "h1:RvyaXv8RgMpJIT73oQFOj4Vos1c1rwJcIm3sPHdd1Js=",
    version = "v0.8.0-dev.2.0.20200612180813-9e99af28df21",
)

go_repository(
    name = "com_github_docker_libtrust",
    importpath = "github.com/docker/libtrust",
    sum = "h1:ZClxb8laGDf5arXfYcAtECDFgAgHklGI8CxgjHnXKJ4=",
    version = "v0.0.0-20150114040149-fa567046d9b1",
)

go_repository(
    name = "com_github_fsouza_go_dockerclient",
    importpath = "github.com/fsouza/go-dockerclient",
    sum = "h1:vuFDnPcds3LvTWGYb9h0Rty14FLgkjHZdwLDROCdgsw=",
    version = "v1.6.5",
)

go_repository(
    name = "com_github_garyburd_redigo",
    importpath = "github.com/garyburd/redigo",
    sum = "h1:LofdAjjjqCSXMwLGgOgnE+rdPuvX9DxCqaHwKy7i/ko=",
    version = "v0.0.0-20150301180006-535138d7bcd7",
)

go_repository(
    name = "com_github_go_ldap_ldap",
    importpath = "github.com/go-ldap/ldap",
    sum = "h1:kD5HQcAzlQ7yrhfn+h+MSABeAy/jAJhvIJ/QDllP44g=",
    version = "v3.0.2+incompatible",
)

go_repository(
    name = "com_github_godbus_dbus",
    importpath = "github.com/godbus/dbus",
    sum = "h1:WqqLRTsQic3apZUK9qC5sGNfXthmPXzUZ7nQPrNITa4=",
    version = "v4.1.0+incompatible",
)

go_repository(
    name = "com_github_gogo_googleapis",
    importpath = "github.com/gogo/googleapis",
    sum = "h1:Z0v3OJDotX9ZBpdz2V+AI7F4fITSZhVE5mg6GQppwMM=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_hashicorp_consul_template",
    importpath = "github.com/hashicorp/consul-template",
    sum = "h1:+D2s8eyRqWyX7GPNxeUi8tsyh8pRn3J6k8giEchPfKQ=",
    version = "v0.25.1",
)

go_repository(
    name = "com_github_hashicorp_cronexpr",
    importpath = "github.com/hashicorp/cronexpr",
    sum = "h1:NJZDd87hGXjoZBdvyCF9mX4DCq5Wy7+A/w+A7q0wn6c=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_hashicorp_go_cty_funcs",
    importpath = "github.com/hashicorp/go-cty-funcs",
    sum = "h1:kgvybwEeu0SXktbB2y3uLHX9lklLo+nzUwh59A3jzQc=",
    version = "v0.0.0-20200930094925-2721b1e36840",
)

go_repository(
    name = "com_github_hashicorp_go_envparse",
    importpath = "github.com/hashicorp/go-envparse",
    sum = "h1:HTmDIaSN95gbdMyrsbNiXSdW4fbGctGQwEqv0H7OhDQ=",
    version = "v0.0.0-20180119215841-310ca1881b22",
)

go_repository(
    name = "com_github_hashicorp_go_gatedio",
    importpath = "github.com/hashicorp/go-gatedio",
    sum = "h1:Jm1X5yP4yCqqWj5L1TgW7iZwCVPGtVc+mro5r/XX7Tg=",
    version = "v0.5.0",
)

go_repository(
    name = "com_github_hashicorp_go_getter",
    importpath = "github.com/hashicorp/go-getter",
    sum = "h1:XDo8LiAcDisiqZdv0TKgz+HtX3WN7zA2JD1R1tjsabE=",
    version = "v1.5.2",
)

go_repository(
    name = "com_github_hashicorp_go_net",
    importpath = "github.com/hashicorp/go.net",
    sum = "h1:sNCoNyDEvN1xa+X0baata4RdcpKwcMS6DH+xwfqPgjw=",
    version = "v0.0.1",
)

go_repository(
    name = "com_github_hashicorp_go_safetemp",
    importpath = "github.com/hashicorp/go-safetemp",
    sum = "h1:2HR189eFNrjHQyENnQMMpCiBAsRxzbTMIgBhEyExpmo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_hashicorp_hcl_v2",
    importpath = "github.com/hashicorp/hcl/v2",
    sum = "h1:iaau0VStfX9CgOlpbceawI94uVEM3sliqnjpHSVQqUo=",
    version = "v2.7.1-0.20201020204811-68a97f93bb48",
)

go_repository(
    name = "com_github_hashicorp_nomad",
    importpath = "github.com/hashicorp/nomad",
    sum = "h1:S2Vdo6mGwegzxcI921ee9OQJd+3OZ8ginf6xORKAbeI=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_hashicorp_nomad_api",
    importpath = "github.com/hashicorp/nomad/api",
    sum = "h1:FJBOmEYvJl+T7N5fQm96C8It1dfy91BgcsQN9srBUpM=",
    version = "v0.0.0-20210115191909-bcd4752fc902",
)

go_repository(
    name = "com_github_jmespath_go_jmespath_internal_testify",
    importpath = "github.com/jmespath/go-jmespath/internal/testify",
    sum = "h1:shLQSRRSCCPj3f2gpwzGwWFoC7ycTf1rcQZHOlsJ6N8=",
    version = "v1.5.1",
)

go_repository(
    name = "com_github_likexian_gokit",
    importpath = "github.com/likexian/gokit",
    sum = "h1:8ypmVXLx8yIvlTwzH8Ybz8LDAfWjdy0W5O354JWPjA4=",
    version = "v0.20.16",
)

go_repository(
    name = "com_github_likexian_simplejson_go",
    importpath = "github.com/likexian/simplejson-go",
    sum = "h1:tNa5q0zLsg6EzUI4npXeUYwk/xuTtMcKdVnRsYdOHzk=",
    version = "v0.0.0-20190502021454-d8787b4bfa0b",
)

go_repository(
    name = "com_github_lk4d4_joincontext",
    importpath = "github.com/LK4D4/joincontext",
    sum = "h1:U7q69tqXiCf6m097GRlNQB0/6SI1qWIOHYHhCEvDxF4=",
    version = "v0.0.0-20171026170139-1724345da6d5",
)

go_repository(
    name = "com_github_marstr_guid",
    importpath = "github.com/marstr/guid",
    sum = "h1:/M4H/1G4avsieL6BbUwCOBzulmoeKVP5ux/3mQNnbyI=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_mattn_go_shellwords",
    importpath = "github.com/mattn/go-shellwords",
    sum = "h1:JhhFTIOslh5ZsPrpa3Wdg8bF0WI3b44EMblmU9wIsXc=",
    version = "v1.0.5",
)

go_repository(
    name = "com_github_mitchellh_colorstring",
    importpath = "github.com/mitchellh/colorstring",
    sum = "h1:KHyL+3mQOF9sPfs26lsefckcFNDcIZtiACQiECzIUkw=",
    version = "v0.0.0-20150917214807-8631ce90f286",
)

go_repository(
    name = "com_github_mitchellh_gox",
    importpath = "github.com/mitchellh/gox",
    sum = "h1:lfGJxY7ToLJQjHHwi0EX6uYBdK78egf954SQl13PQJc=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_mitchellh_iochan",
    importpath = "github.com/mitchellh/iochan",
    sum = "h1:C+X3KsSTLFVBr/tK1eYN/vs4rJcvsiLU338UhYPJWeY=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_mitchellh_osext",
    importpath = "github.com/mitchellh/osext",
    sum = "h1:2+myh5ml7lgEU/51gbeLHfKGNfgEQQIWrlbdaOsidbQ=",
    version = "v0.0.0-20151018003038-5e2d6d41470f",
)

go_repository(
    name = "com_github_ncw_swift",
    importpath = "github.com/ncw/swift",
    sum = "h1:4DQRPj35Y41WogBxyhOXlrI37nzGlyEcsforeudyYPQ=",
    version = "v1.0.47",
)

go_repository(
    name = "com_github_niemeyer_pretty",
    importpath = "github.com/niemeyer/pretty",
    sum = "h1:fD57ERR4JtEqsWbfPhv4DMiApHyliiK5xCTNVSPiaAs=",
    version = "v0.0.0-20200227124842-a10e7caefd8e",
)

go_repository(
    name = "com_github_nvidia_gpu_monitoring_tools",
    importpath = "github.com/NVIDIA/gpu-monitoring-tools",
    sum = "h1:WLyvLAM0QfjAarRzRTG9EgT5McqGWNZMvqqSUSoyUUY=",
    version = "v0.0.0-20180829222009-86f2a9fac6c5",
)

go_repository(
    name = "com_github_opencontainers_runtime_tools",
    importpath = "github.com/opencontainers/runtime-tools",
    sum = "h1:H7DMc6FAjgwZZi8BRqjrAAHWoqEr5e5L6pS4V0ezet4=",
    version = "v0.0.0-20181011054405-1d69bd0f9c39",
)

go_repository(
    name = "com_github_shopify_logrus_bugsnag",
    importpath = "github.com/Shopify/logrus-bugsnag",
    sum = "h1:UrqY+r/OJnIp5u0s1SbQ8dVfLCZJsnvazdBP5hS4iRs=",
    version = "v0.0.0-20171204204709-577dee27f20d",
)

go_repository(
    name = "com_github_skratchdot_open_golang",
    importpath = "github.com/skratchdot/open-golang",
    sum = "h1:fyKiXKO1/I/B6Y2U8T7WdQGWzwehOuGIrljPtt7YTTI=",
    version = "v0.0.0-20160302144031-75fb7ed4208c",
)

go_repository(
    name = "com_github_vmihailenco_msgpack",
    importpath = "github.com/vmihailenco/msgpack",
    sum = "h1:wapg9xDUZDzGCNFlwc5SqI1rvcciqcxEHac4CYj89xI=",
    version = "v3.3.3+incompatible",
)

go_repository(
    name = "com_github_willf_bitset",
    importpath = "github.com/willf/bitset",
    sum = "h1:R43TdZy32XXSXjJn7M/HhALJ9imq6ztLnChfYJpVDnM=",
    version = "v1.1.11-0.20200630133818-d5bec3311243",
)

go_repository(
    name = "com_github_xeipuuv_gojsonpointer",
    importpath = "github.com/xeipuuv/gojsonpointer",
    sum = "h1:J9EGpcZtP0E/raorCMxlFGSTBrsSlaDGf3jU/qvAE2c=",
    version = "v0.0.0-20180127040702-4e3ac2762d5f",
)

go_repository(
    name = "com_github_xeipuuv_gojsonreference",
    importpath = "github.com/xeipuuv/gojsonreference",
    sum = "h1:EzJWgHovont7NscjpAxXsDA8S8BMYve8Y5+7cuRE7R0=",
    version = "v0.0.0-20180127040603-bd5ef7bd5415",
)

go_repository(
    name = "com_github_xeipuuv_gojsonschema",
    importpath = "github.com/xeipuuv/gojsonschema",
    sum = "h1:mvXjJIHRZyhNuGassLTcXTwjiWq7NmjdavZsUnmFybQ=",
    version = "v0.0.0-20180618132009-1d523034197f",
)

go_repository(
    name = "com_github_yvasiyarov_go_metrics",
    importpath = "github.com/yvasiyarov/go-metrics",
    sum = "h1:+lm10QQTNSBd8DVTNGHx7o/IKu9HYDvLMffDhbyLccI=",
    version = "v0.0.0-20140926110328-57bccd1ccd43",
)

go_repository(
    name = "com_github_yvasiyarov_gorelic",
    importpath = "github.com/yvasiyarov/gorelic",
    sum = "h1:hlE8//ciYMztlGpl/VA+Zm1AcTPHYkHJPbHqE6WJUXE=",
    version = "v0.0.0-20141212073537-a9bba5b9ab50",
)

go_repository(
    name = "com_github_yvasiyarov_newrelic_platform_go",
    importpath = "github.com/yvasiyarov/newrelic_platform_go",
    sum = "h1:ERexzlUfuTvpE74urLSbIQW0Z/6hF9t8U4NsJLaioAY=",
    version = "v0.0.0-20140908184405-b21fdbd4370f",
)

go_repository(
    name = "com_github_zclconf_go_cty",
    importpath = "github.com/zclconf/go-cty",
    sum = "h1:Xzr4m4utRDhHDifag1onwwUSq32HLoLBsp+w6tD0880=",
    version = "v1.4.1",
)

go_repository(
    name = "com_github_zclconf_go_cty_yaml",
    importpath = "github.com/zclconf/go-cty-yaml",
    sum = "h1:dNyg4QLTrv2IfJpm7Wtxi55ed5gLGOlPrZ6kMd51hY0=",
    version = "v1.0.2",
)

go_repository(
    name = "in_gopkg_asn1_ber_v1",
    importpath = "gopkg.in/asn1-ber.v1",
    sum = "h1:TxyelI5cVkbREznMhfzycHdkp5cLA7DpE+GKjSslYhM=",
    version = "v1.0.0-20181015200546-f715ec2f112d",
)

go_repository(
    name = "in_gopkg_tomb_v2",
    importpath = "gopkg.in/tomb.v2",
    sum = "h1:EQ3aCG3c3nkUNxx6quE0Ux47RYExj7cJyRMxUXqPf6I=",
    version = "v2.0.0-20140626144623-14b3d72120e8",
)

go_repository(
    name = "io_k8s_api",
    importpath = "k8s.io/api",
    sum = "h1:wG5g5ZmSVgm5B+eHMIbI9EGATS2L8Z72rda19RIEgY8=",
    version = "v0.18.2",
)

go_repository(
    name = "io_k8s_apimachinery",
    importpath = "k8s.io/apimachinery",
    sum = "h1:44CmtbmkzVDAhCpRVSiP2R5PPrC2RtlIv/MoB8xpdRA=",
    version = "v0.18.2",
)

go_repository(
    name = "io_k8s_client_go",
    importpath = "k8s.io/client-go",
    sum = "h1:tTI4hRmb1DRMl4fG6Vclfdi6nTM82oIrTT7HfitmxC4=",
    version = "v8.0.0+incompatible",
)

go_repository(
    name = "io_k8s_cri_api",
    importpath = "k8s.io/cri-api",
    sum = "h1:jvjVvBqgZq3WcaPq07n0h5h9eCnIaR4dhKyHSoZG8Y8=",
    version = "v0.17.3",
)

go_repository(
    name = "io_k8s_gengo",
    importpath = "k8s.io/gengo",
    sum = "h1:JApXBKYyB7l9xx+DK7/+mFjC7A9Bt5A93FPvFD0HIFE=",
    version = "v0.0.0-20201113003025-83324d819ded",
)

go_repository(
    name = "io_k8s_klog",
    importpath = "k8s.io/klog",
    sum = "h1:Pt+yjF5aB1xDSVbau4VsWe+dQNzA0qv1LlXdC2dF6Q8=",
    version = "v1.0.0",
)

go_repository(
    name = "io_k8s_kube_openapi",
    importpath = "k8s.io/kube-openapi",
    sum = "h1:sOHNzJIkytDF6qadMNKhhDRpc6ODik8lVC6nOur7B2c=",
    version = "v0.0.0-20201113171705-d219536bb9fd",
)

go_repository(
    name = "io_k8s_kubernetes",
    importpath = "k8s.io/kubernetes",
    sum = "h1:qTfB+u5M92k2fCCCVP2iuhgwwSOv1EkAkvQY1tQODD8=",
    version = "v1.13.0",
)

go_repository(
    name = "io_k8s_sigs_structured_merge_diff_v3",
    importpath = "sigs.k8s.io/structured-merge-diff/v3",
    sum = "h1:dOmIZBMfhcHS09XZkMyUgkq5trg3/jRyJYFZUiaOp8E=",
    version = "v3.0.0",
)

go_repository(
    name = "io_k8s_sigs_yaml",
    importpath = "sigs.k8s.io/yaml",
    sum = "h1:kr/MCeFWJWTwyaHoR9c8EjH9OumOmoF9YGiZd7lFm/Q=",
    version = "v1.2.0",
)

go_repository(
    name = "io_k8s_utils",
    importpath = "k8s.io/utils",
    sum = "h1:d4vVOjXm687F1iLSP2q3lyPPuyvTUt3aVoBpi2DqRsU=",
    version = "v0.0.0-20200324210504-a9aa75ae1b89",
)

go_repository(
    name = "org_bazil_fuse",
    importpath = "bazil.org/fuse",
    sum = "h1:SC+c6A1qTFstO9qmB86mPV2IpYme/2ZoEQ0hrP+wo+Q=",
    version = "v0.0.0-20160811212531-371fbbdaa898",
)

go_repository(
    name = "org_golang_google_cloud",
    importpath = "google.golang.org/cloud",
    sum = "h1:Cpp2P6TPjujNoC5M2KHY6g7wfyLYfIWRZaSdIKfDasA=",
    version = "v0.0.0-20151119220103-975617b05ea8",
)

go_repository(
    name = "com_github_containerd_stargz_snapshotter_estargz",
    importpath = "github.com/containerd/stargz-snapshotter/estargz",
    sum = "h1:OVQ4FVXeE6OjzuUifzER+7EulqTqw/94oKSqnooEowQ=",
    version = "v0.0.0-20201223015020-a9a0c2d64694",
)

go_repository(
    name = "com_github_google_go_containerregistry",
    importpath = "github.com/google/go-containerregistry",
    sum = "h1:45axtqLd66llqD8R9XgiCQ64foc7I2xkAG40NwR5YFw=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_joefitzgerald_rainbow_reporter",
    importpath = "github.com/joefitzgerald/rainbow-reporter",
    sum = "h1:AuMG652zjdzI0YCCnXAqATtRBpGXMcAnrajcaTrSeuo=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_maxbrunsfeld_counterfeiter_v6",
    importpath = "github.com/maxbrunsfeld/counterfeiter/v6",
    sum = "h1:g+4J5sZg6osfvEfkRZxJ1em0VT95/UOZgi/l7zi1/oE=",
    version = "v6.2.2",
)

go_repository(
    name = "com_github_sclevine_spec",
    importpath = "github.com/sclevine/spec",
    sum = "h1:1Jwdf9jSfDl9NVmt8ndHqbTZ7XCCPbh1jI3hkDBHVYA=",
    version = "v1.2.0",
)

go_repository(
    name = "io_k8s_code_generator",
    importpath = "k8s.io/code-generator",
    sum = "h1:kre3GNich5gbO3d1FyTT8fHI4ZJezZV217yFdWlQaRQ=",
    version = "v0.20.1",
)

go_repository(
    name = "io_k8s_klog_v2",
    importpath = "k8s.io/klog/v2",
    sum = "h1:7+X0fUguPyrKEC4WjH8iGDg3laWgMo5tMnRTIGTTxGQ=",
    version = "v2.4.0",
)

go_repository(
    name = "io_k8s_sigs_structured_merge_diff_v4",
    importpath = "sigs.k8s.io/structured-merge-diff/v4",
    sum = "h1:YHQV7Dajm86OuqnIR6zAelnDWBRjo+YhYV9PmGrh1s8=",
    version = "v4.0.2",
)

go_repository(
    name = "com_github_afex_hystrix_go",
    importpath = "github.com/afex/hystrix-go",
    sum = "h1:rFw4nCn9iMW+Vajsk51NtYIcwSTkXr+JGrMd36kTDJw=",
    version = "v0.0.0-20180502004556-fa1af6a1f4f5",
)

go_repository(
    name = "com_github_aryann_difflib",
    importpath = "github.com/aryann/difflib",
    sum = "h1:pv34s756C4pEXnjgPfGYgdhg/ZdajGhyOvzx8k+23nw=",
    version = "v0.0.0-20170710044230-e206f873d14a",
)

go_repository(
    name = "com_github_aws_aws_lambda_go",
    importpath = "github.com/aws/aws-lambda-go",
    sum = "h1:SuCy7H3NLyp+1Mrfp+m80jcbi9KYWAs9/BXwppwRDzY=",
    version = "v1.13.3",
)

go_repository(
    name = "com_github_aws_aws_sdk_go_v2",
    importpath = "github.com/aws/aws-sdk-go-v2",
    sum = "h1:qZ+woO4SamnH/eEbjM2IDLhRNwIwND/RQyVlBLp3Jqg=",
    version = "v0.18.0",
)

go_repository(
    name = "com_github_casbin_casbin_v2",
    importpath = "github.com/casbin/casbin/v2",
    sum = "h1:bTwon/ECRx9dwBy2ewRVr5OiqjeXSGiTUY74sDPQi/g=",
    version = "v2.1.2",
)

go_repository(
    name = "com_github_clbanning_x2j",
    importpath = "github.com/clbanning/x2j",
    sum = "h1:EdRZT3IeKQmfCSrgo8SZ8V3MEnskuJP0wCYNpe+aiXo=",
    version = "v0.0.0-20191024224557-825249438eec",
)

go_repository(
    name = "com_github_codahale_hdrhistogram",
    importpath = "github.com/codahale/hdrhistogram",
    sum = "h1:qMd81Ts1T2OTKmB4acZcyKaMtRnY5Y44NuXGX2GFJ1w=",
    version = "v0.0.0-20161010025455-3a0bb77429bd",
)

go_repository(
    name = "com_github_eapache_go_resiliency",
    importpath = "github.com/eapache/go-resiliency",
    sum = "h1:1NtRmCAqadE2FN4ZcN6g90TP3uk8cg9rn9eNK2197aU=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_eapache_go_xerial_snappy",
    importpath = "github.com/eapache/go-xerial-snappy",
    sum = "h1:YEetp8/yCZMuEPMUDHG0CW/brkkEp8mzqk2+ODEitlw=",
    version = "v0.0.0-20180814174437-776d5712da21",
)

go_repository(
    name = "com_github_eapache_queue",
    importpath = "github.com/eapache/queue",
    sum = "h1:YOEu7KNc61ntiQlcEeUIoDTJ2o8mQznoNvUhiigpIqc=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_franela_goblin",
    importpath = "github.com/franela/goblin",
    sum = "h1:gb2Z18BhTPJPpLQWj4T+rfKHYCHxRHCtRxhKKjRidVw=",
    version = "v0.0.0-20200105215937-c9ffbefa60db",
)

go_repository(
    name = "com_github_franela_goreq",
    importpath = "github.com/franela/goreq",
    sum = "h1:a9ENSRDFBUPkJ5lCgVZh26+ZbGyoVJG7yb5SSzF5H54=",
    version = "v0.0.0-20171204163338-bcd34c9993f8",
)

go_repository(
    name = "com_github_hudl_fargo",
    importpath = "github.com/hudl/fargo",
    sum = "h1:0U6+BtN6LhaYuTnIJq4Wyq5cpn6O2kWrxAtcqBmYY6w=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_influxdata_influxdb1_client",
    importpath = "github.com/influxdata/influxdb1-client",
    sum = "h1:/WZQPMZNsjZ7IlCpsLGdQBINg5bxKQ1K1sh6awxLtkA=",
    version = "v0.0.0-20191209144304-8bf82d3c094d",
)

go_repository(
    name = "com_github_jpillora_backoff",
    importpath = "github.com/jpillora/backoff",
    sum = "h1:uvFg412JmmHBHw7iwprIxkPMI+sGQ4kzOWsMeHnm2EA=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_knetic_govaluate",
    importpath = "github.com/Knetic/govaluate",
    sum = "h1:1G1pk05UrOh0NlF1oeaaix1x8XzrfjIDK47TY0Zehcw=",
    version = "v3.0.1-0.20171022003610-9aa49832a739+incompatible",
)

go_repository(
    name = "com_github_lightstep_lightstep_tracer_common_golang_gogo",
    importpath = "github.com/lightstep/lightstep-tracer-common/golang/gogo",
    sum = "h1:143Bb8f8DuGWck/xpNUOckBVYfFbBTnLevfRZ1aVVqo=",
    version = "v0.0.0-20190605223551-bc2310a04743",
)

go_repository(
    name = "com_github_lightstep_lightstep_tracer_go",
    importpath = "github.com/lightstep/lightstep-tracer-go",
    sum = "h1:vi1F1IQ8N7hNWytK9DpJsUfQhGuNSc19z330K6vl4zk=",
    version = "v0.18.1",
)

go_repository(
    name = "com_github_lyft_protoc_gen_validate",
    importpath = "github.com/lyft/protoc-gen-validate",
    sum = "h1:KNt/RhmQTOLr7Aj8PsJ7mTronaFyx80mRTT9qF261dA=",
    version = "v0.0.13",
)

go_repository(
    name = "com_github_nats_io_jwt",
    importpath = "github.com/nats-io/jwt",
    sum = "h1:+RB5hMpXUUA2dfxuhBTEkMOrYmM+gKIZYS1KjSostMI=",
    version = "v0.3.2",
)

go_repository(
    name = "com_github_nats_io_nats_go",
    importpath = "github.com/nats-io/nats.go",
    sum = "h1:ik3HbLhZ0YABLto7iX80pZLPw/6dx3T+++MZJwLnMrQ=",
    version = "v1.9.1",
)

go_repository(
    name = "com_github_nats_io_nats_server_v2",
    importpath = "github.com/nats-io/nats-server/v2",
    sum = "h1:i2Ly0B+1+rzNZHHWtD4ZwKi+OU5l+uQo1iDHZ2PmiIc=",
    version = "v2.1.2",
)

go_repository(
    name = "com_github_nats_io_nkeys",
    importpath = "github.com/nats-io/nkeys",
    sum = "h1:6JrEfig+HzTH85yxzhSVbjHRJv9cn0p6n3IngIcM5/k=",
    version = "v0.1.3",
)

go_repository(
    name = "com_github_nats_io_nuid",
    importpath = "github.com/nats-io/nuid",
    sum = "h1:5iA8DT8V7q8WK2EScv2padNa/rTESc1KdnPw4TC2paw=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_oklog_oklog",
    importpath = "github.com/oklog/oklog",
    sum = "h1:wVfs8F+in6nTBMkA7CbRw+zZMIB7nNM825cM1wuzoTk=",
    version = "v0.3.2",
)

go_repository(
    name = "com_github_op_go_logging",
    importpath = "github.com/op/go-logging",
    sum = "h1:lDH9UUVJtmYCjyT0CI4q8xvlXPxeZ0gYCVvWbmPlp88=",
    version = "v0.0.0-20160315200505-970db520ece7",
)

go_repository(
    name = "com_github_opentracing_basictracer_go",
    importpath = "github.com/opentracing/basictracer-go",
    sum = "h1:YyUAhaEfjoWXclZVJ9sGoNct7j4TVk7lZWlQw5UXuoo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_opentracing_contrib_go_observer",
    importpath = "github.com/opentracing-contrib/go-observer",
    sum = "h1:lM6RxxfUMrYL/f8bWEUqdXrANWtrL7Nndbm9iFN0DlU=",
    version = "v0.0.0-20170622124052-a52f23424492",
)

go_repository(
    name = "com_github_openzipkin_contrib_zipkin_go_opentracing",
    importpath = "github.com/openzipkin-contrib/zipkin-go-opentracing",
    sum = "h1:ZCnq+JUrvXcDVhX/xRolRBZifmabN1HcS1wrPSvxhrU=",
    version = "v0.4.5",
)

go_repository(
    name = "com_github_openzipkin_zipkin_go",
    importpath = "github.com/openzipkin/zipkin-go",
    sum = "h1:nY8Hti+WKaP0cRsSeQ026wU03QsM762XBeCXBb9NAWI=",
    version = "v0.2.2",
)

go_repository(
    name = "com_github_pact_foundation_pact_go",
    importpath = "github.com/pact-foundation/pact-go",
    sum = "h1:OYkFijGHoZAYbOIb1LWXrwKQbMMRUv1oQ89blD2Mh2Q=",
    version = "v1.0.4",
)

go_repository(
    name = "com_github_performancecopilot_speed",
    importpath = "github.com/performancecopilot/speed",
    sum = "h1:2WnRzIquHa5QxaJKShDkLM+sc0JPuwhXzK8OYOyt3Vg=",
    version = "v3.0.0+incompatible",
)

go_repository(
    name = "com_github_pkg_profile",
    importpath = "github.com/pkg/profile",
    sum = "h1:F++O52m40owAmADcojzM+9gyjmMOY/T4oYJkgFDH8RE=",
    version = "v1.2.1",
)

go_repository(
    name = "com_github_rcrowley_go_metrics",
    importpath = "github.com/rcrowley/go-metrics",
    sum = "h1:9ZKAASQSHhDYGoxY8uLVpewe1GDZ2vu2Tr/vTdVAkFQ=",
    version = "v0.0.0-20181016184325-3113b8401b8a",
)

go_repository(
    name = "com_github_samuel_go_zookeeper",
    importpath = "github.com/samuel/go-zookeeper",
    sum = "h1:p3Vo3i64TCLY7gIfzeQaUJ+kppEO5WQG3cL8iE8tGHU=",
    version = "v0.0.0-20190923202752-2cc03de413da",
)

go_repository(
    name = "com_github_shopify_sarama",
    importpath = "github.com/Shopify/sarama",
    sum = "h1:9oksLxC6uxVPHPVYUmq6xhr1BOF/hHobWH2UzO67z1s=",
    version = "v1.19.0",
)

go_repository(
    name = "com_github_shopify_toxiproxy",
    importpath = "github.com/Shopify/toxiproxy",
    sum = "h1:TKdv8HiTLgE5wdJuEML90aBgNWsokNbMijUGhmcoBJc=",
    version = "v2.1.4+incompatible",
)

go_repository(
    name = "com_github_sony_gobreaker",
    importpath = "github.com/sony/gobreaker",
    sum = "h1:oMnRNZXX5j85zso6xCPRNPtmAycat+WcoKbklScLDgQ=",
    version = "v0.4.1",
)

go_repository(
    name = "com_github_streadway_amqp",
    importpath = "github.com/streadway/amqp",
    sum = "h1:WhxRHzgeVGETMlmVfqhRn8RIeeNoPr2Czh33I4Zdccw=",
    version = "v0.0.0-20190827072141-edfb9018d271",
)

go_repository(
    name = "com_github_streadway_handy",
    importpath = "github.com/streadway/handy",
    sum = "h1:AhmOdSHeswKHBjhsLs/7+1voOxT+LLrSk/Nxvk35fug=",
    version = "v0.0.0-20190108123426-d5acb3125c2a",
)

go_repository(
    name = "com_github_vividcortex_gohistogram",
    importpath = "github.com/VividCortex/gohistogram",
    sum = "h1:6+hBz+qvs0JOrrNhhmR7lFxo5sINxBCGXrdtl/UvroE=",
    version = "v1.0.0",
)

go_repository(
    name = "com_sourcegraph_sourcegraph_appdash",
    importpath = "sourcegraph.com/sourcegraph/appdash",
    sum = "h1:ucqkfpjg9WzSUubAO62csmucvxl4/JeW3F4I4909XkM=",
    version = "v0.0.0-20190731080439-ebfcffb1b5c0",
)

go_repository(
    name = "org_uber_go_tools",
    importpath = "go.uber.org/tools",
    sum = "h1:0mgffUl7nfd+FpvXMVz4IDEaUSmT1ysygQC7qYo7sG4=",
    version = "v0.0.0-20190618225709-2cfd321de3ee",
)

go_repository(
    name = "ht_sr_git_apreiml_prunef",
    importpath = "git.sr.ht/~apreiml/prunef",
    sum = "h1:2srkQb23j/5UgzvnHOz9joJdkxHybq8pOckP+RnyuRk=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_bruj0_vault_plugin_auth_u2f",
    importpath = "github.com/bruj0/vault-plugin-auth-u2f",
    sum = "h1:iNmRVI+y2FPt0LHrLcUPwNoM2iiwJs99m8Lg98S196o=",
    version = "v0.0.0-20201021131703-3455d429e361",
)

go_repository(
    name = "com_github_ryankurte_go_u2f",
    importpath = "github.com/ryankurte/go-u2f",
    sum = "h1:TbsDz1r7FrSqdEUH9tbVlDHz39P1w2v3ydvu0EtZu3w=",
    version = "v0.1.4",
)

go_repository(
    name = "com_github_cloudflare_cfssl",
    importpath = "github.com/cloudflare/cfssl",
    sum = "h1:Puu1hUwfps3+1CUzYdAZXijuvLuRMirgiXdf3zsM2Ig=",
    version = "v0.0.0-20190726000631-633726f6bcb7",
)

go_repository(
    name = "com_github_duo_labs_webauthn",
    importpath = "github.com/duo-labs/webauthn",
    sum = "h1:eEEfwrmEwl0LVuWz/VkAefdgtPbX174Huu5dxxceihI=",
    version = "v0.0.0-20200714211715-1daaee874e43",
)

go_repository(
    name = "com_github_fxamacker_cbor_v2",
    importpath = "github.com/fxamacker/cbor/v2",
    sum = "h1:6eXqdDDe588rSYAi1HfZKbx6YYQO4mxQ9eC6xYpU/JQ=",
    version = "v2.2.0",
)

go_repository(
    name = "com_github_google_certificate_transparency_go",
    importpath = "github.com/google/certificate-transparency-go",
    sum = "h1:Yf1aXowfZ2nuboBsg7iYGLmwsOARdV86pfH3g95wXmE=",
    version = "v1.0.21",
)

go_repository(
    name = "com_github_x448_float16",
    importpath = "github.com/x448/float16",
    sum = "h1:qLwI1I70+NjRFUR3zs1JPUCgaCXSh3SW62uAKT1mSBM=",
    version = "v0.8.4",
)

go_repository(
    name = "com_github_gorilla_securecookie",
    importpath = "github.com/gorilla/securecookie",
    sum = "h1:miw7JPhV+b/lAHSXz4qd/nN9jRiAFV5FwjeKyCS8BvQ=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_gorilla_sessions",
    importpath = "github.com/gorilla/sessions",
    sum = "h1:DHd3rPN5lE3Ts3D8rKkQ8x/0kqfeNmBAaiSi+o7FsgI=",
    version = "v1.2.1",
)

go_repository(
    name = "com_github_gregdel_pushover",
    importpath = "github.com/gregdel/pushover",
    sum = "h1:jb9bDQ07E8YdegHcvFJ6r0JAQhSxT6BQjDLn/xAMu64=",
    version = "v0.0.0-20201104094836-ddbe0c1d3a38",
)

go_repository(
    name = "io_opentelemetry_go_otel_exporters_otlp",
    importpath = "go.opentelemetry.io/otel/exporters/otlp",
    sum = "h1:ez8agFGbFJJgBU9H3lfX0rxWhZlXqurgZKL4aDcOdqY=",
    version = "v0.19.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_metric",
    importpath = "go.opentelemetry.io/otel/metric",
    sum = "h1:dtZ1Ju44gkJkYvo+3qGqVXmf88tc+a42edOywypengg=",
    version = "v0.19.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_oteltest",
    importpath = "go.opentelemetry.io/otel/oteltest",
    sum = "h1:YVfA0ByROYqTwOxqHVZYZExzEpfZor+MU1rU+ip2v9Q=",
    version = "v0.19.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_sdk_export_metric",
    importpath = "go.opentelemetry.io/otel/sdk/export/metric",
    sum = "h1:9A1PC2graOx3epRLRWbq4DPCdpMUYK8XeCrdAg6ycbI=",
    version = "v0.19.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_sdk_metric",
    importpath = "go.opentelemetry.io/otel/sdk/metric",
    sum = "h1:fka1Zc/lpRMS+KlTP/TRXZuaFtSjUg/maHV3U8rt1Mc=",
    version = "v0.19.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_trace",
    importpath = "go.opentelemetry.io/otel/trace",
    sum = "h1:1ucYlenXIDA1OlHVLDZKX0ObXV5RLaq06DtUKz5e5zc=",
    version = "v0.19.0",
)

go_repository(
    name = "org_golang_x_term",
    importpath = "golang.org/x/term",
    sum = "h1:v+OssWQX+hTHEmOBgwxdZxK4zHq3yOs8F9J7mk0PY8E=",
    version = "v0.0.0-20201126162022-7de9c90e9dd1",
)

go_repository(
    name = "com_github_urfave_cli_v2",
    importpath = "github.com/urfave/cli/v2",
    sum = "h1:qph92Y649prgesehzOrQjdWyxFOp/QVM+6imKHad91M=",
    version = "v2.3.0",
)

gazelle_dependencies()
