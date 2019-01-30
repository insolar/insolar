
local utils = import 'utils.libsonnet' ;
local statefull_set() = import 'statefull_set.libsonnet';

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
        port: 8080,
        name: "prometheus"
      },
      {
        port: 7900,
        name: "network",
        protocol: "TCP"
      },
      {
        port: 19191,
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
    name: "seed-config"
  },
  data:{
            "genesis.yaml": std.manifestYamlDoc(utils.generate_genesis()),
            "insolar.yaml": importstr "insolard.yaml",
            "insolar-genesis.yaml": importstr "insolard-genesis.yaml"
    }

};






//configs()
statefull_set()


