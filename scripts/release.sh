#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

BRING="${BRING:-"./bring"}"
VERSION_NAME="${VERSION_NAME:-"v0.0.0-test"}"

HAS_DEPS=1
env_exists() {
	local name="$1"
    if [ -z "${!name:-}" ]; then
        echo "[!] env \"$1\" is not provided"
        HAS_DEPS=0
    fi
}
env_exists "TARGETARCH"
env_exists "REPO_NAME"
env_exists "GH_TOKEN"
if [ ! -x "$BRING" ]; then
	echo '[!] "bring" executable not found at' "\"$BRING\""
	HAS_DEPS=0
fi
if [ ! "$HAS_DEPS" = '1' ]; then
	echo 'requirements not met; abort'
	exit 1
fi

"$BRING" version | tee .version
source .version
echo "test " "$BRING_VERSION" "==" "$VERSION_NAME"
if [ ! "$BRING_VERSION" = "$VERSION_NAME" ]; then
	echo "expected version name: " "$VERSION_NAME"
	echo "version name not matched; abort"
	exit 1
fi

BRING_NAME="./bring-$VERSION_NAME-linux-$TARGETARCH"
rm "$BRING_NAME" &> /dev/null || true
ln -s "$BRING" "$BRING_NAME"

GH_URL="https://github.com/cli/cli/releases/download/v2.62.0/gh_2.62.0_linux_$TARGETARCH.tar.gz"
curl -fsSL "$GH_URL" -o ./gh.tar.gz
tar -xf ./gh.tar.gz
mv ./gh_*/bin/gh ./gh

./gh version
./gh auth status 
./gh release \
	--repo "$REPO_NAME" \
	upload "$VERSION_NAME" "$BRING_NAME" \
	--clobber
