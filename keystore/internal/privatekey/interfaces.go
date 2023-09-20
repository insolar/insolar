package privatekey

import (
	"crypto"
)

type Loader interface {
	Load(file string) (crypto.PrivateKey, error)
}
