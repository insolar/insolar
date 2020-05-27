package load

import (
	"log"

	"github.com/insolar/loadgen"
)

func AttackerFromName(name string) loadgen.Attack {
	switch name {
	case "get_records":
		return loadgen.WithMonitor(new(GetRecordsAttack))
	case "get_pulses":
		return loadgen.WithMonitor(new(GetPulsesAttack))
	default:
		log.Fatalf("unknown attacker type: %s", name)
		return nil
	}
}
