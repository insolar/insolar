#!/usr/bin/env bash

URL="$(awk 'END{print $1}' /etc/hosts):19191/healthcheck"

http_code=$(curl -s -o /dev/null -w "%{http_code}" $URL)
if [[ ${http_code} != "200" ]]; then
    echo "Bad HTTP healthcheck code on $URL, expected '200', got '${http_code}'"
    exit 1
fi

exit 0
