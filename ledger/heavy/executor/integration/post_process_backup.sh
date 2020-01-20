#!/usr/bin/env bash

# AALEKSEEV TODO delete?

badgerDir=$1
recoverDB=$2

set -o history
set -o histexpand

echo "BADGER DB: $badgerDir"
echo "RECOVER DB: $recoverDB"

# Makes a complete copy of a Badger database directory.
# Repeat rsync if the MANIFEST and SSTables are updated.
rsync -avz --delete $badgerDir/ $recoverDB/
while !! | grep -q "(MANIFEST\|\.sst)$"; do :; done

touch $INSOLAR_CURRENT_BACKUP_DIR/BACKUPED
