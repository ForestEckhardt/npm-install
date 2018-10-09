#!/usr/bin/env bash
set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."
source  ./scripts/install_tools.sh

# Install Pack CLI
host=$([ $(uname -s) == 'Darwin' ] &&  printf "macos" || printf "linux")
version=$(curl --silent "https://api.github.com/repos/buildpack/pack/releases/latest" | jq -r .tag_name)
wget "https://github.com/buildpack/pack/releases/download/$version/pack-$host.tar.gz" -O $GOBIN/pack && chmod +x $GOBIN/pack

GINKGO_NODES=${GINKGO_NODES:-1}
GINKGO_ATTEMPTS=${GINKGO_ATTEMPTS:-1}

export CNB_BUILD_IMAGE=${CNB_BUILD_IMAGE:-cfbuildpacks/cflinuxfs3-cnb-experimental:build}

# TODO: change default to `cfbuildpacks/cflinuxfs3-cnb-experimental:run` when pack cli can use it
export CNB_RUN_IMAGE=${CNB_RUN_IMAGE:-packs/run}

# Always pull latest images
# Most helpful for local testing consistency with CI (which would already pull the latest)
docker pull $CNB_BUILD_IMAGE
docker pull $CNB_RUN_IMAGE

cd integration

echo "Run Buildpack Runtime Integration Tests"
ginkgo -r --flakeAttempts=$GINKGO_ATTEMPTS -nodes $GINKGO_NODES
