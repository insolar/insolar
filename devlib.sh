#!/usr/bin/env bash
alias insps='ps awux | egrep "insolard|insgorund"'
alias launet='./scripts/insolard/launchnet.sh -g'
alias bench='./bin/benchmark -c=64 -r=1024 -k=scripts/insolard/configs/root_member_keys.json'

alias check_same_pulse='find . -name "*.log" -exec grep -m1 --mmap persist {} \; | ' # TODO add awk comparator