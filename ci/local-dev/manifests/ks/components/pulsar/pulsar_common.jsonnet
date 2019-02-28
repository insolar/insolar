local k = import "k.libsonnet";


local pulsar_statefull_set() = import "pulsar_statefull_set.libsonnet";
local pulsar_conf() = {
	apiVersion: "v1",
	kind: "ConfigMap",
	metadata: {
		name: "pulsar-config"
	},
	data:{
		"pulsar.yaml": std.manifestYamlDoc(import "pulsar_config.libsonnet")
	}

};


k.core.v1.list.new([pulsar_statefull_set(), pulsar_conf()])