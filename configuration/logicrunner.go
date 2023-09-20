package configuration

// LogicRunner configuration
type LogicRunner struct {
	// PulseLRUSize - configuration of size of a pulse's cache
	PulseLRUSize int
}

// NewLogicRunner - returns default config of the logic runner
func NewLogicRunner() LogicRunner {
	return LogicRunner{
		PulseLRUSize: 100,
	}
}
