local base_params = import '../params.libsonnet';
local params = std.mergePatch( base_params.components, std.extVar("__ksonnet/params").components );

local insolar_params = params.insolar;
local utils = params.utils;

{
	"pulsar": {
		"connectiontype": "tcp",
		"mainlisteneraddress": "0.0.0.0:58090",
		"storage": {
			"datadirectory": "/opt/insolar/pulsar",
			"txretriesonconflict": 0
		},
		"pulsetime": 10000,
		"receivingsigntimeout": 1000,
		"receivingnumbertimeout": 1000,
		"receivingvectortimeout": 1000,
		"receivingsignsforchosentimeout": 0,
		"neighbours": [],
		"numberofrandomhosts": 1,
		"numberdelta": 10,
		"distributiontransport": {
			"protocol": "TCP",
			"address": "0.0.0.0:58091",
		},
		"pulsedistributor": {
			"bootstraphosts": [
				utils.host_template % [ id , insolar_params.tcp_transport_port] for id in std.range(0, utils.get_num_nodes - 1)
			]
		}
	},
	"keyspath": "/opt/insolar/config/pulsar_keys.json",
	"log": {
		"level": "Debug",
		"adapter": "zerolog"
	}
}
