// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package reference

import (
	"regexp"
)

const LegacyDomainName = "11111111111111111111111111111111"
const RecordDomainName = "record"

var regexObjectName = regexp.MustCompile(`^[[:alpha:]][[:alnum:]]*$`)
var regexDomainName = regexp.MustCompile(`^[[:alpha:]][[:alnum:]]*(\.[[:alpha:]][[:alnum:]]*)*$`)

func IsReservedName(domainName string) bool {
	return domainName == RecordDomainName || domainName == LegacyDomainName
}

func IsValidDomainName(domainName string) bool {
	return regexDomainName.MatchString(domainName)
}

func IsValidObjectName(objectName string) bool {
	return regexObjectName.MatchString(objectName)
}
