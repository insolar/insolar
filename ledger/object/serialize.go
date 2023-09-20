package object

import (
	"github.com/insolar/insolar/insolar"
)

// CalculateIDForBlob calculate id for blob with using current pulse number
func CalculateIDForBlob(scheme insolar.PlatformCryptographyScheme, pulseNumber insolar.PulseNumber, blob []byte) *insolar.ID {
	hasher := scheme.IntegrityHasher()
	_, err := hasher.Write(blob)
	if err != nil {
		panic(err)
	}
	return insolar.NewID(pulseNumber, hasher.Sum(nil))
}
