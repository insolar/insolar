package configuration

// ServiceNetwork is configuration for ServiceNetwork.
type ServiceNetwork struct {
	CacheDirectory string
}

// NewServiceNetwork creates a new ServiceNetwork configuration.
func NewServiceNetwork() ServiceNetwork {
	return ServiceNetwork{
		CacheDirectory: "network_cache",
	}
}
