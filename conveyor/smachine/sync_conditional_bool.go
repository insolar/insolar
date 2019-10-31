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

package smachine

// ConditionalBool allows Acquire() call to pass through when current value is true
func NewConditionalBool(isOpen bool, name string) BoolConditionalLink {
	ctl := &boolConditionalSync{}
	ctl.controller.Init(name, &ctl.mutex, &ctl.controller)

	deps, _ := ctl.AdjustLimit(boolToConditional(isOpen), false)
	if len(deps) != 0 {
		panic("illegal state")
	}
	return BoolConditionalLink{ctl}
}

type BoolConditionalLink struct {
	ctl *boolConditionalSync
}

func boolToConditional(isOpen bool) int {
	if isOpen {
		return 1
	}
	return 0
}

func (v BoolConditionalLink) IsZero() bool {
	return v.ctl == nil
}

// Creates an adjustment that sets the given value when applied with SynchronizationContext.ApplyAdjustment()
// Can be applied multiple times.
func (v BoolConditionalLink) NewValue(isOpen bool) SyncAdjustment {
	if v.ctl == nil {
		panic("illegal state")
	}
	return SyncAdjustment{controller: v.ctl, adjustment: boolToConditional(isOpen), isAbsolute: true}
}

// Creates an adjustment that toggles the conditional when the adjustment is applied with SynchronizationContext.ApplyAdjustment()
// Can be applied multiple times.
func (v BoolConditionalLink) NewToggle() SyncAdjustment {
	if v.ctl == nil {
		panic("illegal state")
	}
	return SyncAdjustment{controller: v.ctl, adjustment: 1, isAbsolute: false}
}

func (v BoolConditionalLink) SyncLink() SyncLink {
	return NewSyncLink(v.ctl)
}

type boolConditionalSync struct {
	conditionalSync
}

func (p *boolConditionalSync) AdjustLimit(limit int, absolute bool) ([]StepLink, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	switch {
	case absolute:
		if limit > 0 {
			if p.controller.state > 0 {
				return nil, false
			}
			return p.setLimit(1)
		}
	case limit == 0:
		return nil, false
	case p.controller.state == 0: // flip-flop
		return p.setLimit(1)
	}
	return p.setLimit(0)
}
