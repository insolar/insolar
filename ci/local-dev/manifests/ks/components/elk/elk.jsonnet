local k = import "k.libsonnet";

local base_params = import '../params.libsonnet';
local params = std.mergePatch( base_params.components, std.extVar("__ksonnet/params").components );
local utils = params.utils;
local elk_params = params.elk;

local config_map() = {
	"apiVersion": "v1",
	"kind": "ConfigMap",
	"metadata": {
		"name": "elk-config"
	},
	"data": {
		"logstash.conf": |||
			input {
				file {
					path => "/logs/*.log"
					start_position => "beginning"
				}
			}

			output {
				elasticsearch {
					hosts => ["elasticsearch:9200"]
				}
			}
		|||
	}
};

local service() = {
	"apiVersion": "v1",
	"kind": "Service",
	"metadata": {
		"name": "elk",
		"labels": {
			"app": "elk"
		}
	},
	"spec": {
		"type": "NodePort",
		"ports": [
			{
				"port": 5601,
				"nodePort": elk_params.kibana_port,
				"name": "kibana"
			},
			{
				"port": 9200,
				"nodePort": elk_params.elasticsearch_port,
				"name": "elasticsearch"
			}
		],
		"selector": {
			"app": "elk"
		}
	}
};

local pod() = {
	"apiVersion": "v1",
	"kind": "Pod",
	"metadata": {
		"name": "elk",
		"labels": {
			"app": "elk"
		}
	},
	"spec": {
		"hostname": "elasticsearch",
		"containers": [
			{
				"name": "elasticsearch",
				"image": "docker.elastic.co/elasticsearch/elasticsearch:6.5.4",
				"tty": true,
				"stdin": true,
				"env": [
					{
						"name": "ES_JAVA_OPTS",
						"value": "-Xmx256m -Xms256m"
					},
					{
						"name": "discovery.type",
						"value": "single-node"
					},
					{
						"name": "xpack.security.enabled",
						"value": "false"
					}
				]
			},
			{
				"name": "kibana",
				"image": "docker.elastic.co/kibana/kibana:6.5.4",
				"tty": true,
				"stdin": true
			},
			{
				"name": "logstash",
				"image": "docker.elastic.co/logstash/logstash:6.6.0",
				"tty": true,
				"stdin": true,
				"volumeMounts": [
					{
						"name": "elk-config",
						"mountPath": "/usr/share/logstash/pipeline/logstash2.conf",
						"subPath": "logstash.conf"
					},
					{
						"name": utils.local_log_volume_name,
						"mountPath": "/logs"
					}
				]
			}
		],
		"volumes": [
			utils.local_log_volume(),
			{
				"name": "elk-config",
				"configMap": {
					"name": "elk-config"
				}
			}
		]
	}
};

k.core.v1.list.new([config_map(), service(), pod()])

