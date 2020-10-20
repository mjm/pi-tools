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
        **kwargs,
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

def k8s_virtual_host(name, namespace, service_name, host_name = None, port = 80, **kwargs):
    _k8s_virtual_host(
        name = name + ".yaml",
        namespace = namespace,
        service_name = service_name,
        host_name = host_name,
        port = port,
    )

    k8s_object(
        name = name,
        kind = "ingress",
        template = name + ".yaml",
        **kwargs,
    )

def _k8s_virtual_host_impl(ctx):
    host_name = ctx.attr.host_name
    if host_name == "":
        host_name = ctx.attr.service_name

    template = ctx.file._template
    subs = {
      "{NAME}": host_name,
      "{NAMESPACE}": ctx.attr.namespace,
      "{SERVICE_NAME}": ctx.attr.service_name,
      "{PORT}": "{}".format(ctx.attr.port),
    }

    out = ctx.actions.declare_file(ctx.label.name)
    ctx.actions.expand_template(
        output = out,
        template = template,
        substitutions = subs,
    )
    return [DefaultInfo(files = depset([out]))]

_k8s_virtual_host = rule(
    implementation = _k8s_virtual_host_impl,
    attrs = {
        "service_name": attr.string(),
        "namespace": attr.string(),
        "host_name": attr.string(),
        "port": attr.int(),
        "_template": attr.label(
            allow_single_file = True,
            default = "//tools:k8s_virtual_host.yaml.tpl",
        ),
    },
)
