local k = import 'ksonnet/ksonnet.beta.3/k.libsonnet';
local env = std.extVar("__ksonnet/environments");

local service = k.core.v1.service;
local servicePort = k.core.v1.service.mixin.spec.portsType;

local grafana_import = import 'grafana/grafana.libsonnet';

local full_grafana = grafana_import + {
                   _config+:: {
                     namespace: env.namespace,
                     "auth.anonymus": { enabled: true },
                     "security"
                   },
                 };


local grafana = full_grafana.grafana; 

k.core.v1.list.new(
  grafana.dashboardDefinitions +
  [
    grafana.dashboardSources,
    grafana.dashboardDatasources,
    grafana.deployment,
    grafana.serviceAccount,
    grafana.service +
    service.mixin.spec.withPorts(servicePort.newNamed('http', 3000, 'http') + servicePort.withNodePort(30910)) +
    service.mixin.spec.withType('NodePort'),
  ]
)
