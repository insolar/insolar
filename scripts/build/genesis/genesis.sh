#!/bin/sh

echo ==============================
echo Making genesis
echo ==============================

pushd `dirname $0`
INSOLAR=`git rev-parse --show-toplevel`
popd

BASE_DIR=$INSOLAR/scripts/insolard

dirs="data keys certs bin log"
for d in $dirs; do rm -rf $d; done
for d in $dirs; do mkdir $d; done
rm -f insolar.cfg.yaml

#pushd $INSOLAR
#make build
#popd

ROOT_MEMBER_KEY=keys/root_member.key.json
BOOTSTRAP_KEY=keys/bootstrap.key.json
[ -f $ROOT_MEMBER_KEY ] || $INSOLAR/bin/insolar -c gen_keys > $ROOT_MEMBER_KEY
[ -f $BOOTSTRAP_KEY ] || $INSOLAR/bin/insolar -c gen_keys > $BOOTSTRAP_KEY


echo Creating node config
cp $INSOLAR/scripts/build/genesis/insolar.yaml insolar.cfg.yaml
echo "keyspath: $BOOTSTRAP_KEY" >> insolar.cfg.yaml



$INSOLAR/bin/insolard  \
                --config insolar.cfg.yaml \
                --genesis $INSOLAR/scripts/build/genesis/genesis.yaml  \
                --keyout certs \
                | tee log/genesis_output.log

