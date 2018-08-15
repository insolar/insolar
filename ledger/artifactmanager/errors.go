package artifactmanager

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidRef        = errors.New("invalid reference")
	ErrClassDeactivated  = errors.New("class is deactivated")
	ErrObjectDeactivated = errors.New("object is deactivated")
	ErrInconsistentIndex = errors.New("inconsistent index")
	ErrWrongObject       = errors.New("provided object is not and instance of provided class")
)
