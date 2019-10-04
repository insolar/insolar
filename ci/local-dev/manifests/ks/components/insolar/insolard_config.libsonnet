local base_params = import '../params.libsonnet';
local params = std.mergePatch( base_params.components, std.extVar("__ksonnet/params").components );
local insolar_params = params.insolar;

{
    "host": {
        "transport": {
            "protocol": "TCP",
            "address": "0.0.0.0:" + insolar_params.tcp_transport_port,
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
        "adapter": "zerolog"
    },
    "logicrunner": {
        "rpclisten": "127.0.0.1:18182",
        "builtin": {},
        "goplugin": {
            "runnerlisten": "127.0.0.1:18181"
        }
    },
    "apirunner": {
        "address": "127.0.0.1:" + insolar_params.api_port,
    },
    "pulsar": {
        "type": "tcp",
        "listenaddress": "0.0.0.0:8090",
        "nodesaddresses": []
    },
    "keyspath": "/opt/insolar/config/node-keys.json",
    "certificatepath": "/opt/insolar/config/node-cert.json",
    "metrics": {
        "listenaddress": "0.0.0.0:" + insolar_params.metrics_port,
    },
    "tracer": {
       	"jaeger": {
            "collectorendpoint": "",
            "agentendpoint": "jaeger-agent:" + params.jaeger.jaeger_agent.port,
            "probabilityrate": 1,
            "samplingrules": {}
        }
    }
}

