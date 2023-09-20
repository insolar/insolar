package coreapi

type PacketVerifyFlags uint32

const DefaultVerify PacketVerifyFlags = 0

const (
	SkipVerify PacketVerifyFlags = 1 << iota
	RequireStrictVerify
	AllowUnverified
	SuccessfullyVerified
)
