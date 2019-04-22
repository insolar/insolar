#! /usr/bin/env bash
set -e
echo "***** start_nodes.sh start *****"

# configurable vars
INSOLAR_ARTIFACTS_DIR=${INSOLAR_ARTIFACTS_DIR:-".artifacts"}/
LAUNCHNET_BASE_DIR=${LAUNCHNET_BASE_DIR:-"${INSOLAR_ARTIFACTS_DIR}launchnet"}/

# dependent vars
LAUNCHNET_LOGS_DIR=${LAUNCHNET_BASE_DIR}logs/
NODES_LOGS=${LAUNCHNET_LOGS_DIR}nodes/

BIN_DIR=bin
INSOLAR_CMD=$BIN_DIR/insolar
INSOLARD_CMD=$BIN_DIR/insolard

NODES_DATA=${LAUNCHNET_BASE_DIR}nodes/

#CERT_DIR=${LAUNCHNET_BASE_DIR}generated_configs/nodes/

GENESIS_CONFIG=${LAUNCHNET_BASE_DIR}genesis.yaml
ROOT_MEMBER_KEYS_FILE=${LAUNCHNET_BASE_DIR}configs/root_member_keys.json
#GENERATED_CONFIGS_DIR=${LAUNCHNET_BASE_DIR}/configs/generated_configs/nodes

NUM_NODES=$(sed -n '/^nodes:/,$p' $GENESIS_CONFIG | grep "host:" | grep -cv "#" )
ROLES=($(sed -n '/^nodes:/,$p' ./scripts/insolard/genesis_template.yaml | grep "role" | cut -d: -f2))
(>&2 echo "ROLES=$ROLES")
(>&2 echo "NUM_NODES=$NUM_NODES")
#exit


create_nodes_dir()
{
    echo "prepare nodes dir"
    for i in `seq 1 $NUM_NODES`
    do
        set -x
        mkdir -vp ${NODES_LOGS}${i}
        { set +x; } 2>/dev/null
    done
}

clear_dirs()
{
    echo "clear nodes dirs"
    rm -rf $NODES_DATA/certs/

    for i in `seq 1 $NUM_NODES`
    do
        set -x
        rm -rvf ${NODES_LOGS}${i}
        { set +x; } 2>/dev/null
    done
}

generate_nodes_certs()
{
    echo "generate_nodes_certs() starts ..."
    mkdir -p $NODES_DATA/certs/
    for i in `seq 1 $NUM_NODES`
    do
        role="${ROLES[$i]//\"}"
        set -x
        ${INSOLAR_CMD} certgen \
            --root-keys ${ROOT_MEMBER_KEYS_FILE} \
            --url "http://127.0.0.1:19101/api" \
            --node-cert ${NODES_DATA}${i}/cert.json \
            --node-keys ${NODES_DATA}${i}/keys.json \
            --role ${role}
        { set +x; } 2>/dev/null
    done
    echo "generate_nodes_certs() end."
}

clear_dirs
create_nodes_dir
generate_nodes_certs

for i in `seq 1 $NUM_NODES`
do
    set -x
    $INSOLARD_CMD \
        --config ${NODES_DATA}${i}/insolard.yaml \
        --trace &> ${NODES_LOGS}${i}/output.log &
    { set +x; } 2>/dev/null
    echo "NODE $i STARTED in background"
done

echo "nodes started ..."
