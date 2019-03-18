#! /usr/bin/env bash
#
# Example of profiling 60 second profile on all insolard node
# (by default profiles 30 seconds):
#
# ./scripts/insolard/profile.sh [60]

prof_time=${1:-"30"}
prof_files_dir=pprof
current_dir=$( dirname $0 )
web_profile_port=8080

trap 'killall' INT TERM EXIT

killall() {
    trap '' INT TERM     # ignore INT and TERM while shutting down
    echo "Shutting down..."
    kill -TERM 0         # fixed order, send TERM not INT
    wait
    echo DONE
}

confs=${current_dir}"/configs/generated_configs/*nodes/insolar_*.yaml"
prof_ports=$( grep listenaddress ${confs} |  grep -o ":\d\+" | grep -o "\d\+" | tr '\n' ' ' )

mkdir -p ${current_dir}/${prof_files_dir}

echo "Fetching profile data from insolar nodes..."
for port in ${prof_ports}
do
    curl "http://localhost:${port}/debug/pprof/profile?seconds=${prof_time}" --output ${current_dir}/${prof_files_dir}/prof_${port} &> /dev/null &
done
wait

echo "Starting web servers with profile info..."

i=${web_profile_port}
for port in ${prof_ports}
do
    go tool pprof -http=:${i} ${current_dir}/${prof_files_dir}/prof_${port} &
    echo "Started web profile server on localhost:${i}/ui"
    i=$((i + 1))
done
wait
