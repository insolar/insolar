package common

import (
	"github.com/pkg/errors"
)

func RecoverError(msg string, recovered interface{}) error {
	if recovered == nil {
		return nil
	}
	return errors.WithStack(errors.WithMessage(errors.Errorf("%v", recovered), msg))
}
