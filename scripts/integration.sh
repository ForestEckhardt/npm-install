#!/usr/bin/env bash
set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."
source  ./scripts/install_tools.sh

GINKGO_NODES=${GINKGO_NODES:-1}
GINKGO_ATTEMPTS=${GINKGO_ATTEMPTS:-1}

export CNB_BUILD_IMAGE=${CNB_BUILD_IMAGE:-cfbuildpacks/cflinuxfs3-cnb-experimental:build}

# TODO: change default to `cfbuildpacks/cflinuxfs3-cnb-experimental:run` when pack cli can use it
export CNB_RUN_IMAGE=${CNB_RUN_IMAGE:-packs/run}
docker pull $CNB_RUN_IMAGE

cd integration

echo "Run Buildpack Runtime Integration Tests"
ginkgo -r --flakeAttempts=$GINKGO_ATTEMPTS -nodes $GINKGO_NODES
