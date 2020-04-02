#!/usr/bin/env bash

set -u

exitcode=0
end=$((SECONDS+TESTS_TIMEOUT))
while [ $SECONDS -lt $end ]; do
  make ${MAKE_TARGET} || exitcode=${exitcode}+1
done

exit ${exitcode}
