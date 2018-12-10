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
NODES_DATA=scripts/insolard/nodes/

NODE_1=$NODES_DATA/1
NODE_2=$NODES_DATA/2
NODE_3=$NODES_DATA/3
NODE_4=$NODES_DATA/4
NODE_5=$NODES_DATA/5


DISCOVERY_NODES_KEYS_DIR=$TEST_DATA/scripts/discovery_nodes

stop_listening()
{
    stop_insgorund=$1
    ports="13831 13832 23832 23833 33833 33834"
    if [ "$stop_insgorund" == "true" ]
    then
        ports="$ports $INSGORUND_LISTEN_PORT $INSGORUND_RPS_PORT"
    fi
    
    echo "Stop listening..."
    for port in $ports
    do
        echo "port: $port"
        pids=$(lsof -i :$port | grep "LISTEN\|UDP" | awk '{print $2}')
        for pid in $pids
        do
            echo "killing pid $pid"
            kill -9 $pid
        done
    done
}

clear_dirs()
{
    echo "Cleaning directories ... "
    rm -rfv $CONTRACT_STORAGE/*
    rm -rfv $LEDGER_DIR/*
    rm -rfv $NODES_DATA/*
}

create_required_dirs()
{
    mkdir -p $TEST_DATA/functional
    mkdir -p $CONTRACT_STORAGE
    mkdir -p $LEDGER_DIR
    mkdir -p $NODES_DATA/certs

    mkdir -p $NODE_1/data
    mkdir -p $NODE_2/data
    mkdir -p $NODE_3/data
    mkdir -p $NODE_4/data
    mkdir -p $NODE_5/data

    mkdir -p scripts/insolard/$CONFIGS_DIR
}

prepare()
{
    stop_listening $run_insgorund
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
    bin/insolar -c gen_keys > $NODE_1/keys.json
    bin/insolar -c gen_keys > $NODE_2/keys.json
    bin/insolar -c gen_keys > $NODE_3/keys.json
    bin/insolar -c gen_keys > $NODE_4/keys.json
    bin/insolar -c gen_keys > $NODE_5/keys.json
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
    echo "usage: $0 [options]"
    echo "possible options: "
    echo -e "\t-h - show help"
    echo -e "\t-n - don't run insgorund"
    echo -e "\t-g - preventively generate initial ledger"
    echo -e "\t-l - clear all and exit"
}

process_input_params()
{
    OPTIND=1
    while getopts "h?ngl" opt; do
        case "$opt" in
        h|\?)
            usage
            exit 0
            ;;
        n)
            run_insgorund=false
            ;;
        g)
            genesis
            ;;
        l)
            prepare
            exit 0
            ;;
        esac
    done
}

launch_insgorund()
{
    host=127.0.0.1
    $INSGORUND -l $host:$INSGORUND_LISTEN_PORT --rpc $host:$INSGORUND_RPS_PORT
}

copy_data()
{
    cp $LEDGER_DIR/* $NODE_1/data
    cp $LEDGER_DIR/* $NODE_2/data
    cp $LEDGER_DIR/* $NODE_3/data
    cp $LEDGER_DIR/* $NODE_4/data
    cp $LEDGER_DIR/* $NODE_5/data
}

copy_certs()
{
    cp $NODES_DATA/certs/discovery_cert_1.json $NODE_1/cert.json
    cp $NODES_DATA/certs/discovery_cert_2.json $NODE_2/cert.json
    cp $NODES_DATA/certs/discovery_cert_3.json $NODE_3/cert.json
#    cp $NODES_DATA/certs/discovery_cert_4.json $NODE_4/cert.json
#    cp $NODES_DATA/certs/discovery_cert_5.json $NODE_5/cert.json
}

genesis()
{
    prepare
    build_binaries
    generate_bootstrap_keys
    generate_root_member_keys
    generate_discovery_nodes_keys

    printf "start genesis ... \n"
    $INSOLARD --config scripts/insolard/insolar.yaml --genesis scripts/insolard/genesis.yaml --keyout $NODES_DATA/certs
    printf "genesis is done\n"

    copy_data
    copy_certs
}

trap 'stop_listening true' EXIT

run_insgorund=true
check_working_dir
process_input_params $@

printf "start pulsar ... \n"
$PULSARD -c scripts/insolard/pulsar.yaml &> $NODES_DATA/pulsar_output.txt &

if [ "$run_insgorund" == "true" ]
then
    printf "start insgorund ... \n"
    launch_insgorund &
else
    echo "INSGORUND IS NOT LAUNCHED"
fi

printf "start nodes ... \n"

$INSOLARD --config scripts/insolard/insolar_1.yaml &> $NODE_1/output.txt &
echo "STARTED 1"
$INSOLARD --config scripts/insolard/insolar_2.yaml &> $NODE_2/output.txt &
echo "STARTED 2"
$INSOLARD --config scripts/insolard/insolar_3.yaml &> $NODE_3/output.txt
#echo "STARTED 3"
#$INSOLARD --config scripts/insolard/insolar_4.yaml &> $NODE_4/output.txt &
#echo "STARTED 4"
#$INSOLARD --config scripts/insolard/insolar_5.yaml &> $NODE_5/output.txt
#echo "FINISHING ..."
