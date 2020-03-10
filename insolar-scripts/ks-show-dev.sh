#! /usr/bin/env bash
# Helper which shows ksonnet output.
# Useful for syntax error checks and checking actual ksonnet output.

pushd ci/local-dev/manifests/ks
echo "start check"
echo "hint: ksonnet is slow, please be patient"
if which yamlyaml; then
    ks show dev | yamlyaml
else
    ks show dev
    echo "If you want to show nice output for stringified json in ConfigMap, install yamlyaml tool:"
    echo "  go install github.com/nordicdyno/yamlyaml"
fi
