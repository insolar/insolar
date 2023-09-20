package configuration

// Introspection holds introspection configuration.
type Introspection struct {
	// Addr specifies address where introspection server starts
	Addr string
}

// NewIntrospection creates new default configuration for introspection.
func NewIntrospection() Introspection {
	return Introspection{}
}
