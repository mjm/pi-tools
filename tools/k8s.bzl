load("@io_bazel_rules_k8s//k8s:object.bzl", "k8s_object")

def k8s_http_service(name, namespace, service_name, app = None, target_port = "http", **kwargs):
    _k8s_http_service(
        name = name + ".yaml",
        namespace = namespace,
        service_name = service_name,
        app = app,
        target_port = target_port,
    )

    k8s_object(
        name = name,
        kind = "service",
        template = name + ".yaml",
    )

def _k8s_http_service_impl(ctx):
    app = ctx.attr.app
    if app == "":
        app = ctx.attr.service_name

    template = ctx.file._template
    subs = {
      "{NAME}": ctx.attr.service_name,
      "{NAMESPACE}": ctx.attr.namespace,
      "{APP}": app,
      "{TARGET_PORT}": ctx.attr.target_port,
    }

    out = ctx.actions.declare_file(ctx.label.name)
    ctx.actions.expand_template(
        output = out,
        template = template,
        substitutions = subs,
    )
    return [DefaultInfo(files = depset([out]))]

_k8s_http_service = rule(
    implementation = _k8s_http_service_impl,
    attrs = {
        "service_name": attr.string(),
        "namespace": attr.string(),
        "app": attr.string(),
        "target_port": attr.string(),
        "_template": attr.label(
            allow_single_file = True,
            default = "//tools:k8s_service.yaml.tpl",
        ),
    },
)
