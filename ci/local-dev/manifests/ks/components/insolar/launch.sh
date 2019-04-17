CONFIG_DIR=/opt/insolar/config
GENESIS_CONFIG=$CONFIG_DIR/genesis.yaml
NODES_DATA=$CONFIG_DIR/nodes
DISCOVERY_KEYS=$CONFIG_DIR/discovery
CERTS_KEYS=$CONFIG_DIR/certs

NUM_NODES=$(fgrep '"host":' $GENESIS_CONFIG | grep -cv "#" )

ls -alhR /opt
if [ "$HOSTNAME" = seed-0 ] && ! ( test -e /opt/insolar/config/finished )
then
    echo "generate bootstrap key"
    insolar gen-key-pair > $CONFIG_DIR/bootstrap_keys.json

    echo "generate root member key"
    insolar gen-key-pair > $CONFIG_DIR/root_member_keys.json

    echo "generate genesis"
    mkdir -vp $NODES_DATA
    mkdir -vp $CERTS_KEYS
    mkdir -vp $CONFIG_DIR/data
    mkdir -vp $DISCOVERY_KEYS
    insolard --config $CONFIG_DIR/insolar-genesis.yaml --genesis $GENESIS_CONFIG --keyout $CERTS_KEYS
    touch /opt/insolar/config/finished
else
    while ! (/usr/bin/test -e /opt/insolar/config/finished)
    do
        echo "Waiting for genesis ... ( sleep 5 sec )"
        sleep 5s
    done
fi

echo next step
if [ -f /opt/work/config/node-cert.json ]
then
    echo "skip work"
else    
    echo "copy genesis"
    cp -vR $CONFIG_DIR/data /opt/work/
    mkdir -vp /opt/work/config
    cp -v $CERTS_KEYS/$(hostname | awk -F'-' '{ printf "seed-%d-cert.json", $2 }')  /opt/work/config/node-cert.json
    cp -v $DISCOVERY_KEYS/$(hostname | awk -F'-' '{ printf "key-%02d.json", $2 }')  /opt/work/config/node-keys.json
fi
