package record

// LifelineIndex represents meta information for record object
type LifelineIndex struct {
	LatestStateID ID
	LatestStateType TypeID
	AppendIDs []ID
}
