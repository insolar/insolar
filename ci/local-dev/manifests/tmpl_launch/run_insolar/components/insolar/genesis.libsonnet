local params = import '../params.libsonnet';

local make_min_roles() = {
  virtual:  1,
  heavy_material: 1,
  light_material: 1
};

{
    // It generates nodes in particular order: 1) heavy_material 2) light_material 3) virtual
    generate_genesis() :: {

      // common fields
      root_keys_file: "/opt/insolar/config/root_member_keys.json",
      root_balance: 1000000000,
      majority_rule: 0,
      min_roles: make_min_roles(),
      pulsar_public_keys: [ "pulsar_public_key" ],

      // generating discovery_nodes
      local discovery_nodes_tmpl() = {
        host: params.global.utils.host_template,
        role: "%s",
        keys_file: "/opt/insolar/config/nodes/keys/%s-%d.json",
        cert_name: "%s-%d-cert.json"
      },

      discovery_nodes:
      [
         {
           insolar_params :: params.components.insolar,
           host: discovery_nodes_tmpl().host % [ id ] ,
           keys_file: discovery_nodes_tmpl().keys_file % [ self.insolar_params.hostname, id ],
           cert_name: discovery_nodes_tmpl().cert_name % [ self.insolar_params.hostname, id ],

           role: discovery_nodes_tmpl().role % params.global.utils.id_to_node_type( id ), 
         }
         for id in std.range(0, params.global.utils.get_num_nodes - 1)
      ]
    }
}

