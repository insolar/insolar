package proc

import (
	"testing"
)

// TestRunner is a helper structure. Its only purpose to run some code before and after main execution.
type TestRunner struct {
	before, after func()
	t             *testing.T
}

func NewTestRunner(t *testing.T) *TestRunner {
	return &TestRunner{t: t}
}

// Before sets a callback to call before each `Run`. Never place assertions here, only initialization.
func (s *TestRunner) Before(f func()) {
	s.before = f
}

// Before sets a callback to call after each `Run`. Never place assertions here, only initialization.
func (s *TestRunner) After(f func()) {
	s.after = f
}

// Executes `Before` callback then passed func then `After` callback. If you use closures with shared state, do not call
// t.Parallel() inside because it can mess up shared state.
func (s *TestRunner) Run(toRun func()) {
	if s.before != nil {
		s.before()
	}
	defer func() {
		if s.after != nil {
			s.after()
		}
	}()

	toRun()
}
