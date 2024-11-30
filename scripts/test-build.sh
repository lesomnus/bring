#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" # Directory where this script exists.
__root="$(cd "$(dirname "${__dir}")" && pwd)"         # Root directory of project.

cd "$__root"

REPO_NAME="lesomnus/bring"
VERSION_NAME=""
if [ -z "${GH_TOKEN:-""}" ]; then
	export GH_TOKEN="gho_invalid"
fi

docker buildx build \
	--no-cache \
	--progress=plain \
	--secret id=github_token,env=GH_TOKEN \
	--build-arg "REPO_NAME=$REPO_NAME" \
	--build-arg "VERSION_NAME=$VERSION_NAME" \
	--tag bring:test \
	.
