local make_min_roles() = {
  virtual:  1,
  heavy_material: 1,
  light_material: 1
};

{
    generate_genesis( num_heavies = 5, num_lights=1, num_virtuals=1, hostname = "seed", domain = "bootstrap" ) :: {

      // common fields
      root_keys_file: "/opt/insolar/config/root_member_keys.json",
      root_balance: 1000000000,
      majority_rule: 0,
      min_roles: make_min_roles(),
      pulsar_public_keys: [ "pulsar_public_key" ],

      // generating discovery_nodes
      local discovery_nodes_tmpl() = {
        host: "%s-%d.%s:7900",
        role: "%s",
        keys_file: "/opt/insolar/config/nodes/%s-%d/keys.json",
        cert_name: "%s-%d-cert.json"
      },

      discovery_nodes:
      [
         {
           host: discovery_nodes_tmpl().host % [ hostname, id, domain ] ,
           keys_file: discovery_nodes_tmpl().keys_file % [ hostname, id ],
           cert_name: discovery_nodes_tmpl().cert_name % [ hostname, id ],

           role: discovery_nodes_tmpl().role %
             if id < num_heavies then [ "heavy_material" ]
             else if id < num_heavies + num_lights then [ "light_material" ]
             else [ "virtual" ]
         }
         for id in std.range(0, num_heavies + num_lights + num_virtuals - 1)
      ]

    }

}

