local params = std.extVar("__ksonnet/params").components;

local make_min_roles() = {
  virtual:  1,
  heavy_material: 1,
  light_material: 1
};

{
    generate_genesis() :: {

      // common fields
      root_keys_file: "/opt/insolar/config/root_member_keys.json",
      root_balance: 1000000000,
      majority_rule: 0,
      min_roles: make_min_roles(),
      pulsar_public_keys: [ "pulsar_public_key" ],

      // generating discovery_nodes
      local discovery_nodes_tmpl() = {
        host: params.utils.host_template,
        role: "%s",
        keys_file: "/opt/insolar/config/nodes/keys/%s-%d.json",
        cert_name: "%s-%d-cert.json"
      },

      discovery_nodes:
      [
         {
           host: discovery_nodes_tmpl().host % [ id ] ,
           keys_file: discovery_nodes_tmpl().keys_file % [ params.insolar.hostname, id ],
           cert_name: discovery_nodes_tmpl().cert_name % [  params.insolar.hostname, id ],

           role: discovery_nodes_tmpl().role %
             if id < params.insolar.num_heavies then [ "heavy_material" ]
             else if id < params.insolar.num_heavies + params.insolar.num_lights then [ "light_material" ]
             else [ "virtual" ]
         }
         for id in std.range(0, params.utils.get_num_nodes - 1)
      ]
    }
}

