
local base_params = import '../params.libsonnet';
local params = std.mergePatch( base_params.components, std.extVar("__ksonnet/params").components );

local utils = params.utils;
local image_params = params.insolar.image;

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
		"replicas": utils.get_num_nodes,
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
						"imagePullPolicy": image_params.image_pull_policy,
						"image": image_params.image + ":" + image_params.tag,
						"tty": true,
						"stdin": true,
						"command": [
							"/bin/sh", "-ec", importstr "launch.sh"
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
								"name": "node-config",
								"mountPath": "/opt/insolar/config/insolar-genesis.yaml",
								"subPath": "insolar-genesis.yaml"
							},
							{
								"name": "node-config",
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
						"imagePullPolicy": image_params.image_pull_policy,
						"image": image_params.image + ":" + image_params.tag,
						"workingDir": "/opt/insolar",
						"tty": true,
						"stdin": true,
						launch_cmd :: "/go/bin/insgorund -l 127.0.0.1:18181 --rpc 127.0.0.1:18182 -d /tmp/code 2>&1",
						"command": [
							"bash",
							"-c",
							if params.insolar.local_launch == true
							then
								self.launch_cmd + " | tee /logs/$(POD_NAME).insolard.log 2>&1"
							else
								self.launch_cmd

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
							},
							{
								"name": "node-log",
								"mountPath": "/logs"
							}
						]
					},
					{
						"name": "insolard",
						"imagePullPolicy": image_params.image_pull_policy,
						"image": image_params.image + ":" + image_params.tag,
						"workingDir": "/opt/insolar",
						"tty": true,
						"stdin": true,
						launch_cmd :: "/go/bin/insolard --config /opt/insolar/config/node-insolar.yaml --trace 2>&1",
						"command": [
							"bash",
							"-c",
							if params.insolar.local_launch == true
							then
								self.launch_cmd + " | tee /logs/$(POD_NAME).insolard.log 2>&1"
							else
								self.launch_cmd

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
								"value": "$(POD_IP):" + params.insolar.api_port
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
							},
							{
								"name": utils.local_log_volume_name,
								"mountPath": "/logs"
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
						"name": "node-config",
						"configMap": {
							"name": "node-config"
						}
					},
					{
						"name": "work",
						"emptyDir": {}
					},
					utils.local_log_volume()
				]
			}
		},
		"updateStrategy": {
			"type": "OnDelete"
		},
		"podManagementPolicy": "Parallel"
	}
}
