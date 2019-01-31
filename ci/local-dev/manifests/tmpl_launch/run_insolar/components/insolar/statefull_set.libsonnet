local params = std.extVar("__ksonnet/params").components.insolar;

{
	"apiVersion": "apps/v1beta1",
	"kind": "StatefulSet",
	"metadata": {
		"name": "seed",
		"labels": {
			"app": "bootstrap"
		}
	},
	"spec": {
		"serviceName": "bootstrap",
		"replicas": params.num_heavies + params.num_lights + params.num_virtuals,
		"template": {
			"metadata": {
				"labels": {
					"app": "bootstrap"
				}
			},
			"spec": {
				"initContainers": [
					{
						"name": "init-bootstrap",
						"imagePullPolicy": "Never",
						"image": "base",
						"tty": true,
						"stdin": true,
						"command":[
                            "/bin/sh",  "-ec", importstr "launch.sh"
                            ],
						"env": [
							{
								"name": "HOME",
								"value": "/opt/insolar"
							},
							{
								"name": "INSOLAR_LEDGER_STORAGE_DATADIRECTORY",
								"value": "/opt/insolar/config/data"
							}
						],
						"volumeMounts": [
							{
								"name": "bootstrap-config",
								"mountPath": "/opt/insolar/config"
							},
							{
								"name": "code",
								"mountPath": "/tmp/code"
							},
							{
								"name": "seed-config",
								"mountPath": "/opt/insolar/config/insolar-genesis.yaml",
								"subPath": "insolar-genesis.yaml"
							},
							{
								"name": "seed-config",
								"mountPath": "/opt/insolar/config/genesis.yaml",
								"subPath": "genesis.yaml"
							},
							{
								"name": "work",
								"mountPath": "/opt/work"
							}
						]
					}
				],
				"containers": [
					{
						"name": "insgorund",
						"imagePullPolicy": "Never",
						"image": "base",
						"workingDir": "/opt/insolar",
						"tty": true,
						"stdin": true,
						"command": [
							"insgorund",
							"-l",
							"127.0.0.1:18181",
							"--rpc",
							"127.0.0.1:18182",
							"-d",
							"/tmp/code"
						],
						"env": [
							{
								"name": "HOME",
								"value": "/opt/insolar"
							}
						],
						"volumeMounts": [
							{
								"name": "work",
								"mountPath": "/opt/insolar"
							},
							{
								"name": "code",
								"mountPath": "/tmp/code"
							}
						]
					},
					{
						"name": "insolard",
						"imagePullPolicy": "Never",
						"image": "base",
						"workingDir": "/opt/insolar",
						"tty": true,
						"stdin": true,
						"command": [
							"/go/bin/insolard",
							"--config",
							"/opt/insolar/config/node-insolar.yaml"
						],
						"env": [
							{
								"name": "HOME",
								"value": "/opt/insolar"
							},
							{
								"name": "POD_NAME",
								"valueFrom": {
									"fieldRef": {
										"fieldPath": "metadata.name"
									}
								}
							},
							{
								"name": "POD_IP",
								"valueFrom": {
									"fieldRef": {
										"fieldPath": "status.podIP"
									}
								}
							},
							{
								"name": "INSOLAR_HOST_TRANSPORT_ADDRESS",
								"value": "$(POD_IP):7900"
							},
							{
								"name": "INSOLAR_APIRUNNER_ADDRESS",
								"value": "$(POD_IP):19191"
							}
						],
						"volumeMounts": [
							{
								"name": "work",
								"mountPath": "/opt/insolar"
							},
							{
								"name": "bootstrap-config",
								"mountPath": "/opt/bootstrap-config"
							},
							{
								"name": "code",
								"mountPath": "/tmp/code"
							},
							{
								"name": "node-config",
								"mountPath": "/opt/insolar/config/node-insolar.yaml",
								"subPath": "insolar.yaml"
							}
						]
					}
				],
				"volumes": [
					{
						"name": "bootstrap-config",
						"persistentVolumeClaim": {
							"claimName": "bootstrap-config"
						}
					},
					{
						"name": "code",
						"emptyDir": {}
					},
					{
						"name": "seed-config",
						"configMap": {
							"name": "seed-config"
						}
					},
					{
						"name": "node-config",
						"configMap": {
							"name": "node-config"
						}
					},
					{
						"name": "work",
						"emptyDir": {}
					}
				]
			}
		},
		"updateStrategy": {
			"type": "OnDelete"
		},
		"podManagementPolicy": "Parallel"
	}
}