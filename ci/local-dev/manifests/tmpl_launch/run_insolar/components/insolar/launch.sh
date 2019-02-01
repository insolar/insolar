CONFIG_DIR=/opt/insolar/config
GENESIS_CONFIG=$CONFIG_DIR/genesis.yaml
NODES_DATA=$CONFIG_DIR/nodes
KEYS_DIR=$NODES_DATA/keys

NUM_NODES=$(fgrep '"host":' $GENESIS_CONFIG | grep -cv "#" )

ls -alhR /opt
if [ -f $CONFIG_DIR/bootstrap_keys.json ]
then
    echo skip generate
else
    echo generate bootstrap key
    insolar -c gen_keys > $CONFIG_DIR/bootstrap_keys.json
    echo generate root member key
    insolar -c gen_keys > $CONFIG_DIR/root_member_keys.json
    echo generate discovery node keys

    mkdir -p $KEYS_DIR
    for i in `seq 0 $((NUM_NODES-1))`
    do
        insolar -c gen_keys > $KEYS_DIR/seed-$i.json
    done

    echo generate genesis
    mkdir -p $NODES_DATA/certs
    mkdir -p $CONFIG_DIR/data
    insolard --config $CONFIG_DIR/insolar-genesis.yaml --genesis $GENESIS_CONFIG --keyout $NODES_DATA/certs
fi

echo next step
if [ -f /opt/work/config/node-cert.json ]
then
    echo skip work
else    
    echo copy genesis
    cp -R $CONFIG_DIR/data /opt/work/
    mkdir -p /opt/work/config
    cp $NODES_DATA/certs/$(ls $NODES_DATA/certs/ | grep $(hostname | sed 's/[^0-9]*//g')) /opt/work/config/node-cert.json
    cp $KEYS_DIR/$(hostname).json /opt/work/config/node-keys.json
fi
