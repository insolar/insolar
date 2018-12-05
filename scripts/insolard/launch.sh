#!/usr/bin/env bash
set -e

BIN_DIR=bin
TEST_DATA=testdata
INSOLARD=$BIN_DIR/insolard
INSGORUND=$BIN_DIR/insgorund
CONTRACT_STORAGE=contractstorage
LEDGER_DIR=data
INSGORUND_LISTEN_PORT=18181
INSGORUND_RPS_PORT=18182
INSGOLAR_METRICS_PORT=9090
INSGORUND_METRICS_PORT=9091
CONFIGS_DIR=configs
ROOT_MEMBER_KEYS_FILE=scripts/insolard/$CONFIGS_DIR/root_member_keys.json

stop_listening()
{
    ports="19191 $INSGORUND_LISTEN_PORT $INSGORUND_RPS_PORT 8090 8080 $INSGOLAR_METRICS_PORT $INSGORUND_METRICS_PORT"
    if [ "$1" != "" ]
    then
        ports=$@
    fi
    echo "Stop listening..."
    for port in $ports
    do
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
    mkdir -p scripts/insolard/$CONFIGS_DIR
    mkdir -p scripts/insolard/$CONFIGS_DIR/certs
}

prepare()
{
    if [ "$gorund_only" == "1" ]
    then
        stop_listening $INSGORUND_LISTEN_PORT $INSGORUND_RPS_PORT
    else
        stop_listening
    fi
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

generate_root_member_keys()
{
	bin/insolar -c gen_keys > $ROOT_MEMBER_KEYS_FILE
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
    echo "usage: $0 <clear|rebuild|gorund_only>"
}

process_input_params()
{
    param=$1
    if [  "$param" == "clear" ]
    then
        prepare
        exit 0
    fi

    if [ "$param" == "rebuild" ]
    then
        rebuild_binaries
        exit 0
    fi

    if [ "$param" == "gorund_only" ]
    then
        gorund_only=1
        return 0
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
    $INSGORUND -l $host:$INSGORUND_LISTEN_PORT --rpc $host:$INSGORUND_RPS_PORT --metrics $host:$INSGORUND_METRICS_PORT
}

trap stop_listening EXIT

gorund_only=0
param=$1
check_working_dir
process_input_params $param

prepare
build_binaries
generate_root_member_keys

if [ "$gorund_only" == "1" ]
then
    run_insgorund
else
    run_insgorund &
    $INSOLARD --config scripts/insolard/one_node_insolar.yaml --genesis scripts/insolard/one_node_genesis.yaml --keyout scripts/insolard/$CONFIGS_DIR/certs
    $INSOLARD --config scripts/insolard/one_node_insolar.yaml
fi
