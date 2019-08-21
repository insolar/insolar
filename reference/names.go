//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
