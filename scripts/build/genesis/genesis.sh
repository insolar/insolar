#!/bin/sh
set -e

echo ==============================
echo Making genesis
echo ==============================

SCRIPTPATH=`dirname $0`
INSOLAR=`git -C $SCRIPTPATH rev-parse --show-toplevel`

dirs="data certs bin log"
for d in $dirs; do rm -rf $d; done
for d in $dirs; do mkdir $d; done
rm -f insolar.cfg.yaml

make -C $INSOLAR build

cp $INSOLAR/bin/insgorund bin/

ROOT_MEMBER_KEY=keys/root_member.key.json
BOOTSTRAP_KEY=keys/bootstrap.key.json
[[ -d keys ]] || mkdir keys
[[ -f $ROOT_MEMBER_KEY ]] || $INSOLAR/bin/insolar -c gen_keys > $ROOT_MEMBER_KEY
[[ -f $BOOTSTRAP_KEY ]] || $INSOLAR/bin/insolar -c gen_keys > $BOOTSTRAP_KEY

echo Creating node config
cp $INSOLAR/scripts/build/genesis/insolar.yaml insolar.cfg.yaml
echo "keyspath: $BOOTSTRAP_KEY" >> insolar.cfg.yaml

$INSOLAR/bin/insolard  \
                --config insolar.cfg.yaml \
                --genesis $INSOLAR/scripts/build/genesis/genesis.yaml  \
                --keyout certs \
                | tee log/genesis_output.log

tar -czf data.tgz data
tar -czf keys.tgz keys
