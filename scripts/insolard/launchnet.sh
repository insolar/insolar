#!/usr/bin/env bash
set -e
# requires: lsof, awk, sed, grep, pgrep

# Changeable environment variables (parameters)
INSOLAR_ARTIFACTS_DIR=${INSOLAR_ARTIFACTS_DIR:-".artifacts"}/
LAUNCHNET_BASE_DIR=${LAUNCHNET_BASE_DIR:-"${INSOLAR_ARTIFACTS_DIR}launchnet"}/

INSOLAR_LOG_FORMATTER=${INSOLAR_LOG_FORMATTER:-"text"}
INSOLAR_LOG_LEVEL=${INSOLAR_LOG_LEVEL:-"debug"}
GORUND_LOG_LEVEL=${GORUND_LOG_LEVEL:-${INSOLAR_LOG_LEVEL}}

# predefined/dependent environment variables

LAUNCHNET_LOGS_DIR=${LAUNCHNET_BASE_DIR}logs/
DISCOVERY_NODE_LOGS=${LAUNCHNET_LOGS_DIR}discoverynodes/

BIN_DIR=bin
INSOLARD=$BIN_DIR/insolard
INSGORUND=$BIN_DIR/insgorund
PULSARD=$BIN_DIR/pulsard
PULSEWATCHER=$BIN_DIR/pulsewatcher

LEDGER_DIR=${LAUNCHNET_BASE_DIR}data
LEDGER_NEW_DIR=${LAUNCHNET_BASE_DIR}new-data

# TODO: move to launchnet dir
PULSAR_DATA_DIR=${LAUNCHNET_BASE_DIR}pulsar_data
PULSAR_CONFIG=${LAUNCHNET_BASE_DIR}pulsar.yaml

SCRIPTS_DIR=scripts/insolard/

CONFIGS_DIR=${LAUNCHNET_BASE_DIR}configs

KEYS_FILE=$CONFIGS_DIR/bootstrap_keys.json
ROOT_MEMBER_KEYS_FILE=$CONFIGS_DIR/root_member_keys.json

# TODO: use only heavy matereal data dir
DISCOVERY_NODES_DATA=${LAUNCHNET_BASE_DIR}discoverynodes/

DISCOVERY_NODES_HEAVY_DATA=${DISCOVERY_NODES_DATA}1/

NODES_DATA=${LAUNCHNET_BASE_DIR}nodes/
INSGORUND_DATA=${LAUNCHNET_BASE_DIR}insgorund/

GENESIS_TEMPLATE=${SCRIPTS_DIR}genesis_template.yaml
GENESIS_CONFIG=${LAUNCHNET_BASE_DIR}genesis.yaml
INSOLARD_CONFIG=${LAUNCHNET_BASE_DIR}insolard.yaml

PULSEWATCHER_CONFIG=${LAUNCHNET_BASE_DIR}/pulsewatcher.yaml

INSGORUND_PORT_FILE=$CONFIGS_DIR/insgorund_ports.txt

export INSOLAR_LOG_FORMATTER
export INSOLAR_LOG_LEVEL

NUM_DISCOVERY_NODES=$(sed '/^nodes:/ q' $GENESIS_TEMPLATE | grep "host:" | grep -v "#" | wc -l | tr -d '[:space:]')
NUM_NODES=$(sed -n '/^nodes:/,$p' $GENESIS_TEMPLATE | grep "host:" | grep -v "#" | wc -l | tr -d '[:space:]')
echo "discovery+other nodes: ${NUM_DISCOVERY_NODES}+${NUM_NODES}"

for i in `seq 1 $NUM_DISCOVERY_NODES`
do
    DISCOVERY_NODE_DIRS+=($DISCOVERY_NODES_DATA/$i)
done


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
#    set -x
    echo "stop_listening(): starts ..."
    stop_insgorund=$1
    ports="$ports 58090" # Pulsar
    if [[ "$stop_insgorund" == "true" ]]
    then
        echo "stop_listening(): stop insgorund"
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

    for port in $ports
    do
        echo "kill port '$port' owner"
        kill_port $port &
    done
    wait
    echo "stop_listening() end."
}

clear_dirs()
{
    echo "clear_dirs() starts ..."
    set -x
    rm -rfv ${LEDGER_DIR}
    rm -rfv ${DISCOVERY_NODES_DATA}
    rm -rfv ${INSGORUND_DATA}
    rm -rfv ${NODES_DATA}
    rm -rfv ${LAUNCHNET_LOGS_DIR}
    { set +x; } 2>/dev/null

    for i in `seq 1 $NUM_DISCOVERY_NODES`
    do
        set -x
        rm -rfv ${DISCOVERY_NODE_LOGS}${i}
        { set +x; } 2>/dev/null
    done
}

create_required_dirs()
{
    echo "create_required_dirs() starts ..."
    set -x
    mkdir -p $LEDGER_DIR
    mkdir -p $DISCOVERY_NODES_DATA/certs
    mkdir -p $INSGORUND_DATA
    mkdir -p $CONFIGS_DIR

    touch $INSGORUND_PORT_FILE
    { set +x; } 2>/dev/null

    for i in `seq 1 $NUM_DISCOVERY_NODES`
    do
        set -x
        mkdir -p ${DISCOVERY_NODE_LOGS}${i}
        { set +x; } 2>/dev/null
    done

    echo "create_required_dirs() end."
}

generate_insolard_configs()
{
    echo "generate configs"
    set -x
    go run scripts/generate_insolar_configs.go \
        -p $INSGORUND_PORT_FILE
    { set +x; } 2>/dev/null
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
    echo "generate bootstrap keys: $KEYS_FILE"
    bin/insolar gen-key-pair > $KEYS_FILE
}

generate_root_member_keys()
{
    echo "generate root member_keys: $ROOT_MEMBER_KEYS_FILE"
    bin/insolar gen-key-pair > $ROOT_MEMBER_KEYS_FILE
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

        $INSGORUND -l $host:$listen_port --rpc $host:$rpc_port --log-level=$GORUND_LOG_LEVEL --metrics :$metrics_port &> $INSGORUND_DATA$rpc_port.log &

    done < "$INSGORUND_PORT_FILE"
}

copy_data()
{
    echo "copy_data() starts ..."

    echo "prepare heavy data"
    set -x
    mv ${LEDGER_DIR}/ ${DISCOVERY_NODES_HEAVY_DATA}/data
    { set +x; } 2>/dev/null

    echo "prepare legacy dirs"
    for node in "${DISCOVERY_NODE_DIRS[@]}"
    do
        set -x
        mkdir -p ${node}/data
        { set +x; } 2>/dev/null
    done
    echo "copy_data() end."
}

copy_discovery_certs()
{
    echo "copy_certs() starts ..."
    i=0
    for node in "${DISCOVERY_NODE_DIRS[@]}"
    do
        i=$((i + 1))
        set -x
        cp -v $DISCOVERY_NODES_DATA/certs/discovery_cert_$i.json $node/cert.json
        { set +x; } 2>/dev/null
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

    echo "start genesis ..."
    CMD="$INSOLARD --config $INSOLARD_CONFIG --genesis $GENESIS_CONFIG --keyout $DISCOVERY_NODES_DATA/certs"

    GENESIS_EXIT_CODE=0
    set +e
    if [[ "$NO_GENESIS_LOG_REDIRECT" != "1" ]]; then
        set -x
        ${CMD} &> ${LAUNCHNET_LOGS_DIR}genesis_output.log
        GENESIS_EXIT_CODE=$?
        { set +x; } 2>/dev/null
    else
        set -x
        ${CMD}
        GENESIS_EXIT_CODE=$?
        { set +x; } 2>/dev/null
    fi
    set -e
    if [[ $GENESIS_EXIT_CODE -ne 0 ]]; then
        echo "Genesis failed"
        if [[ "$NO_GENESIS_LOG_REDIRECT" != "1" ]]; then
            echo "check log: ${LAUNCHNET_LOGS_DIR}/genesis_output.log"
        fi
        exit $GENESIS_EXIT_CODE
    fi
    echo "genesis is done"

    copy_data
    copy_discovery_certs
}

run_insgorund=true
watch_pulse=true
check_working_dir
process_input_params $@

trap 'stop_listening true' INT TERM EXIT

echo "start pulsar ..."
echo "   log: ${LAUNCHNET_LOGS_DIR}pulsar_output.log"
set -x
mkdir -p $PULSAR_DATA_DIR
${PULSARD} -c ${PULSAR_CONFIG} --trace &> ${LAUNCHNET_LOGS_DIR}pulsar_output.log &
{ set +x; } 2>/dev/null

if [[ "$run_insgorund" == "true" ]]
then
    echo "start insgorund ..."
    launch_insgorund
else
    echo "insgorund launch skip"
fi

echo "start discovery nodes ..."

for i in `seq 1 $NUM_DISCOVERY_NODES`
do
    set -x
    INSOLAR_LOG_LEVEL=$insolar_log_level $INSOLARD \
        --config ${DISCOVERY_NODES_DATA}${i}/insolard.yaml \
        --trace &> ${DISCOVERY_NODE_LOGS}${i}/output.log &
    { set +x; } 2>/dev/null
    echo "discovery node $i started in background"
done

echo "discovery nodes started ..."

if [[ "$NUM_NODES" -ne "0" ]]
then
    wait_for_complete_network_state
    ./scripts/insolard/start_nodes.sh
fi

if [[ "$watch_pulse" == "true" ]]
then
    echo "starting pulse watcher..."
    echo "${PULSEWATCHER} -c ${PULSEWATCHER_CONFIG}"
    ${PULSEWATCHER} -c ${PULSEWATCHER_CONFIG}
else
    echo "waiting..."
    wait
fi

echo "FINISHING ..."
