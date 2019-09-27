#!/usr/bin/env bash
# bash strict mode
set -euo pipefail
IFS=$'\n\t'

DATE_CMD=date
if which gdate >/dev/null; then
# if Mac and `brew install coreutils` performed
    DATE_CMD=gdate
fi

usage()
{
    echo "usage: $0 [options]"
    echo "possible options: "
    echo -e "\t-h - show help"
    echo -e "\t-t - git HEAD ref time"
    echo -e "\t-t - git HEAD ref date"
}


gnu_date_available()
{
    ${DATE_CMD} -d "@1568213332" &> /dev/null && RC=$? || RC=$?
    return $RC
}


while getopts "h?td" opt; do
    case "$opt" in
    h|\?)
        usage
        exit 0
        ;;
    t)
        if gnu_date_available ; then
            GIT_TIMESTAMP=$(git show -s --format=%ct)
            ${DATE_CMD} -d "@$GIT_TIMESTAMP" +'%H:%M:%S'
        else
            echo "<TIME>"
        fi
        ;;
    d)
        if gnu_date_available ; then
            GIT_TIMESTAMP=$(git show -s --format=%ct)
            ${DATE_CMD} -d "@$GIT_TIMESTAMP" +'%Y-%m-%d'
        else
            echo "<DATE>"
        fi
        ;;
    esac
done
