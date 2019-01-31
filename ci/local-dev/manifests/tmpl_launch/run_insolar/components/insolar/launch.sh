ls -alhR /opt
if [ -f /opt/insolar/config/bootstrap_keys.json ]
then
    echo skip generate
else
    echo generate bootstrap key
    insolar -c gen_keys > /opt/insolar/config/bootstrap_keys.json
    echo generate root member key
    insolar -c gen_keys > /opt/insolar/config/root_member_keys.json
    echo generate discovery node keys
    mkdir -p /opt/insolar/config/nodes/seed-0
    mkdir -p /opt/insolar/config/nodes/seed-1
    mkdir -p /opt/insolar/config/nodes/seed-2
    mkdir -p /opt/insolar/config/nodes/seed-3
    mkdir -p /opt/insolar/config/nodes/seed-4
    insolar -c gen_keys > /opt/insolar/config/nodes/seed-0/keys.json
    insolar -c gen_keys > /opt/insolar/config/nodes/seed-1/keys.json
    insolar -c gen_keys > /opt/insolar/config/nodes/seed-2/keys.json
    insolar -c gen_keys > /opt/insolar/config/nodes/seed-3/keys.json
    insolar -c gen_keys > /opt/insolar/config/nodes/seed-4/keys.json
    echo generate genesis
    mkdir -p /opt/insolar/config/nodes/certs
    mkdir -p /opt/insolar/config/data
    insolard --config /opt/insolar/config/insolar-genesis.yaml --genesis /opt/insolar/config/genesis.yaml --keyout /opt/insolar/config/nodes/certs
fi

echo next step
if [ -f /opt/work/config/node-cert.json ]; then
    echo skip work
else
    echo copy genesis
    cp -R /opt/insolar/config/data /opt/work/
    mkdir -p /opt/work/config
    cp /opt/insolar/config/nodes/certs/$(ls /opt/insolar/config/nodes/certs/ | grep $(hostname | sed 's/[^0-9]*//g')) /opt/work/config/node-cert.json
    cp /opt/insolar/config/nodes/$(hostname)/keys.json /opt/work/config/node-keys.json
fi
