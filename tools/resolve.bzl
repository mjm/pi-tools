load("@io_bazel_rules_docker//container:providers.bzl", "PushInfo")

def _impl(ctx):
    output_file = ctx.actions.declare_file(ctx.label.name)

    args = ctx.actions.args()
    args.add(output_file)

    inputs = []

    for push in ctx.attr.pushes:
        push_info = push[PushInfo]
        args.add("{}/{}".format(push_info.registry, push_info.repository))
        args.add(push_info.digest.path)
        inputs.append(push_info.digest)

    ctx.actions.run(
        mnemonic = "AssembleDigests",
        executable = ctx.executable._helper,
        arguments = [args],
        inputs = inputs,
        outputs = [output_file],
    )

    return [
        DefaultInfo(
            files = depset([output_file]),
        ),
    ]

image_digests = rule(
    implementation = _impl,
    attrs = {
        "pushes": attr.label_list(
            providers = [PushInfo],
        ),
        "_helper": attr.label(
            default = Label("//tools:image_digests.sh"),
            allow_single_file = True,
            executable = True,
            cfg = "exec",
        )
    },
)
