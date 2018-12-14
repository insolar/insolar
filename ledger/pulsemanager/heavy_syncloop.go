/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package pulsemanager

import (
	"context"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func (m *PulseManager) syncloop(ctx context.Context) {
	defer close(m.syncdone)

	inslog := inslogger.FromContext(ctx)

	for {
		inslog.Info("syncronization waiting new pulse signal")
		_, ok := <-m.gotpulse
		if !ok {
			inslog.Debug("stop is called, so we are should just stop syncronization loop")
			return
		}

		pulse, err := m.Current(ctx)
		if err != nil {
			err = errors.Wrap(err, "syncloop failed get current pulse")
			inslog.Error(err)
			continue
		}

		tosyncJetPulses := m.syncstates.unshiftJetPulses()
		if len(tosyncJetPulses) == 0 {
			continue
		}

		inslog.Infof("syncronization got next chunk of work")
		// TODO: reimplement retry logic - 14.Dec.2018 @nordicdyno
		var g errgroup.Group
		for _, tosync := range tosyncJetPulses {
			// it locks until Set call (LR.OnPulse) is finished
			// TODO: think how to parallelize replication per jet
			g.Go(func() error {
				err := m.HeavySync(ctx, *pulse, tosync.jet, tosync.pn, false)
				if err == nil {
					m.db.SetReplicatedPulse(ctx, tosync.jet, tosync.pn)
				}
				return err
			})
		}
		if err := g.Wait(); err != nil {
			inslog.Error(err)
			continue
		}
	}
}
