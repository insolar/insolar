#!/usr/bin/env bash
set -em
# requires: lsof, awk, sed, grep, pgrep

export GO111MODULE=on

# Changeable environment variables (parameters)
INSOLAR_ARTIFACTS_DIR=${INSOLAR_ARTIFACTS_DIR:-".artifacts"}/
LAUNCHNET_BASE_DIR=${LAUNCHNET_BASE_DIR:-"${INSOLAR_ARTIFACTS_DIR}launchnet"}/

INSOLAR_LOG_FORMATTER=${INSOLAR_LOG_FORMATTER:-""}
INSOLAR_LOG_LEVEL=${INSOLAR_LOG_LEVEL:-"debug"}
GORUND_LOG_LEVEL=${GORUND_LOG_LEVEL:-${INSOLAR_LOG_LEVEL}}
# we can skip build binaries (by default in CI environment they skips)
SKIP_BUILD=${SKIP_BUILD:-${CI_ENV}}
BUILD_TAGS=${BUILD_TAGS:-'-tags "debug functest"'}

# predefined/dependent environment variables

LAUNCHNET_LOGS_DIR=${LAUNCHNET_BASE_DIR}logs/
DISCOVERY_NODE_LOGS=${LAUNCHNET_LOGS_DIR}discoverynodes/
INSGORUND_LOGS=${LAUNCHNET_LOGS_DIR}insgorund/

BIN_DIR=bin
INSOLAR_CLI=${BIN_DIR}/insolar
INSOLARD=$BIN_DIR/insolard
INSGORUND=$BIN_DIR/insgorund
KEEPERD=$BIN_DIR/keeperd
PULSARD=$BIN_DIR/pulsard
PULSEWATCHER=$BIN_DIR/pulsewatcher

# DUMP_METRICS_ENABLE enables metrics dump to logs dir after every functest
DUMP_METRICS_ENABLE=${DUMP_METRICS_ENABLE:-"1"}

PULSAR_DATA_DIR=${LAUNCHNET_BASE_DIR}pulsar_data
PULSAR_CONFIG=${LAUNCHNET_BASE_DIR}pulsar.yaml

SCRIPTS_DIR=scripts/insolard/

CONFIGS_DIR=${LAUNCHNET_BASE_DIR}configs/

PULSAR_KEYS=${CONFIGS_DIR}pulsar_keys.json
HEAVY_GENESIS_CONFIG_FILE=${CONFIGS_DIR}heavy_genesis.json
CONTRACTS_PLUGINS_DIR=${LAUNCHNET_BASE_DIR}contracts

DISCOVERY_NODES_DATA=${LAUNCHNET_BASE_DIR}discoverynodes/

DISCOVERY_NODES_HEAVY_DATA=${DISCOVERY_NODES_DATA}1/

BOOTSTRAP_TEMPLATE=${SCRIPTS_DIR}bootstrap_template.yaml
BOOTSTRAP_CONFIG=${LAUNCHNET_BASE_DIR}bootstrap.yaml
BOOTSTRAP_INSOLARD_CONFIG=${LAUNCHNET_BASE_DIR}insolard.yaml

KEEPERD_CONFIG=${LAUNCHNET_BASE_DIR}keeperd.yaml
KEEPERD_LOG=${LAUNCHNET_LOGS_DIR}keeperd.log

PULSEWATCHER_CONFIG=${LAUNCHNET_BASE_DIR}pulsewatcher.yaml

INSGORUND_PORT_FILE=${CONFIGS_DIR}insgorund_ports.txt

set -x
export INSOLAR_LOG_FORMATTER=${INSOLAR_LOG_FORMATTER}
export INSOLAR_LOG_LEVEL=${INSOLAR_LOG_LEVEL}
{ set +x; } 2>/dev/null

NUM_DISCOVERY_NODES=$(sed '/^nodes:/ q' $BOOTSTRAP_TEMPLATE | grep "host:" | grep -v "#" | wc -l | tr -d '[:space:]')
NUM_NODES=$(sed -n '/^nodes:/,$p' $BOOTSTRAP_TEMPLATE | grep "host:" | grep -v "#" | wc -l | tr -d '[:space:]')
echo "discovery+other nodes: ${NUM_DISCOVERY_NODES}+${NUM_NODES}"

for i in `seq 1 $NUM_DISCOVERY_NODES`
do
    DISCOVERY_NODE_DIRS+=(${DISCOVERY_NODES_DATA}${i})
done

# LOGROTATOR_ENABLE enables log rotation before every functest start
LOGROTATOR_ENABLE=${LOGROTATOR_ENABLE:-""}
LOGROTATOR=tee
LOGROTATOR_BIN=${LAUNCHNET_BASE_DIR}inslogrotator
if [[ "$LOGROTATOR_ENABLE" == "1" ]]; then
  LOGROTATOR=${LOGROTATOR_BIN}
fi

build_logger()
{
    echo "build logger binaries"
    set -x
    pushd scripts/_logger
    GO111MODULE=on go build -o inslogrotator .
    popd
    mv scripts/_logger/inslogrotator ${LOGROTATOR_BIN}
    { set +x; } 2>/dev/null
}

kill_port()
{
    port=$1
    pids=$(lsof -i :$port | grep "LISTEN\|UDP" | awk '{print $2}')
    for pid in $pids
    do
        echo -n "killing pid $pid at "
        date
        kill $pid
    done
}

kill_all()
{
  echo "kill all processes: insgorund, insolard, pulsard"
  set +e
  killall insgorund
  killall insolard
  killall pulsard
  set -e
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

    transport_ports=$( grep "host:" ${BOOTSTRAP_CONFIG} | grep -o ":\d\+" | grep -o "\d\+" | tr '\n' ' ' )
    keeperd_port=$( grep "listenaddress:" ${KEEPERD_CONFIG} | grep -o ":\d\+" | grep -o "\d\+" | tr '\n' ' ' )
    ports="$ports $transport_ports $keeperd_port"

    for port in $ports
    do
        echo "killing process using port '$port'"
        kill_port $port
    done

    echo "stop_listening() end."
}

stop_all()
{
  stop_listening true
  kill_all
}

clear_dirs()
{
    echo "clear_dirs() starts ..."
    set -x
    rm -rfv ${DISCOVERY_NODES_DATA}
    rm -rfv ${LAUNCHNET_LOGS_DIR}
    rm -rfv ${CONTRACTS_PLUGINS_DIR}
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
    mkdir -p ${DISCOVERY_NODES_DATA}certs
    mkdir -p ${CONFIGS_DIR}

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
    go run -mod=vendor scripts/generate_insolar_configs.go -p ${INSGORUND_PORT_FILE}
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
    echo "build binaries"
    set -x
    export BUILD_TAGS
    make build
    { set +x; } 2>/dev/null
}

rebuild_binaries()
{
    echo "rebuild binaries"
    make clean
    build_binaries
}

generate_pulsar_keys()
{
    echo "generate pulsar keys: ${PULSAR_KEYS}"
    bin/insolar gen-key-pair > ${PULSAR_KEYS}
}

generate_root_member_keys()
{
    echo "generate members keys in dir: $CONFIGS_DIR"
    bin/insolar gen-key-pair > ${CONFIGS_DIR}root_member_keys.json
    bin/insolar gen-key-pair > ${CONFIGS_DIR}fee_member_keys.json
    bin/insolar gen-key-pair > ${CONFIGS_DIR}migration_admin_member_keys.json
    for (( b = 0; b < 10; b++ ))
    do
    bin/insolar gen-key-pair > ${CONFIGS_DIR}migration_daemon_${b}_member_keys.json
    done

    for (( b = 0; b < 30; b++ ))
    do
    bin/insolar gen-key-pair > ${CONFIGS_DIR}network_incentives_${b}_member_keys.json
    done

    for (( b = 0; b < 30; b++ ))
    do
    bin/insolar gen-key-pair > ${CONFIGS_DIR}application_incentives_${b}_member_keys.json
    done

    for (( b = 0; b < 30; b++ ))
    do
    bin/insolar gen-key-pair > ${CONFIGS_DIR}foundation_${b}_member_keys.json
    done

    for (( b = 0; b < 1; b++ ))
    do
    bin/insolar gen-key-pair > ${CONFIGS_DIR}funds_${b}_member_keys.json
    done

    for (( b = 0; b < 4; b++ ))
    do
    bin/insolar gen-key-pair > ${CONFIGS_DIR}enterprise_${b}_member_keys.json
    done
}

generate_migration_addresses()
{
    echo "generate migration addresses: ${CONFIGS_DIR}migration_addresses.json"
    bin/insolar gen-migration-addresses > ${CONFIGS_DIR}migration_addresses.json
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
    echo -e "\t-g - start launchnet"
    echo -e "\t-b - do bootstrap only and exit, show bootstrap logs"
    echo -e "\t-l - clear all and exit"
    echo -e "\t-C - generate configs only"
    echo -e "\t-w - start without pulse watcher"
}

process_input_params()
{
    # shell does not reset OPTIND automatically;
    # it must be manually reset between multiple calls to getopts
    # within the same shell invocation if a new set of parameters is to be used
    OPTIND=1
    while getopts "h?ngblwC" opt; do
        case "$opt" in
        h|\?)
            usage
            exit 0
            ;;
        n)
            run_insgorund=false
            ;;
        g)
            GENESIS=1
            bootstrap
            ;;
        b)
            NO_BOOTSTRAP_LOG_REDIRECT=1
            NO_STOP_LISTENING_ON_PREPARE=${NO_STOP_LISTENING_ON_PREPARE:-"1"}
            bootstrap
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

launch_keeperd()
{
    echo "launch_keeperd() starts ..."
    ${KEEPERD} --config=${KEEPERD_CONFIG} \
    &> ${KEEPERD_LOG} &

    echo "launch_keeperd() end."
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
        num=`scripts/insolard/check_status.sh 2>/dev/null | grep "CompleteNetworkState" | wc -l`
        echo "$num/$NUM_DISCOVERY_NODES discovery nodes ready"
        if [[ "$num" -eq "$NUM_DISCOVERY_NODES" ]]
        then
            break
        fi
        sleep 5s
    done
}

bootstrap()
{
    echo "bootstrap start"
    prepare
    if [[ "$SKIP_BUILD" != "1" ]]; then
        build_binaries
    else
        echo "SKIP: build binaries (SKIP_BUILD=$SKIP_BUILD)"
    fi
    generate_pulsar_keys
    generate_root_member_keys
    generate_insolard_configs
    generate_migration_addresses

    echo "start bootstrap ..."
    CMD="${INSOLAR_CLI} bootstrap --config=${BOOTSTRAP_CONFIG} --certificates-out-dir=${DISCOVERY_NODES_DATA}certs"

    GENESIS_EXIT_CODE=0
    set +e
    if [[ "$NO_BOOTSTRAP_LOG_REDIRECT" != "1" ]]; then
        set -x
        ${CMD} &> ${LAUNCHNET_LOGS_DIR}bootstrap.log
        GENESIS_EXIT_CODE=$?
        { set +x; } 2>/dev/null
        echo "bootstrap log: ${LAUNCHNET_LOGS_DIR}bootstrap.log"
    else
        set -x
        ${CMD}
        BOOTSTRAP_EXIT_CODE=$?
        { set +x; } 2>/dev/null
    fi
    set -e
    if [[ ${BOOTSTRAP_EXIT_CODE} -ne 0 ]]; then
        echo "Genesis failed"
        if [[ "${NO_BOOTSTRAP_LOG_REDIRECT}" != "1" ]]; then
            echo "check log: ${LAUNCHNET_LOGS_DIR}/bootstrap.log"
        fi
        exit ${BOOTSTRAP_EXIT_CODE}
    fi
    echo "bootstrap is done"

    copy_discovery_certs
}

run_insgorund=true
watch_pulse=true
check_working_dir
process_input_params $@

kill_all
trap 'stop_all' INT TERM EXIT

echo "start pulsar ..."
echo "   log: ${LAUNCHNET_LOGS_DIR}pulsar_output.log"
set -x
mkdir -p ${PULSAR_DATA_DIR}
${PULSARD} -c ${PULSAR_CONFIG} &> ${LAUNCHNET_LOGS_DIR}pulsar_output.log &
{ set +x; } 2>/dev/null
echo "pulsar log: ${LAUNCHNET_LOGS_DIR}pulsar_output.log"

launch_keeperd

if [[ "$run_insgorund" == "true" ]]
then
    echo "start insgorund ..."
    launch_insgorund
else
    echo "insgorund launch skip"
fi

handle_sigchld()
{
  jobs -pn
  echo "someone left the network"
}

if [[ "$LOGROTATOR_ENABLE" == "1" ]]; then
  echo "prepare logger"
  build_logger
fi

trap 'handle_sigchld' SIGCHLD

echo "start heavy node"
set -x
$INSOLARD \
    --config ${DISCOVERY_NODES_DATA}1/insolard.yaml \
    --heavy-genesis ${HEAVY_GENESIS_CONFIG_FILE} \
    2>&1 | ${LOGROTATOR} ${DISCOVERY_NODE_LOGS}1/output.log > /dev/null &
{ set +x; } 2>/dev/null
echo "heavy node started in background"
echo "log: ${DISCOVERY_NODE_LOGS}1/output.log"

echo "start discovery nodes ..."
for i in `seq 2 $NUM_DISCOVERY_NODES`
do
    set -x
    $INSOLARD \
        --config ${DISCOVERY_NODES_DATA}${i}/insolard.yaml \
        2>&1 | ${LOGROTATOR} ${DISCOVERY_NODE_LOGS}${i}/output.log > /dev/null &
    { set +x; } 2>/dev/null
    echo "discovery node $i started in background"
    echo "log: ${DISCOVERY_NODE_LOGS}${i}/output.log"
done

echo "discovery nodes started ..."

if [[ "$NUM_NODES" -ne "0"  && "$run_insgorund" == "true" ]]
then
    wait_for_complete_network_state
    if [[ "$GENESIS" == "1" ]]; then
        ./scripts/insolard/start_nodes.sh -g
    else
        ./scripts/insolard/start_nodes.sh
    fi
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
