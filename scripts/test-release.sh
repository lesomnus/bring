#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" # Directory where this script exists.
__root="$(cd "$(dirname "${__dir}")" && pwd)"         # Root directory of project.

cd "$__root"
./scripts/gen-version-file.sh > /dev/null

TMP=$(mktemp -d -t bring-XXXX)
echo "$TMP"

go build -o "$TMP/bring"
cd "$TMP"

./bring version > /dev/null
export BRING="$TMP/bring"

export TARGETARCH="$(dpkg --print-architecture)"
export REPO_NAME="lesomnus/bring"
export VERSION_NAME=""

if [ -z "${GH_TOKEN:-""}" ]; then
	export GH_TOKEN="gho_invalid"
fi

"$__root/scripts/release.sh"

rm -rf "$TMP"


