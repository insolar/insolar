#!/usr/bin/env bash


alias cdins='cd ~/go/src/github.com/insolar/insolar'
alias insps='ps awux | egrep "insolard|insgorund"'
alias launet='./scripts/insolard/launchnet.sh -g'
alias bench='./bin/benchmark -c=64 -r=1024 -k=scripts/insolard/configs/root_member_keys.json'

check_same_pulse() { # network launched with same pulse
    find . -name "*.log" -exec grep -m1 --mmap persist {} \; | \
    perl -ne '/current_pulse=(\d+)/; print "$1 ";'
}

get_pprof_web() {
    curl http://localhost:8005/debug/pprof/profile\?seconds\=20 >tmp.prof
    go tool pprof -http localhost:3001 tmp.prof
}