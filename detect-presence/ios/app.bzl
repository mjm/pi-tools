load("@build_bazel_rules_apple//apple:ios.bzl", "ios_application")

def presence_app(name, bundle_name, provisioning_profile):
    ios_application(
        name = name,
        app_icons = ["//detect-presence/ios/Resources:app_icon"],
        bundle_id = "com.mattmoriarity.Presence",
        bundle_name = bundle_name,
        families = ["iphone"],
        infoplists = [":Info.plist"],
        minimum_os_version = "14.0",
        provisioning_profile = provisioning_profile,
        deps = [":app-lib"],
    )
