#! /usr/bin/env bash

set -euo pipefail

OS="$(go env GOHOSTOS)"
ARCH="$(go env GOARCH)"

cd $(git rev-parse --show-toplevel)

echo -e ">>> Compiling the Go proto"
for label in $(bazel query 'kind(go_proto_library, //...)'); do
  echo -e ">>> Found go_proto_library ${label}"
	package="${label%%:*}"
	package="${package##//}"
	target="${label##*:}"

	# do not continue if the package does not exist
	[[ -d "${package}" ]] || continue

	# compute the path where bazel put the files
	out_path="bazel-bin/${package}/${target}_/github.com/mjm/pi-tools/${package}"

	# compute the relative_path to the
	count_paths="$(echo -n "${package}" | tr '/' '\n' | wc -l)"
	relative_path=""
	for i in $(seq 0 ${count_paths}); do
		relative_path="../${relative_path}"
	done

  echo -e ">>> Building ${label}"
	bazel build "${label}"

	found=0
	for f in ${out_path}/*.pb.go; do
		if [[ -f "${f}" ]]; then
			found=1
			ln -nsf "${relative_path}${f}" "${package}/"
		fi
	done
	if [[ "${found}" == "0" ]]; then
		echo "ERR: no .pb.go file was found inside $out_path for the package ${package}"
		exit 1
	fi
done
