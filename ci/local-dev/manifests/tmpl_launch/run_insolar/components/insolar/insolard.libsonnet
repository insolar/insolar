local base_params = import '../params.libsonnet';
local params = std.mergePatch( base_params.components.insolar, std.extVar("__ksonnet/params").components.insolar );

{
    "versionmanager": {
        "minalowedversion": "v0.3.0"
    },
    "host": {
        "transport": {
            "protocol": "TCP",
            "address": "0.0.0.0:" + params.tcp_transport_port,
            "behindnat": false
        },
        "bootstraphosts": [],
        "isrelay": false,
        "infinitybootstrap": false,
        "timeout": 4
    },
    "ledger": {
        "storage": {
            "datadirectory": "/opt/insolar/data",
            "txretriesonconflict": 3
        }
    },
    "log": {
        "level": "Debug",
        "adapter": "logrus"
    },
    "logicrunner": {
        "rpclisten": "127.0.0.1:18182",
        "builtin": {},
        "goplugin": {
            "runnerlisten": "127.0.0.1:18181"
        }
    },
    "apirunner": {
        "address": "127.0.0.1:19191"
    },
    "pulsar": {
        "type": "tcp",
        "listenaddress": "0.0.0.0:8090",
        "nodesaddresses": []
    },
    "keyspath": "/opt/insolar/config/node-keys.json",
    "certificatepath": "/opt/insolar/config/node-cert.json",
    "metrics": {
        "listenaddress": "0.0.0.0:" + params.metrics_port
    },
    "tracer": {
       	"jaeger": {
            "collectorendpoint": "",
            "agentendpoint": "jaeger-agent:6831",
            "probabilityrate": 1,
            "samplingrules": {}
        }
    }
}

