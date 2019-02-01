#!/usr/bin/env bash

BIN_DIR=bin
BASE_DIR=scripts/insolard
NODES_DATA=$BASE_DIR/nodes
GENESIS_CONFIG=$BASE_DIR/genesis.yaml
INSOLARD=$BIN_DIR/insolard
CONFIGS_DIR=configs
ROOT_MEMBER_KEYS_FILE=$BASE_DIR/$CONFIGS_DIR/root_member_keys.json
GENERATED_CONFIGS_DIR=$BASE_DIR/$CONFIGS_DIR/generated_configs/nodes
CERT_GETERATOR=$BIN_DIR/certgen

insolar_log_level=Debug

NUM_NODES=$(sed -n '/^nodes:/,$p' $GENESIS_CONFIG | grep "host:" | grep -cv "#" )

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
        i=$((i + 1))
        $CERT_GETERATOR --root-conf $ROOT_MEMBER_KEYS_FILE -c $NODES_DATA/certs/node_cert_$i.json -k $node/keys.json
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
    if [[ "$i" -eq "$NUM_NODES" ]]
    then
        echo "NODE $i STARTED in foreground"
        INSOLAR_LOG_LEVEL=$insolar_log_level $INSOLARD --config $GENERATED_CONFIGS_DIR/insolar_$((i+NUM_DISCOVERY_NODES)).yaml --trace &> $node/output.txt
        lastNodePID=`echo \$!`
        break
    fi
    INSOLAR_LOG_LEVEL=$insolar_log_level $INSOLARD --config $GENERATED_CONFIGS_DIR/insolar_$((i+NUM_DISCOVERY_NODES)).yaml --trace &> $node/output.txt &
    echo "NODE $i STARTED in background"
done

sleep 10s   #time to consensus
printf "nodes started ... \n"
