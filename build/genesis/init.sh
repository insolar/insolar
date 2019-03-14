#!/bin/sh
set -e
GENESIS_OUTPUT_DIR=$GENESIS_OUTPUT_DIR
if [ -z "$GENESIS_OUTPUT_DIR" ]; then
    echo "GENESIS_OUTPUT_DIR variable not set. Exit."
    exit 1
fi

if [ -f "${GENESIS_OUTPUT_DIR}/*.sst" ]; then
    echo "GENESIS_OUTPUT_DIR is not empty. Exit."
    exit
fi

echo "copy genesis to ${GENESIS_OUTPUT_DIR}"
cp -r /genesis/. ${GENESIS_OUTPUT_DIR}/
