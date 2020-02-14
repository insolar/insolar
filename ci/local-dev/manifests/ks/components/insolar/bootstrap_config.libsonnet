local base_params = import '../params.libsonnet';
local params = std.mergePatch(base_params.components, std.extVar('__ksonnet/params').components);
local utils = params.utils;
local insolar_params = params.insolar;

local make_min_roles() = {
  virtual: 1,
  heavy_material: 1,
  light_material: 1,
};

{
  // It generates nodes in particular order: 1) heavy_material 2) light_material 3) virtual
  generate():: {

    // common fields
    members_keys_dir: '/opt/insolar/config/',
    node_keys_dir: '/opt/insolar/config/nodes',
    discovery_keys_dir: '/opt/insolar/config/discovery',
    heavy_genesis_config_file: "/opt/insolar/config/heavy_genesis.json",
    heavy_genesis_plugins_dir: "/opt/insolar/plugins",
    contracts: {
      insgocc: "/go/bin/insgocc",
      outdir: "/opt/insolar/config/plugins",
    },
    root_balance: 1000000000,
    majority_rule: 0,
    min_roles: make_min_roles(),
    pulsar_public_keys: ['pulsar_public_key'],

    // generating discovery_nodes
    local discovery_nodes_tmpl() = {
      host: utils.host_template,
      role: '%s',
      cert_name: '%s-%d-cert.json',
      key_name: '%s-%d-key.json',
    },

    discovery_nodes:
      [
        {
          host: discovery_nodes_tmpl().host % [id, insolar_params.tcp_transport_port],
          cert_name: discovery_nodes_tmpl().cert_name % [insolar_params.hostname, id],
          key_name:  discovery_nodes_tmpl().key_name  % [insolar_params.hostname, id],

          role: discovery_nodes_tmpl().role % utils.id_to_node_type(id),
        }
        for id in std.range(0, utils.get_num_nodes - 1)
      ],
  },
}
