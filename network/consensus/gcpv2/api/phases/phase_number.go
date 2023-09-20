package phases

type Number uint8

const (
	Phase0 Number = iota
	Phase1
	Phase2
	Phase3
	Phase4
	upperPhaseNumber
)

const Count = int(upperPhaseNumber)
const Mask = 1<<upperPhaseNumber - 1

type Type uint8

const (
	TypePulse Type = iota
	TypeSoloState
	TypeNeighborsState
	TypeStateVectors
)

var phaseTypes = [upperPhaseNumber]Type{TypePulse, TypeSoloState, TypeNeighborsState, TypeStateVectors, TypeStateVectors}

func (v Number) GetPhaseType() Type {
	return phaseTypes[v]
}

type Flag uint8

const (
	FlagUnique Flag = iota
)
