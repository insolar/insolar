{
  global: {
  },
  components: { 
      "jaeger":{
          port: 30686,
          jaeger_agent: {
              port: 6831
          }
      },
      "insolar": { 
          num_heavies: 1,
          num_lights: 2,
          num_virtuals: 2,
          hostname: "seed",
          domain: "bootstrap",
          tcp_transport_port: 7900,
          metrics_port: 8080,
          api_port: 19191,
          image: {
              image: "base",
              tag: "latest",
              image_pull_policy: "Never"
          }
      },
      "elk": {
        kibana_port: 30601,
        elasticsearch_port: 30200
      },
      "prometheus": {
          port: 30090
      },
      "utils":{
        insolar_conf : $.components.insolar,
        get_num_nodes : self.insolar_conf.num_heavies + self.insolar_conf.num_lights + self.insolar_conf.num_virtuals,
        host_template : self.insolar_conf.hostname + "-%d." + self.insolar_conf.domain + ":%d",
        id_to_node_type( id ) ::  if id < self.insolar_conf.num_heavies then "heavy_material" 
                                 else if id < self.insolar_conf.num_heavies + self.insolar_conf.num_lights then "light_material"
                                 else "virtual",

        local_log_volume_name: "node-log",
        local_log_volume() :: {
            "name": $.components.utils.local_log_volume_name,
            "hostPath": {
                "path": "/tmp/insolar_logs/",
                "type": "DirectoryOrCreate"
            }
        }
      }
  }
}
