#!/usr/bin/env bash
set -e
INSOLAR_BIN=${INSOLAR_BIN:-"insolar"}
BOOTSTRAP_CONFIG=${BOOTSTRAP_CONFIG:-"/etc/bootstrap/bootstrap.yaml"}
CONFIGS_DIR=${CONFIGS_DIR:-"/var/data/bootstrap/configs/"}

PULSAR_KEYS=${CONFIGS_DIR}pulsar-keys.json

generate_pulsar_keys()
{
    echo "generate pulsar keys: ${PULSAR_KEYS}"
    ${INSOLAR_BIN} gen-key-pair --target=node > ${PULSAR_KEYS}
}

generate_root_member_keys()
{
    echo "generate members keys in dir: $CONFIGS_DIR"
    ${INSOLAR_BIN} gen-key-pair --target=user > ${CONFIGS_DIR}root_member_keys.json

    ${INSOLAR_BIN} gen-key-pair --target=user > ${CONFIGS_DIR}fee_member_keys.json
    ${INSOLAR_BIN} gen-key-pair --target=user > ${CONFIGS_DIR}migration_admin_member_keys.json
    for (( b = 0; b < 10; b++ ))
    do
        ${INSOLAR_BIN} gen-key-pair --target=node > ${CONFIGS_DIR}migration_daemon_${b}_member_keys.json
    done

    for (( b = 0; b < 140; b++ ))
    do
        ${INSOLAR_BIN} gen-key-pair --target=user > ${CONFIGS_DIR}network_incentives_${b}_member_keys.json
    done

    for (( b = 0; b < 40; b++ ))
    do
        ${INSOLAR_BIN} gen-key-pair --target=user > ${CONFIGS_DIR}application_incentives_${b}_member_keys.json
    done

    for (( b = 0; b < 40; b++ ))
    do
        ${INSOLAR_BIN} gen-key-pair --target=user > ${CONFIGS_DIR}foundation_${b}_member_keys.json
    done

    for (( b = 0; b < 1; b++ ))
    do
        ${INSOLAR_BIN} gen-key-pair --target=user > ${CONFIGS_DIR}funds_${b}_member_keys.json
    done

    for (( b = 0; b < 8; b++ ))
    do
        ${INSOLAR_BIN} gen-key-pair --target=user > ${CONFIGS_DIR}enterprise_${b}_member_keys.json
    done
}

generate_migration_addresses()
{
    echo "generate migration addresses: ${CONFIGS_DIR}migration_addresses.json"
    ${INSOLAR_BIN} gen-migration-addresses -c 10000 > ${CONFIGS_DIR}migration_addresses.json
}

bootstrap()
{
    echo "bootstrap start"
    generate_pulsar_keys
    generate_root_member_keys
    generate_migration_addresses
}

mkdir -p ${CONFIGS_DIR}
bootstrap

$INSOLAR_BIN bootstrap -c $BOOTSTRAP_CONFIG

if [[ "$HEAVY_GENESIS" ]]; then
    insolard -c /app/bootstrap/insolard.yaml --heavy-genesis=/var/data/bootstrap/heavy_genesis.json --genesis-only
fi

MY_BIN_DIR=$( dirname "${BASH_SOURCE[0]}" )
cp ${MY_BIN_DIR}/kustomization.yaml ${CONFIGS_DIR}/../
