#!/usr/bin/env bash
set -e
# requires: lsof, awk, sed, grep, pgrep

# Changeable environment variables (parameters)
INSOLAR_ARTIFACTS_DIR=${INSOLAR_ARTIFACTS_DIR:-".artifacts"}/
LAUNCHNET_BASE_DIR=${LAUNCHNET_BASE_DIR:-"${INSOLAR_ARTIFACTS_DIR}launchnet"}/

INSOLAR_LOG_FORMATTER=${INSOLAR_LOG_FORMATTER:-"text"}
INSOLAR_LOG_LEVEL=${INSOLAR_LOG_LEVEL:-"debug"}
GORUND_LOG_LEVEL=${GORUND_LOG_LEVEL:-${INSOLAR_LOG_LEVEL}}
# we can skip build binaries (by default in CI environment they skips)
SKIP_BUILD=${SKIP_BUILD:-${CI_ENV}}

# predefined/dependent environment variables

LAUNCHNET_LOGS_DIR=${LAUNCHNET_BASE_DIR}logs/
DISCOVERY_NODE_LOGS=${LAUNCHNET_LOGS_DIR}discoverynodes/
INSGORUND_LOGS=${LAUNCHNET_LOGS_DIR}insgorund/

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

KEYS_FILE=${CONFIGS_DIR}/bootstrap_keys.json
ROOT_MEMBER_KEYS_FILE=${CONFIGS_DIR}/root_member_keys.json
HEAVY_GENESIS_CONFIG_FILE=${CONFIGS_DIR}/heavy_genesis.json

MD_ADMIN_MEMBER_KEYS_FILE=${CONFIGS_DIR}/md_admin_member_keys.json
ORACLE0_MEMBER_KEYS_FILE=${CONFIGS_DIR}/oracle0_member_keys.json
ORACLE1_MEMBER_KEYS_FILE=${CONFIGS_DIR}/oracle1_member_keys.json
ORACLE2_MEMBER_KEYS_FILE=${CONFIGS_DIR}/oracle2_member_keys.json
ORACLE3_MEMBER_KEYS_FILE=${CONFIGS_DIR}/oracle3_member_keys.json
ORACLE4_MEMBER_KEYS_FILE=${CONFIGS_DIR}/oracle4_member_keys.json
ORACLE5_MEMBER_KEYS_FILE=${CONFIGS_DIR}/oracle5_member_keys.json
ORACLE6_MEMBER_KEYS_FILE=${CONFIGS_DIR}/oracle6_member_keys.json
ORACLE7_MEMBER_KEYS_FILE=${CONFIGS_DIR}/oracle7_member_keys.json
ORACLE8_MEMBER_KEYS_FILE=${CONFIGS_DIR}/oracle8_member_keys.json
ORACLE9_MEMBER_KEYS_FILE=${CONFIGS_DIR}/oracle9_member_keys.json

# TODO: use only heavy matereal data dir
DISCOVERY_NODES_DATA=${LAUNCHNET_BASE_DIR}discoverynodes/

DISCOVERY_NODES_HEAVY_DATA=${DISCOVERY_NODES_DATA}1/

GENESIS_TEMPLATE=${SCRIPTS_DIR}bootstrap/genesis_template.yaml
BOOTSTRAP_GENESIS_CONFIG=${LAUNCHNET_BASE_DIR}genesis.yaml
BOOTSTRAP_INSOLARD_CONFIG=${LAUNCHNET_BASE_DIR}insolard.yaml

PULSEWATCHER_CONFIG=${LAUNCHNET_BASE_DIR}pulsewatcher.yaml

INSGORUND_PORT_FILE=${CONFIGS_DIR}/insgorund_ports.txt

set -x
export INSOLAR_LOG_FORMATTER=${INSOLAR_LOG_FORMATTER}
export INSOLAR_LOG_LEVEL=${INSOLAR_LOG_LEVEL}
{ set +x; } 2>/dev/null

NUM_DISCOVERY_NODES=$(sed '/^nodes:/ q' $GENESIS_TEMPLATE | grep "host:" | grep -v "#" | wc -l | tr -d '[:space:]')
NUM_NODES=$(sed -n '/^nodes:/,$p' $GENESIS_TEMPLATE | grep "host:" | grep -v "#" | wc -l | tr -d '[:space:]')
echo "discovery+other nodes: ${NUM_DISCOVERY_NODES}+${NUM_NODES}"

for i in `seq 1 $NUM_DISCOVERY_NODES`
do
    DISCOVERY_NODE_DIRS+=(${DISCOVERY_NODES_DATA}${i})
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

    transport_ports=$( grep "host:" ${BOOTSTRAP_GENESIS_CONFIG} | grep -o ":\d\+" | grep -o "\d\+" | tr '\n' ' ' )
    ports="$ports $transport_ports"

    for port in $ports
    do
        echo "killing process using port '$port'"
        kill_port $port
    done

    echo "stop_listening() end."
}

clear_dirs()
{
    echo "clear_dirs() starts ..."
    set -x
    rm -rfv ${LEDGER_DIR}
    rm -rfv ${DISCOVERY_NODES_DATA}
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
    mkdir -p ${DISCOVERY_NODES_DATA}certs
    mkdir -p $CONFIGS_DIR

    mkdir -p ${INSGORUND_LOGS}
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

    echo "generate admin member_keys: $MD_ADMIN_MEMBER_KEYS_FILE"
    bin/insolar gen-key-pair > $MD_ADMIN_MEMBER_KEYS_FILE

    echo "generate oracles member_keys"
    bin/insolar gen-key-pair > $ORACLE0_MEMBER_KEYS_FILE
    bin/insolar gen-key-pair > $ORACLE1_MEMBER_KEYS_FILE
    bin/insolar gen-key-pair > $ORACLE2_MEMBER_KEYS_FILE
    bin/insolar gen-key-pair > $ORACLE3_MEMBER_KEYS_FILE
    bin/insolar gen-key-pair > $ORACLE4_MEMBER_KEYS_FILE
    bin/insolar gen-key-pair > $ORACLE5_MEMBER_KEYS_FILE
    bin/insolar gen-key-pair > $ORACLE6_MEMBER_KEYS_FILE
    bin/insolar gen-key-pair > $ORACLE7_MEMBER_KEYS_FILE
    bin/insolar gen-key-pair > $ORACLE8_MEMBER_KEYS_FILE
    bin/insolar gen-key-pair > $ORACLE9_MEMBER_KEYS_FILE
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

        ${INSGORUND} \
            -l ${host}:${listen_port} \
            --rpc ${host}:${rpc_port} \
            --log-level=${GORUND_LOG_LEVEL} \
            --metrics :${metrics_port} \
            &> ${INSGORUND_LOGS}${rpc_port}.log &

    done < "${INSGORUND_PORT_FILE}"
}

copy_data()
{
    echo "copy data dir to heavy"
    set -x
    mv ${LEDGER_DIR}/ ${DISCOVERY_NODES_HEAVY_DATA}data
    { set +x; } 2>/dev/null
}

copy_discovery_certs()
{
    echo "copy_certs() starts ..."
    i=0
    for node in "${DISCOVERY_NODE_DIRS[@]}"
    do
        i=$((i + 1))
        set -x
        cp -v ${DISCOVERY_NODES_DATA}certs/discovery_cert_$i.json ${node}/cert.json
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
        echo "SKIP: build binaries (SKIP_BUILD=$SKIP_BUILD)"
    fi
    generate_bootstrap_keys
    generate_root_member_keys
    generate_insolard_configs

    echo "start genesis ..."
    CMD="$INSOLARD --config ${BOOTSTRAP_INSOLARD_CONFIG} --genesis ${BOOTSTRAP_GENESIS_CONFIG} --keyout ${DISCOVERY_NODES_DATA}certs"

    GENESIS_EXIT_CODE=0
    set +e
    if [[ "$NO_GENESIS_LOG_REDIRECT" != "1" ]]; then
        set -x
        ${CMD} &> ${LAUNCHNET_LOGS_DIR}genesis_output.log
        GENESIS_EXIT_CODE=$?
        { set +x; } 2>/dev/null
        echo "genesis log: ${LAUNCHNET_LOGS_DIR}genesis_output.log"
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
echo "pulsar log: ${LAUNCHNET_LOGS_DIR}pulsar_output.log"

if [[ "$run_insgorund" == "true" ]]
then
    echo "start insgorund ..."
    launch_insgorund
else
    echo "insgorund launch skip"
fi

echo "start heavy node"
set -x
$INSOLARD \
    --config ${DISCOVERY_NODES_DATA}1/insolard.yaml \
    --heavy-genesis ${HEAVY_GENESIS_CONFIG_FILE} \
    --trace &> ${DISCOVERY_NODE_LOGS}1/output.log &
{ set +x; } 2>/dev/null
echo "heavy node started in background"
echo "log: ${DISCOVERY_NODE_LOGS}1/output.log"

echo "start discovery nodes ..."
for i in `seq 2 $NUM_DISCOVERY_NODES`
do
    set -x
    $INSOLARD \
        --config ${DISCOVERY_NODES_DATA}${i}/insolard.yaml \
        --trace &> ${DISCOVERY_NODE_LOGS}${i}/output.log &
    { set +x; } 2>/dev/null
    echo "discovery node $i started in background"
    echo "log: ${DISCOVERY_NODE_LOGS}${i}/output.log"
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
