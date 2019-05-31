local base_params = import '../params.libsonnet';
local params = std.mergePatch( base_params.components.insolar, std.extVar("__ksonnet/params").components.insolar );

local bootstrap_config = import 'bootstrap_config.libsonnet' ;
local statefull_set() = import 'statefull_set.libsonnet';
local genesis_insolard_conf() = import "insolard_genesis_config.libsonnet";
local insolard_conf() = import "insolard_config.libsonnet";
local k = import "k.libsonnet";
local pulsar() = import 'pulsar/pulsar_common.libsonnet';

local perisitant_claim() = {
  kind: "PersistentVolumeClaim",
  apiVersion: "v1",
  metadata: {
    name: "bootstrap-config",
    labels: {
      app: "bootstrap"
    }
  },
  spec: {
    accessModes: [
      "ReadWriteMany"
    ],
    resources: {
      requests: {
        storage: "2Gi"
      }
    }
  }
};

local service() = {
  apiVersion: "v1",
  kind: "Service",
  metadata: {
    name: "bootstrap",
    labels: {
      app: "bootstrap"
    }
  },
  spec: {
    ports: [
      {
        port: params.metrics_port,
        name: "metrics"
      },
      {
        port: params.tcp_transport_port,
        name: "network",
        protocol: "TCP"
      },
      {
        port: params.api_port,
        name: "api",
        protocol: "TCP"
      }
    ],
    clusterIP: "None",
    selector: {
      app: "bootstrap"
    }
  }
};

local configs() = {
  apiVersion: "v1",
  kind: "ConfigMap",
  metadata: {
    name: "node-config"
  },
  data:{
            "bootstrap.yaml": std.manifestYamlDoc(bootstrap_config.generate()),
            "insolar-genesis.yaml": std.manifestYamlDoc(genesis_insolard_conf()),
            "insolar.yaml": std.manifestYamlDoc(insolard_conf()),
    }

};

k.core.v1.list.new([configs(), service(), perisitant_claim(), statefull_set()])

