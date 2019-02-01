
{
	pulsar_statefull_set() :: import "pulsar_statefull_set.libsonnet",
	pulsar_conf() :: {
		apiVersion: "v1",
		kind: "ConfigMap",
		metadata: {
			name: "pulsar-config"
		},
		data:{
			"pulsar.yaml": std.manifestYamlDoc(import "pulsar_config.libsonnet")
		}

	}
}