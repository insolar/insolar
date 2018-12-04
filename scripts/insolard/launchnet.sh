#!/usr/bin/env bash
set -e

BIN_DIR=bin
TEST_DATA=testdata
INSOLARD=$BIN_DIR/insolard
INSGORUND=$BIN_DIR/insgorund
PULSARD=$BIN_DIR/pulsard
CONTRACT_STORAGE=contractstorage
LEDGER_DIR=data
INSGORUND_LISTEN_PORT=18181
INSGORUND_RPS_PORT=18182
CONFIGS_DIR=configs
KEYS_FILE=scripts/insolard/$CONFIGS_DIR/bootstrap_keys.json
ROOT_MEMBER_KEYS_FILE=scripts/insolard/$CONFIGS_DIR/root_member_keys.json
CERTIFICATE_FILE=scripts/insolard/$CONFIGS_DIR/certificate.json
NODES_DATA=scripts/insolard/nodes/
FIRST_NODE=$NODES_DATA/first
SECOND_NODE=$NODES_DATA/second
THIRD_NODE=$NODES_DATA/third

DISCOVERY_NODES_KEYS_DIR=$TEST_DATA/scripts/discovery_nodes

stop_listening()
{
    ports="$INSGORUND_LISTEN_PORT $INSGORUND_RPS_PORT 53835 53837 53839 38181"
    if [ "$1" != "" ]
    then
        ports=$@
    fi
    echo "Stop listening..."
    for port in $ports
    do
        echo "port: $port"
        lsof -i :$port | grep LISTEN | awk '{print $2}' | xargs kill
    done
}

clear_dirs()
{
    echo "Cleaning directories ... "
    rm -rfv $CONTRACT_STORAGE/*
    rm -rfv $LEDGER_DIR/*
}

create_required_dirs()
{
    mkdir -p $TEST_DATA/functional
    mkdir -p $CONTRACT_STORAGE
    mkdir -p $LEDGER_DIR
    mkdir -p $NODES_DATA
    mkdir -p $NODES_DATA/certs
    mkdir -p $FIRST_NODE
    mkdir -p $FIRST_NODE/data
    mkdir -p $SECOND_NODE
    mkdir -p $SECOND_NODE/data
    mkdir -p $THIRD_NODE
    mkdir -p $THIRD_NODE/data
    mkdir -p scripts/insolard/$CONFIGS_DIR
}

prepare()
{
    stop_listening
    clear_dirs
    create_required_dirs
}

build_binaries()
{
    make build
}

rebuild_binaries()
{
    make clean
    build_binaries
}

generate_bootstrap_keys()
{
	bin/insolar -c gen_keys > $KEYS_FILE
}

generate_root_member_keys()
{
	bin/insolar -c gen_keys > $ROOT_MEMBER_KEYS_FILE
}

generate_discovery_nodes_keys()
{
    bin/insolar -c gen_keys > $FIRST_NODE/keys.json
    bin/insolar -c gen_keys > $SECOND_NODE/keys.json
    bin/insolar -c gen_keys > $THIRD_NODE/keys.json
}

generate_certificate()
{
    bin/insolar -c gen_certificate -g $KEYS_FILE > $CERTIFICATE_FILE
}

check_working_dir()
{
    if ! pwd | grep -q "src/github.com/insolar/insolar$"
    then
        echo "Run me from insolar root"
        exit 1
    fi
}

usage()
{
    echo "usage: $0 <clear>"
}

process_input_params()
{
    param=$1
    if [  "$param" == "clear" ]
    then
        prepare
        exit 0
    fi

    if [ "$param" == "help" ] || [ "$param" == "-h" ] || [ "$param" == "--help" ]
    then
        usage
        exit 0
    fi
}

run_insgorund()
{
    host=127.0.0.1
    $INSGORUND -l $host:$INSGORUND_LISTEN_PORT --rpc $host:$INSGORUND_RPS_PORT
}

copy_data()
{
    cp $LEDGER_DIR/* $FIRST_NODE/data
    cp $LEDGER_DIR/* $SECOND_NODE/data
    cp $LEDGER_DIR/* $THIRD_NODE/data
}

copy_serts()
{
    cp $NODES_DATA/certs/discovery_cert_1.json $FIRST_NODE/cert.json
    cp $NODES_DATA/certs/discovery_cert_2.json $SECOND_NODE/cert.json
    cp $NODES_DATA/certs/discovery_cert_3.json $THIRD_NODE/cert.json
}

trap stop_listening EXIT

param=$1
check_working_dir
process_input_params $param

prepare
build_binaries
generate_bootstrap_keys
generate_root_member_keys
generate_certificate
generate_discovery_nodes_keys

printf "start pulsar ... \n"
$PULSARD -c scripts/insolard/pulsar.yaml &> $NODES_DATA/pulsar_output.txt &

printf "start insgorund ... \n"
run_insgorund &

printf "start genesis ... \n"
$INSOLARD --config scripts/insolard/insolar.yaml --genesis scripts/insolard/genesis.yaml --keyout $NODES_DATA/certs
printf "genesis is done\n"

copy_data
copy_serts

printf "start nodes ... \n"

$INSOLARD --config scripts/insolard/first_insolar.yaml &> $FIRST_NODE/output.txt &
echo "FIRST STARTED"
$INSOLARD --config scripts/insolard/second_insolar.yaml &> $SECOND_NODE/output.txt &
echo "SECOND STARTED"
$INSOLARD --config scripts/insolard/third_insolar.yaml &> $THIRD_NODE/output.txt
echo "FINISHING ..."
