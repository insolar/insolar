package resolver

// PublicAddressResolver is network address resolver interface.
type PublicAddressResolver interface {

	// Resolve returns public network address from given internal address.
	Resolve(address string) (string, error)
}

// Resolve resolves public address
func Resolve(fixedPublicAddress, address string) (string, error) {
	var r PublicAddressResolver
	if fixedPublicAddress != "" {
		r = NewFixedAddressResolver(fixedPublicAddress)
	} else {
		r = NewExactResolver()
	}

	return r.Resolve(address)
}
