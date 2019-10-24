set -o pipefail

RSYNC_CODE=0
RSYNC_CMD="rsync t 3 "
output=$( $RSYNC_CMD 2>&1  ) || RSYNC_CODE=$?

if [ "$RSYNC_CODE" == "23" ]
then
    echo ZHOOOOOOPA
fi

echo "OUTPUT: $output"
echo "EXIT CODE: $RSYNC_CODE"

echo "TEST"

while true
do
    echo "PPPP"
done

exit 0

////////////////////////////////////



#!/usr/bin/env bash
set -e
set -o pipefail

BADGER_DIR=/opt/insolar/data
REMOTE_HOST=dev-insolar-alpha-6.vm.insolar.io
DATE=$(date +%H:%m:%S-%s)

function log() {
  echo "$(date --rfc-3339=seconds) $1"
}

log "before finding ${DATE}"

find /opt/insolar/backup/target -mmin +3 -type d -name 'pulse-*' -exec rm -rv {} +


log "starting backup ${DATE}"

RSYNC_CMD="rsync -avW --info=PROGRESS2 --delete ${BADGER_DIR} ${REMOTE_HOST}::backup"
RSYNC_CODE=0
$RSYNC_CMD || RSYNC_CODE=$?
if [ "$RSYNC_CODE" == "24" ]
then
    echo "Skip error code 24: skip vanishing"
else
    exit $RSYNC_CODE
fi

iter=0
log "between rsyncs ${DATE}"
while true
do
  ((iter++))
  "rsync iter $iter. ${DATE}"

  output=$( $RSYNC_CMD || RSYNC_CODE=$? )
  if [ "$RSYNC_CODE" == "24" ]
  then
    echo "Skip error code 24: skip vanishing"
  else
    exit $RSYNC_CODE
  fi
  if ! echo "$output" | grep -q "(MANIFEST\|\.sst)$"
  then
      break
  fi
done

touch ${INSOLAR_CURRENT_BACKUP_DIR}/success
log "backup done ${DATE}"











