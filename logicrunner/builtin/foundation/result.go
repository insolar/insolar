// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package foundation

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
)

// Result struct that serialized and saved on ledger as result of request
type Result struct {
	// Error - logic error, this field is not NIL if error is created outside
	// contract code, for example "object not found"
	Error *Error
	// Returns - all return values of contract's method, whatever it returns,
	// including errors
	Returns []interface{}
}

// MarshalMethodErrorResult creates result data with all `returns` returned
// by contract
func MarshalMethodResult(returns ...interface{}) ([]byte, error) {
	result, err := insolar.Serialize(Result{Returns: returns})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't serialize method returns")
	}

	return result, nil
}

// MarshalMethodErrorResult creates result data with logic error for any method.
func MarshalMethodErrorResult(err error) ([]byte, error) {
	result, err := insolar.Serialize(Result{Error: &Error{S: err.Error()}})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't serialize method returns")
	}

	return result, nil
}

// UnmarshalMethodResult extracts method result into provided pointers to existing
// variables. You should know what types method returns and number of returned
// variables. If result contains logic error that was created outside of contract
// then method returns this error *Error. Errors of serialized returned as second
// return value. Example:
//
//     var i int
//     var contractError *foundation.Error
//     logicError, err := UnmarshalMethodResult(data, &i, &contractError)
//     if err != nil {
//         ... system error ...
//     }
//     if logicError != nil {
//         ... logic error created with MarshalMethodErrorResult ...
//     }
//     if contractError != nil {
//         ... contract returned error ...
//     }
//     ...
//
func UnmarshalMethodResult(data []byte, returns ...interface{}) (*Error, error) {
	res := Result{
		Returns: returns,
	}
	err := insolar.Deserialize(data, &res)
	if err != nil {
		return nil, errors.Wrap(err, "can't deserialize method result")
	}

	return res.Error, nil
}

// UnmarshalMethodResultSimplified is simplified version of UnmarshalMethodResult
// that finds *foundation.Error in `returns` and saves there logicError in case
// it's not empty. This works as we force all methods in contracts to return error.
// Top level logic error has priority over error returned by contract.
// Example:
//
//     var i int
//     var contractError *foundation.Error
//     err := UnmarshalMethodResultSimplified(data, &i, &contractError)
//     if err != nil {
//         ... system error ...
//     }
//     if contractError != nil {
//         ... logic error set by system of returned by contract ...
//     }
//     ...
func UnmarshalMethodResultSimplified(data []byte, returns ...interface{}) error {
	contractErr, err := UnmarshalMethodResult(data, returns...)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal result")
	}

	// this magic helper that injects logic error into one of returns, just sugar
	if contractErr != nil {
		found := false
		for i := 0; i < len(returns); i++ {
			if e, ok := returns[i].(**Error); ok {
				if e == nil {
					return errors.New("nil pointer in unmarshal")
				}
				*e = contractErr
				found = true
			}
		}
		if !found {
			return errors.New("no place for error in returns")
		}
	}

	return nil
}
