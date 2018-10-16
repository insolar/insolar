#!/usr/bin/env bash
set -e
source ./scripts/insolard/.env

stop_listening()
{
    ports="19191 $INSGORUND_LISTEN_PORT $INSGORUND_RPS_PORT 8090 8080"
    echo "Stop listening..."
    for port in $ports
    do
        lsof -i :$port | grep LISTEN | awk '{print $2}' | xargs kill
    done
}

clear_dirs()
{
    echo "Cleaning directories ..."
    rm -rfv $CONTRACT_STORAGE/*
    rm -rfv $LEDGER_DIR/*
}

create_required_dirs()
{
    mkdir -p $TEST_DATA/functional
    mkdir -p $CONTRACT_STORAGE
    mkdir -p $LEDGER_DIR
}

prepare()
{
    stop_listening
    create_required_dirs
    clear_dirs
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
	${BIN_DIR}/insolar -c gen_keys > scripts/insolard/bootstrap_keys.json
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
    $INSGORUND -l $host:$INSGORUND_LISTEN_PORT --rpc $host:$INSGORUND_RPS_PORT
}

trap stop_listening EXIT

gorund_only=0
param=$1
check_working_dir
process_input_params $param

prepare
build_binaries
generate_bootstrap_keys

if [ "$gorund_only" == "1" ]
then
    run_insgorund
else
    run_insgorund &
    $INSOLARD --config scripts/insolard/insolar.yaml
fi
