#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" # Directory where this script exists.
__root="$(cd "$(dirname "${__dir}")" && pwd)"         # Root directory of project.

export GITHUB_TOKEN=""
VERSION_NAME=${BRING_VERSION:-"v0.0.0-test"}

cd "$__root"

docker buildx build \
	--no-cache \
	--progress=plain \
	--secret id=github_token,env=GITHUB_TOKEN \
	--build-arg "VERSION_NAME=$VERSION_NAME" \
	--tag bring:test \
	.
