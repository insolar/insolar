local k = import "k.libsonnet";

local params = import '../params.libsonnet';
local insolar_params = params.components.insolar;

local service() = {
	"apiVersion": "v1",
	"kind": "Service",
	"metadata": {
		"name": "prometheus",
		"labels": {
			"app": "prometheus"
		}
	},
	"spec": {
		"type": "NodePort",
		"ports": [
			{
				"port": 9090,
				"nodePort": 30090,
				"name": "prometheus"
			}
		],
		"selector": {
			"app": "prometheus"
		}
	}
};

local pod() = {
	"apiVersion": "v1",
	"kind": "Pod",
	"metadata": {
		"name": "prometheus",
		"labels": {
			"app": "prometheus"
		}
	},
	"spec": {
		"containers": [
			{
				"name": "prometheus",
				"image": "prom/prometheus:v2.6.0",
				"tty": true,
				"stdin": true,
				"command": [
					"/bin/prometheus",
					"--config.file=/etc/prometheus/prometheus.yml"
				],
				"volumeMounts": [
					{
						"name": "prometheus-config",
						"mountPath": "/etc/prometheus/prometheus.yml",
						"subPath": "prometheus.yml"
					}
				]
			}
		],
		"volumes": [
			{
				"name": "prometheus-config",
				"configMap": {
					"name": "prometheus-config"
				}
			}
		]
	}
};

local get_typed_nodes( node_type ) = {
	tmp:: [  
			if params.global.utils.id_to_node_type( id ) == node_type then params.global.utils.host_template % [ id ]
			for id in std.range(0, params.global.utils.get_num_nodes - 1)
		 ],

    result : std.prune( self.tmp )

};

local config_map() = {
	"apiVersion": "v1",
	"kind": "ConfigMap",
	"metadata": {
		"name": "prometheus-config"
	},

	"data":{
		"prometheus.yml": std.manifestYamlDoc( {
							scrape_configs: [
								{
									job_name: "virtual",
									static_configs: [
										{ targets: get_typed_nodes( "virtual" ).result }
									],
								},
								{
									job_name: "heavy_material",
									static_configs: [
										{ targets: get_typed_nodes( "heavy_material" ).result }
									],
								},
								{
									job_name: "light_material",
									static_configs: [
										{ targets: get_typed_nodes( "light_material" ).result }
									],
								}
							 ]
						 } )
	}
};


k.core.v1.list.new([config_map(), service(), pod()])


