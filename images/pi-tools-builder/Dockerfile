FROM ubuntu:focal

RUN apt-get update -y && apt-get install -y wget build-essential python3 python-is-python3 python3-distutils

RUN wget https://github.com/bazelbuild/bazelisk/releases/download/v1.7.2/bazelisk-linux-arm64 -O /usr/local/bin/bazel && chmod +x /usr/local/bin/bazel

# Force Bazelisk to pull down the latest Bazel version
RUN bazel version

COPY bazelrc /etc/bazel.bazelrc


