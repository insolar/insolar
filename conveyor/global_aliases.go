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

package conveyor

import (
	"sync"

	"github.com/insolar/insolar/conveyor/smachine"
)

var _ smachine.SlotAliasRegistry = &GlobalAliases{}

type GlobalAliases struct {
	m sync.Map
}

func (p *GlobalAliases) UnpublishAlias(key interface{}) {
	p.m.Delete(key)
}

func (p *GlobalAliases) GetPublishedAlias(key interface{}) smachine.SlotLink {
	if v, ok := p.m.Load(key); ok {
		if link, ok := v.(smachine.SlotLink); ok {
			return link
		}
	}
	return smachine.SlotLink{}
}

func (p *GlobalAliases) PublishAlias(key interface{}, slot smachine.SlotLink) bool {
	_, loaded := p.m.LoadOrStore(key, slot)
	return !loaded
}
