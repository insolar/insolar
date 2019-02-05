{
  global: {
    "utils":{
        insolar_conf :: $.components.insolar,
        get_num_nodes : self.insolar_conf.num_heavies + self.insolar_conf.num_lights + self.insolar_conf.num_virtuals,
        host_template : self.insolar_conf.hostname + "-%d." + self.insolar_conf.domain + ":" + self.insolar_conf.tcp_transport_port,
        id_to_node_type( id ) :  if id < self.insolar_conf.num_heavies then "heavy_material" 
                                 else if id < self.insolar_conf.num_heavies + self.insolar_conf.num_lights then "light_material"
                                 else "virtual",
      }
  },
  components: {
    "insolar": { 
      nodes:[
        { num_heavies: 1 },
        { num_lights: 5 },
        { num_virtuals: 4 }
      ],
      num_heavies: 1,
      num_lights: 5,
      num_virtuals: 4,
      hostname: "seed",
      domain: "bootstrap",
      tcp_transport_port: 7900,
      },
  }
}
