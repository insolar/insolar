// Copyright 2020 Insolar Network Ltd.
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

package executor

import (
	"context"
	fmt "fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

func (j *JetInfo) updateSplit(split bool) error {
	if !j.IsSplitSet {
		j.Split = split
		j.IsSplitSet = true
	} else if j.Split != split {
		return errors.New(fmt.Sprintf("try to change split from %t to %t ", j.Split, split))
	}
	return nil
}

func (j *JetInfo) addDrop(newJetID insolar.JetID, split bool) error {
	if j.DropConfirmed {
		return errors.New("addDrop. try to rewrite drop confirmation. existing: " + j.JetID.DebugString() +
			", new: " + newJetID.DebugString())
	}

	if err := j.updateSplit(split); err != nil {
		return errors.Wrap(err, "updateSplit return error")
	}

	j.DropConfirmed = true
	j.JetID = newJetID

	return nil
}

func (j *JetInfo) checkIncomingHot(incomingJetID insolar.JetID) error {
	if len(j.HotConfirmed) >= 2 {
		return errors.New("num hot confirmations exceeds 2. existing: " + insolar.JetIDCollection(j.HotConfirmed).DebugString() +
			", new: " + incomingJetID.DebugString())
	}

	if len(j.HotConfirmed) == 1 && j.HotConfirmed[0].Equal(incomingJetID) {
		return errors.New("try add already existing hot confirmation: " + incomingJetID.DebugString())
	}

	return nil
}

func (j *JetInfo) addBackup() {
	j.BackupConfirmed = true
}

func (j *JetInfo) addHot(newJetID insolar.JetID, parentID insolar.JetID, split bool) error {
	err := j.checkIncomingHot(newJetID)
	if err != nil {
		return errors.Wrap(err, "incorrect incoming jet")
	}

	j.HotConfirmed = append(j.HotConfirmed, newJetID)
	j.JetID = parentID
	if err := j.updateSplit(split); err != nil {
		return errors.Wrap(err, "updateSplit return error")
	}

	return nil
}

func (j *JetInfo) isConfirmed(ctx context.Context, checkBackup bool) bool {
	if checkBackup && !j.BackupConfirmed {
		return false
	}

	if !j.DropConfirmed {
		return false
	}

	if len(j.HotConfirmed) == 0 {
		return false
	}

	if !j.IsSplitSet {
		inslogger.FromContext(ctx).Error("IsSplitJet must be set before calling for isConfirmed")
		return false
	}

	if !j.Split {
		return j.HotConfirmed[0].Equal(j.JetID)
	}

	if len(j.HotConfirmed) != 2 {
		return false
	}

	parentFirst := jet.Parent(j.HotConfirmed[0])
	parentSecond := jet.Parent(j.HotConfirmed[1])

	return parentFirst.Equal(parentSecond) && parentSecond.Equal(j.JetID)
}
