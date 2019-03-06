#!/usr/bin/env bash

BIN_DIR=bin
BASE_DIR=scripts/insolard
NODES_DATA=$BASE_DIR/nodes
GENESIS_CONFIG=$BASE_DIR/genesis.yaml
INSOLARD=$BIN_DIR/insolard
CONFIGS_DIR=configs
ROOT_MEMBER_KEYS_FILE=$BASE_DIR/$CONFIGS_DIR/root_member_keys.json
GENERATED_CONFIGS_DIR=$BASE_DIR/$CONFIGS_DIR/generated_configs/nodes
CERT_GENERATOR=$BIN_DIR/certgen

insolar_log_level=Debug

NUM_NODES=$(sed -n '/^nodes:/,$p' $GENESIS_CONFIG | grep "host:" | grep -cv "#" )
ROLES=($(sed -n '/^nodes:/,$p' ./scripts/insolard/genesis.yaml | grep "role" | cut -d: -f2))

for i in `seq 1 $NUM_NODES`
do
    NODES+=($NODES_DATA/$i)
done

create_nodes_dir()
{
    for node in "${NODES[@]}"
    do
        mkdir -vp $node/data
    done
}

clear_dirs()
{
    rm -rfv $NODES_DATA/*
}

generate_nodes_certs()
{
    echo "generate_nodes_certs() starts ..."
    mkdir $NODES_DATA/certs/
    i=0
    for node in "${NODES[@]}"
    do
        role="${ROLES[$i]//\"}"
        i=$((i + 1))
        $CERT_GENERATOR --root-conf $ROOT_MEMBER_KEYS_FILE -h "http://127.0.0.1:19101/api" -c $NODES_DATA/certs/node_cert_$i.json -k $node/keys.json -r $role
        cp -v $NODES_DATA/certs/node_cert_$i.json $node/cert.json
    done
    echo "generate_nodes_certs() end."
}

printf "start nodes ... \n"

clear_dirs
create_nodes_dir
generate_nodes_certs

i=0
for node in "${NODES[@]}"
do
    i=$((i + 1))
    INSOLAR_LOG_LEVEL=$insolar_log_level $INSOLARD --config $GENERATED_CONFIGS_DIR/insolar_$i.yaml --trace &> $node/output.log &
    echo "NODE $i STARTED in background"
done

printf "nodes started ... \n"
