package hash

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy/algorithmprovider"
)

type myProvider struct{}

func NewMyProvider() algorithmprovider.HashAlgorithmProvider {
	return &myProvider{}
}

func (*myProvider) Hash224bits() insolar.Hasher {
	return nil
}

func (*myProvider) Hash256bits() insolar.Hasher {
	return nil
}

func (*myProvider) Hash512bits() insolar.Hasher {
	return nil
}
