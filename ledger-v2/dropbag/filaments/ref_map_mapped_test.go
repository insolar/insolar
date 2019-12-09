//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package filaments

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMappedRefMap_LoadBucket(t *testing.T) {
	require.NotNil(t, emptyBucketMarker)

	//require.LessOrEqual(t, reference.LocalBinarySize, int(bucketKeyType.Size()))
	//require.LessOrEqual(t, reference.LocalBinarySize, bucketKeySize)
	//vf, _ := bucketKeyL1type.FieldByName("value")
	//require.Equal(t, reference.LocalBinarySize, int(vf.Offset))
}
