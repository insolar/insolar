#!/usr/bin/env bash

set -u

# pass make target(s) as arguments when running the script
[[ "$#" == 0 ]] && return 1

exitcode=0
end=$((SECONDS+TESTS_TIMEOUT))
while [ $SECONDS -lt $end ]; do
  make $@ || exitcode=${exitcode}+1
done

exit ${exitcode}
