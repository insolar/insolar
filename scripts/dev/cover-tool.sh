#!/usr/bin/env bash
#
# Tool for simple code coverage inspection.
#
# Usage examples:
#
# ./scripts/dev/cover-tool.sh
#
# ./scripts/dev/cover-tool.sh ledger
#
# ./scripts/dev/cover-tool.sh instrumentation -v
#
# COVERPROFILE=./coverage.txt ./scripts/dev/cover-tool.sh -v
#
# depends on gocov:
#
#   go get github.com/axw/gocov/gocov
#
# optional tool (adds browser report in verbose mode):
#
#   go get gopkg.in/matm/v1/gocov-html
#

# set 'bash strict mode'
set -euo pipefail
IFS=$'\n\t'
# set -x

# parse arguments
COVER_TARGET=${1:-""}
shift || true
# echo "ok"

ARG=${1:-""}
VERBOSE=""
if [[ ! -z "$ARG" ]]; then
    case $ARG in
        -v|--verbose)
            VERBOSE="1"
        ;;
        *)
            echo "Unknown option $ARG"
            exit 1
    esac
fi

# set vars
OUTPUT_DIR=${OUTPUT_DIR:-""}
if [[ -z "$OUTPUT_DIR" ]]; then
    OUTPUT_DIR=$(mktemp -d)
fi

TESTED_PACKAGES=./$COVER_TARGET/...
if [[ -z "$COVER_TARGET" ]]; then
    TESTED_PACKAGES=./...
fi

COVERPROFILE=${COVERPROFILE:-""}
if [[ -z "$COVERPROFILE" ]]; then
    # collect coverage
    export COVERPROFILE=$OUTPUT_DIR/coverage.out
    export TESTED_PACKAGES
    make test_with_coverage
fi

# produce report
GOCOV_FILE=$OUTPUT_DIR/coverage.gocov
gocov convert $COVERPROFILE > $GOCOV_FILE
if [[ -z "$VERBOSE" ]]; then
    gocov report $GOCOV_FILE | grep -E -v '\S+.go\s+' | grep '.' | awk '/^(.*)([[:space:]]+-+[[:space:]]+)(.*)$/{print $1 "\t" $3}'
    gocov report $GOCOV_FILE | grep 'Total'
else
    gocov report $GOCOV_FILE

    if which gocov-html; then
        gocov convert $COVERPROFILE > $GOCOV_FILE.json
        gocov-html $GOCOV_FILE.json > $GOCOV_FILE.html
        open $GOCOV_FILE.html
    fi
    go tool cover -html=$COVERPROFILE
fi
