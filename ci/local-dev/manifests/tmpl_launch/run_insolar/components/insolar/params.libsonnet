{
  global: {
    // User-defined global parameters; accessible to all component and environments, Ex:
    // replicas: 4,
  },
  components: {
    "insolar":{ 
    	num_heavies: 1,
    	num_lights: 2,
    	num_virtuals: 2,
    	hostname: "seed",
    	domain: "bootstrap",	
    	},
    // Component-level parameters, defined initially from 'ks prototype use ...'
    // Each object below should correspond to a component in the components/ directory
  },
}
