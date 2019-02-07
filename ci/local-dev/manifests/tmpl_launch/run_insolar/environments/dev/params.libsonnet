local params = std.extVar("__ksonnet/params");
local globals = import "globals.libsonnet";
local envParams = params + {
  components +: {
    "insolar.insolar"+:{
    image+: {
              image: "registry.insolar.io/insolard",
              tag: "v0.7.6",
              image_pull_policy: "IfNotPresent"
          }
        },
    "pulsar.insolar"+:{
    image+: {
              image: "registry.insolar.io/insolard",
              tag: "v0.7.6",
              image_pull_policy: "IfNotPresent"
          }
        }
  },
};

{
  components: {
    [x]: envParams.components[x] + globals, for x in std.objectFields(envParams.components)
  },
}
