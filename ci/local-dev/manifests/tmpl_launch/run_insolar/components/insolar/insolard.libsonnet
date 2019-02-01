{
    "versionmanager": {
        "minalowedversion": "v0.3.0"
    },
    "host": {
        "transport": {
            "protocol": "TCP",
            "address": "0.0.0.0:7900",
            "behindnat": false
        },
        "bootstraphosts": [],
        "isrelay": false,
        "infinitybootstrap": false,
        "timeout": 4
    },
    "service": {
        "service": {}
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
    "stats": {
        "listenaddress": "0.0.0.0:8080"
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
        "listenaddress": "0.0.0.0:8080"
    }
}

