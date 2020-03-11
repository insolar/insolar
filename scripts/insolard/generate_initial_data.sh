#!/usr/bin/env bash
set -e

# configurable vars
INSOLAR_ARTIFACTS_DIR=${INSOLAR_ARTIFACTS_DIR:-".artifacts"}/
LAUNCHNET_BASE_DIR=${LAUNCHNET_BASE_DIR:-"${INSOLAR_ARTIFACTS_DIR}/launchnet"}/

# dependent vars
CONFIGS_DIR=${LAUNCHNET_BASE_DIR}configs/

echo "generate members keys in dir: $CONFIGS_DIR"
bin/insolar gen-key-pair --target=user > ${CONFIGS_DIR}/root_member_keys.json

bin/insolar gen-key-pair --target=user > ${CONFIGS_DIR}/fee_member_keys.json
bin/insolar gen-key-pair --target=user > ${CONFIGS_DIR}/migration_admin_member_keys.json
for (( b = 0; b < 10; b++ ))
do
bin/insolar gen-key-pair --target=node > ${CONFIGS_DIR}/migration_daemon_${b}_member_keys.json
done

for (( b = 0; b < 140; b++ ))
do
bin/insolar gen-key-pair --target=user > ${CONFIGS_DIR}/network_incentives_${b}_member_keys.json
done

for (( b = 0; b < 40; b++ ))
do
bin/insolar gen-key-pair --target=user > ${CONFIGS_DIR}/application_incentives_${b}_member_keys.json
done

for (( b = 0; b < 40; b++ ))
do
bin/insolar gen-key-pair --target=user > ${CONFIGS_DIR}/foundation_${b}_member_keys.json
done

for (( b = 0; b < 1; b++ ))
do
bin/insolar gen-key-pair --target=user > ${CONFIGS_DIR}/funds_${b}_member_keys.json
done

for (( b = 0; b < 8; b++ ))
do
bin/insolar gen-key-pair --target=user > ${CONFIGS_DIR}/enterprise_${b}_member_keys.json
done

echo "generate migration addresses: ${CONFIGS_DIR}/migration_addresses.json"
bin/insolar gen-migration-addresses > ${CONFIGS_DIR}/migration_addresses.json
