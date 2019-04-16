#! /usr/bin/env bash
set -e
# Changeable environment variables (parameters)
# LOG_SYMLINKS_BY_NODE_REF activates log_symlinks_by_node_refs function
LOG_SYMLINKS_BY_NODE_REF=
INSOLAR_LOG_FORMATTER=${INSOLAR_LOG_FORMATTER:-"text"}
INSOLAR_LOG_LEVEL=${INSOLAR_LOG_LEVEL:-"debug"}
GORUND_LOG_LEVEL=${GORUND_LOG_LEVEL:-${INSOLAR_LOG_LEVEL}}

# predefined environment variables
BIN_DIR=bin
TEST_DATA=testdata
INSOLARD=$BIN_DIR/insolard
INSGORUND=$BIN_DIR/insgorund
PULSARD=$BIN_DIR/pulsard
PULSEWATCHER=$BIN_DIR/pulsewatcher
CONTRACT_STORAGE=contractstorage
LEDGER_DIR=data
LEDGER_NEW_DIR=new-data
CONFIGS_DIR=configs
BASE_DIR=scripts/insolard
KEYS_FILE=$BASE_DIR/$CONFIGS_DIR/bootstrap_keys.json
ROOT_MEMBER_KEYS_FILE=$BASE_DIR/$CONFIGS_DIR/root_member_keys.json
DISCOVERY_NODES_DATA=$BASE_DIR/discoverynodes/
NODES_DATA=$BASE_DIR/nodes/
INSGORUND_DATA=$BASE_DIR/insgorund/
GENESIS_CONFIG=$BASE_DIR/genesis.yaml
GENERATED_CONFIGS_DIR=$BASE_DIR/$CONFIGS_DIR/generated_configs
GENERATED_DISCOVERY_CONFIGS_DIR=$GENERATED_CONFIGS_DIR/discoverynodes
INSGORUND_PORT_FILE=$BASE_DIR/$CONFIGS_DIR/insgorund_ports.txt
PULSEWATCHER_CONFIG=$GENERATED_CONFIGS_DIR/utils/pulsewatcher.yaml


export INSOLAR_LOG_FORMATTER
export INSOLAR_LOG_LEVEL

NUM_DISCOVERY_NODES=$(sed '/^nodes:/ q' $GENESIS_CONFIG | grep "host:" | grep -v "#" | wc -l)

for i in `seq 1 $NUM_DISCOVERY_NODES`
do
    DISCOVERY_NODES+=($DISCOVERY_NODES_DATA/$i)
done

DISCOVERY_NODES_KEYS_DIR=$TEST_DATA/scripts/discovery_nodes

NUM_NODES=$(sed -n '/^nodes:/,$p' $GENESIS_CONFIG | grep "host:" | grep -v "#" | wc -l)

if [[ "$NUM_NODES" -ne "0" ]]
then
    for i in `seq 1 $NUM_NODES`
    do
        NODES+=($NODES_DATA/$i)
    done
fi

kill_port()
{
    port=$1
    pids=$(lsof -i :$port | grep "LISTEN\|UDP" | awk '{print $2}')
    for pid in $pids
    do
        echo "killing pid $pid"
        kill -9 $pid
    done
}

stop_listening()
{
    echo "stop_listening() starts ..."
    stop_insgorund=$1
    ports="$ports 58090" # Pulsar
    ports="$ports 53837" # Genesis
    if [[ "$stop_insgorund" == "true" ]]
    then
        gorund_ports=
        while read -r line; do

            listen_port=$( echo "$line" | awk '{print $1}' )
            rpc_port=$( echo "$line" | awk '{print $2}' )

            gorund_ports="$gorund_ports $listen_port $rpc_port"

        done < "$INSGORUND_PORT_FILE"

        gorund_ports="$gorund_ports $(echo $(pgrep insgorund ))"

        ports="$ports $gorund_ports"

    fi

    transport_ports=$( grep "host:" $GENESIS_CONFIG | grep -o ":\d\+" | grep -o "\d\+" | tr '\n' ' ' )
    ports="$ports $transport_ports"

    echo "Stop listening..."

    for port in $ports
    do
        echo "port: $port"
        kill_port $port
    done
    echo "stop_listening() end."
}

clear_dirs()
{
    echo "clear_dirs() starts ..."
    set -x
    rm -rfv $CONTRACT_STORAGE/*
    rm -rfv $LEDGER_DIR/*
    rm -rfv $LEDGER_NEW_DIR/*
    rm -rfv $DISCOVERY_NODES_DATA/*
    rm -rfv $GENERATED_CONFIGS_DIR/*
    rm -rfv $INSGORUND_DATA/*
    rm -rfv $NODES_DATA/*
    set +x
}

create_required_dirs()
{
    echo "create_required_dirs() starts ..."
    mkdir -vp $CONTRACT_STORAGE
    mkdir -vp $LEDGER_DIR
    mkdir -vp $LEDGER_NEW_DIR
    mkdir -vp $DISCOVERY_NODES_DATA/certs
    mkdir -vp $GENERATED_CONFIGS_DIR
    mkdir -vp $INSGORUND_DATA
    touch $INSGORUND_PORT_FILE

    for node in "${DISCOVERY_NODES[@]}"
    do
        mkdir -vp $node/data
        mkdir -vp $node/new-data
    done

    mkdir -p scripts/insolard/$CONFIGS_DIR

    echo "create_required_dirs() end."
}

generate_insolard_configs()
{
    echo "generate configs"
    go run scripts/generate_insolar_configs.go -o $GENERATED_CONFIGS_DIR -p $INSGORUND_PORT_FILE -g $GENESIS_CONFIG -t $BASE_DIR/pulsar_template.yaml
}

prepare()
{
    echo "prepare() starts ..."
    if [[ "$NO_STOP_LISTENING_ON_PREPARE" != "1" ]]; then
        stop_listening $run_insgorund
    fi
    clear_dirs
    create_required_dirs
    echo "prepare() end."
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
    echo "generate bootstrap keys"
    bin/insolar -c gen_keys > $KEYS_FILE
}

generate_root_member_keys()
{
    echo "generate root member_keys"
    bin/insolar -c gen_keys > $ROOT_MEMBER_KEYS_FILE
}

check_working_dir()
{
    echo "check_working_dir() starts ..."
    if ! pwd | grep -q "src/github.com/insolar/insolar$"
    then
        echo "Run me from insolar root"
        exit 1
    fi
    echo "check_working_dir() end."
}

usage()
{
    echo "usage: $0 [options]"
    echo "possible options: "
    echo -e "\t-h - show help"
    echo -e "\t-n - don't run insgorund"
    echo -e "\t-g - generate genesis"
    echo -e "\t-G - generate genesis and exit, show generation log"
    echo -e "\t-l - clear all and exit"
    echo -e "\t-C - generate configs only"
}

process_input_params()
{
    OPTIND=1
    while getopts "h?ngGlwC" opt; do
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
        G)
            NO_GENESIS_LOG_REDIRECT=1
            NO_STOP_LISTENING_ON_PREPARE=${NO_STOP_LISTENING_ON_PREPARE:-"1"}
            SKIP_BUILD=${SKIP_BUILD:-"1"}
            genesis
            exit 0
            ;;
        l)
            prepare
            exit 0
            ;;
        w)
            watch_pulse=false
            ;;
        C)
            generate_insolard_configs
            exit $?
        esac
    done
}

launch_insgorund()
{
    host=127.0.0.1
    metrics_port=28223
    while read -r line; do

        metrics_port=$((metrics_port + 20))
        listen_port=$( echo "$line" | awk '{print $1}' )
        rpc_port=$( echo "$line" | awk '{print $2}' )

        $INSGORUND -l $host:$listen_port --rpc $host:$rpc_port --log-level=$GORUND_LOG_LEVEL --metrics :$metrics_port &> $INSGORUND_DATA/$rpc_port.log &

    done < "$INSGORUND_PORT_FILE"
}

copy_data()
{
    echo "copy_data() starts ..."
    for node in "${DISCOVERY_NODES[@]}"
    do
        cp -v $LEDGER_DIR/* $node/data
        cp -v $LEDGER_NEW_DIR/* $node/new-data
    done
    echo "copy_data() end."
}

copy_discovery_certs()
{
    echo "copy_certs() starts ..."
    i=0
    for node in "${DISCOVERY_NODES[@]}"
    do
        i=$((i + 1))
        cp -v $DISCOVERY_NODES_DATA/certs/discovery_cert_$i.json $node/cert.json
    done
    echo "copy_certs() end."
}

wait_for_complete_network_state()
{
    while true
    do
        num=`scripts/insolard/check_status.sh 2>/dev/null | grep "NodeReady" | wc -l`
        echo "$num/$NUM_DISCOVERY_NODES discovery nodes ready"
        if [[ "$num" -eq "$NUM_DISCOVERY_NODES" ]]
        then
            break
        fi
        sleep 5s
    done
}

genesis()
{
    prepare
    if [[ "$SKIP_BUILD" != "1" ]]; then
        echo "build binaries"
        build_binaries
    else
        echo "SKIP: build binaries"
    fi
    generate_bootstrap_keys
    generate_root_member_keys
    generate_insolard_configs

    printf "start genesis ... \n"
    CMD="$INSOLARD --config $BASE_DIR/insolar.yaml --genesis $GENESIS_CONFIG --keyout $DISCOVERY_NODES_DATA/certs"
    if [[ "$NO_GENESIS_LOG_REDIRECT" != "1" ]]; then
        ${CMD} &> $DISCOVERY_NODES_DATA/genesis_output.log
    else
        ${CMD}
    fi
    printf "genesis is done\n"

    copy_data
    copy_discovery_certs

    if [[ "$LOG_SYMLINKS_BY_NODE_REF" != "" ]]; then
        log_symlinks_by_node_refs
    fi
}

log_symlinks_by_node_refs() {
    if which jq ; then
        NL=$BASE_DIR/loglinks
        mkdir  $NL || \
        rm -f $NL/*.log
        for node in "${DISCOVERY_NODES[@]}" ; do
            ref=`jq -r '.reference' $node/cert.json`
            [[ $ref =~ .+\. ]]
            ln -s `pwd`/$node/output.log $NL/${BASH_REMATCH[0]}log
        done
    else
        echo "no jq =("
    fi
}

run_insgorund=true
watch_pulse=true
check_working_dir
process_input_params $@

trap 'stop_listening true' INT TERM EXIT

printf "start pulsar ... \n"
ARTIFACTS_DIR=${ARTIFACTS_DIR:-".artifacts"}
PULSAR_DATA_DIR=$ARTIFACTS_DIR/pulsar_data
mkdir -p $PULSAR_DATA_DIR
$PULSARD -c $GENERATED_CONFIGS_DIR/pulsar.yaml --trace &> $DISCOVERY_NODES_DATA/pulsar_output.log &

if [[ "$run_insgorund" == "true" ]]
then
    printf "start insgorund ... \n"
    launch_insgorund
else
    echo "INSGORUND IS NOT LAUNCHED"
fi

printf "start discovery nodes ... \n"

i=0
for node in "${DISCOVERY_NODES[@]}"
do
    i=$((i + 1))
    INSOLAR_LOG_LEVEL=$insolar_log_level $INSOLARD --config $GENERATED_DISCOVERY_CONFIGS_DIR/insolar_$i.yaml --trace &> $node/output.log &
    echo "DISCOVERY NODE $i STARTED in background"
done

printf "discovery nodes started ... \n"

if [[ "$NUM_NODES" -ne "0" ]]
then
    wait_for_complete_network_state
    scripts/insolard/start_nodes.sh
fi

if [[ "$watch_pulse" == "true" ]]
then
    echo "Starting pulse watcher..."
    $PULSEWATCHER -c $PULSEWATCHER_CONFIG
else
    echo "Waiting..."
    wait
fi

echo "FINISHING ..."
