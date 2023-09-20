package consensuskit

// BftMajority function guarantees that (float(bftMajorityCount)/nodeCount > 2.0/3.0)	AND	(float(bftMajorityCount - 1)/nodeCount <= 2.0/3.0)
func BftMajority(nodeCount int) int {
	return nodeCount - BftMinority(nodeCount)
}

func BftMinority(nodeCount int) int {
	return (nodeCount - 1) / 3
}
