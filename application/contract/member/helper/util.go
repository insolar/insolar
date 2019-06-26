package helper

import "github.com/insolar/insolar/insolar"

// Contains tells whether a contains x.
func Contains(a []insolar.Reference, x insolar.Reference) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
