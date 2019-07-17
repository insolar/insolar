package common

import (
	"github.com/tylerb/gls"
)

const glsSystemErrorKey = "systemError"

type SystemError interface {
	GetSystemError() error
	SetSystemError(err error)
}

type systemError struct{}

func NewSystemError() *systemError {
	return &systemError{}
}

func (h *systemError) GetSystemError() error {
	// SystemError means an error in the system (platform), not a particular contract.
	// For instance, timed out external call or failed deserialization means a SystemError.
	// In case of SystemError all following external calls during current method call return
	// an error and the result of the current method call is discarded (not registered).
	callContextInterface := gls.Get(glsSystemErrorKey)
	if callContextInterface == nil {
		return nil
	}
	return callContextInterface.(error)
}

func (h *systemError) SetSystemError(err error) {
	gls.Set(glsSystemErrorKey, err)
}
