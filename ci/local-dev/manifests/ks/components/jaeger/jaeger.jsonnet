local k = import "k.libsonnet";

local base_params = import '../params.libsonnet';
local params = std.mergePatch( base_params.components, std.extVar("__ksonnet/params").components );

local service_jaeger() = {
	"apiVersion": "v1",
	"kind": "Service",
	"metadata": {
		"name": "jaeger",
		"labels": {
			"app": "jaeger"
		}
	},
	"spec": {
		"type": "NodePort",
		"ports": [
			{
				"port": 16686,
				"nodePort": params.jaeger.port,
				"name": "jaeger"
			}
		],
		"selector": {
			"app": "jaeger"
		}
	}
};

local service_jaeger_agent() = {
	"apiVersion": "v1",
	"kind": "Service",
	"metadata": {
		"name": "jaeger-agent",
		"labels": {
			"app": "jaeger-agent"
		}
	},
	"spec": {
		"ports": [
			{
				"port": params.jaeger.jaeger_agent.port,
				"protocol": "UDP",
				"name": "agent-compact"
			}
		],
		"selector": {
			"app": "jaeger"
		}
	}
};

local pod() = {
	"apiVersion": "v1",
	"kind": "Pod",
	"metadata": {
		"name": "jaeger",
		"labels": {
			"app": "jaeger"
		}
	},
	"spec": {
		"containers": [
			{
				"name": "jaeger",
				"image": "jaegertracing/all-in-one:1.8",
				"imagePullPolicy": "IfNotPresent",
				"tty": true,
				"stdin": true
			}
		],
		"env": [
			{
				"name": "SPAN_STORAGE_TYPE",
				"value": "elasticsearch"
			},
			{
				"name": "ES_SERVER_URLS",
				"value": "http://elk:" + params.elk.elasticsearch_port
			},
			{
				"name": "ES_TAGS_AS_FIELDS",
				"value": "true"
			}
		]
	}
};

k.core.v1.list.new([ service_jaeger(), service_jaeger_agent(), pod()])
