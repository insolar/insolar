local base_params = import '../params.libsonnet';
local params = std.mergePatch( base_params.components, std.extVar("__ksonnet/params").components );

local image_params = params.insolar.image;

{
	"apiVersion": "apps/v1beta1",
	"kind": "StatefulSet",
	"metadata": {
		"name": "pulsar",
		"labels": {
			"app": "pulsar"
		}
	},
	"spec": {
		"serviceName": "pulsar",
		"replicas": 1,
		"template": {
			"metadata": {
				"labels": {
					"app": "pulsar"
				}
			},
			"spec": {
				"initContainers": [
					{
						"name": "init-register",
						"imagePullPolicy": image_params.image_pull_policy,
						"image": image_params.image + ":" + image_params.tag,
						"tty": true,
						"stdin": true,
						"command": [
							"/bin/sh",
							"-c",
							"/go/bin/insolar gen-key-pair > /opt/insolar/config/bootstrap_keys.json;"
						],
						"env": [
							{
								"name": "HOME",
								"value": "/opt/insolar"
							}
						],
						"volumeMounts": [
							{
								"name": "config",
								"mountPath": "/opt/insolar/config"
							}
						]
					}
				],
				"containers": [
					{
						"name": "pulsar",
						"imagePullPolicy": image_params.image_pull_policy,
						"image": image_params.image + ":" + image_params.tag,
						"workingDir": "/opt/insolar",
						"tty": true,
						"stdin": true,
						"command": [
							"/go/bin/pulsard",
							"-c",
							"/opt/insolar/config/pulsar.yaml"
						],
						"env": [
							{
								"name": "HOME",
								"value": "/opt/insolar"
							},
							{
								"name": "INSOLAR_KEYSPATH",
								"value": "/opt/insolar/config/bootstrap_keys.json"
							},
							{
								"name": "INSOLAR_PULSAR_STORAGE_DATADIRECTORY",
								"value": "/opt/insolar/pulsar"
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
								"name": "INSOLAR_PULSAR_MAINLISTENERADDRESS",
								"value": "$(POD_IP):58090"
							},
							{
								"name": "INSOLAR_PULSAR_DISTRIBUTIONTRANSPORT_ADDRESS",
								"value": "$(POD_IP):58091"
							}
						],
						"resources": {
							"requests": {
								"cpu": "300m",
								"memory": "200M"
							}
						},
						"volumeMounts": [
							{
								"name": "config",
								"mountPath": "/opt/insolar/config"
							},
							{
								"name": "pulsar",
								"mountPath": "/opt/insolar/pulsar"
							},
							{
								"name": "code",
								"mountPath": "/tmp/code"
							},
							{
								"name": "pulsar-config",
								"mountPath": "/opt/insolar/config/pulsar.yaml",
								"subPath": "pulsar.yaml"
							}
						]
					}
				],
				"volumes": [
					{
						"name": "config",
						"emptyDir": {}
					},
					{
						"name": "pulsar",
						"emptyDir": {}
					},
					{
						"name": "code",
						"emptyDir": {}
					},
					{
						"name": "pulsar-config",
						"configMap": {
							"name": "pulsar-config"
						}
					}
				],
				"imagePullSecrets": [
					{
						"name": "registry-insolar-io"
					}
				]
			}
		},
		"updateStrategy": {
			"type": "RollingUpdate"
		},
		"podManagementPolicy": "Parallel"
	}
}
