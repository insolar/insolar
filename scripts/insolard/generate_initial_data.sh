#!/usr/bin/env bash
set -e

# configurable vars
INSOLAR_ARTIFACTS_DIR=${INSOLAR_ARTIFACTS_DIR:-".artifacts"}/
LAUNCHNET_BASE_DIR=${LAUNCHNET_BASE_DIR:-"${INSOLAR_ARTIFACTS_DIR}/launchnet"}/

# dependent vars
CONFIGS_DIR=${LAUNCHNET_BASE_DIR}configs/

echo "generate members keys in dir: $CONFIGS_DIR"
bin/insolar gen-key-pair --target=user > ${CONFIGS_DIR}/root_member_keys.json
