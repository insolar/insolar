// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package args

import (
	"fmt"
	"time"
)

func LazyStr(fn func() string) fmt.Stringer {
	return &lazyStringer{fn}
}

func LazyFmt(format string, a ...interface{}) fmt.Stringer {
	return &lazyStringer{func() string {
		return fmt.Sprintf(format, a...)
	}}
}

func LazyTimeFmt(format string, t time.Time) fmt.Stringer {
	return &lazyStringer{func() string {
		return t.Format(format)
	}}
}

type lazyStringer struct {
	fn func() string
}

func (v *lazyStringer) String() string {
	return v.fn()
}
