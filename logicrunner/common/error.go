// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
