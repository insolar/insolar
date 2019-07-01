package serialization

type ClaimHeader struct {
	TypeAndLength uint16 `insolar-transport:"header;[0-9]=length;[10-15]=header:ClaimType;group=Claims"` // [00-09] ByteLength [10-15] ClaimClass
	// actual payload
}

type GenericClaim struct {
	// ByteSize>=1
	ClaimHeader
	Payload []byte
}

type EmptyClaim struct {
	// ByteSize=1
	ClaimHeader `insolar-transport:"delimiter;ClaimType=0;length=header"`
}

type ClaimList struct {
	// ByteSize>=1
	Claims      []GenericClaim
	EndOfClaims EmptyClaim // ByteSize=1 - indicates end of claims
}
