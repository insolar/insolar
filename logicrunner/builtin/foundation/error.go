package foundation

// Error elementary string based error struct satisfying builtin error interface
//    foundation.Error{"some err"}
type Error struct {
	S string
}

// Error returns error in string format
func (e *Error) Error() string {
	return e.S
}
