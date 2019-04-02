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

package component

import (
	"context"
)

// Initer interface provides method to init a component. During initialization components may NOT be ready. Only safe
// methods (e.g. not dependant on other components) can be called during initialization.
type Initer interface {
	Init(ctx context.Context) error
}

// Starter interface provides method to start a component.
type Starter interface {
	Start(ctx context.Context) error
}

// GracefulStopper interface provides method to end work with other components.
type GracefulStopper interface {
	GracefulStop(ctx context.Context) error
}

// Stopper interface provides method to stop a component.
type Stopper interface {
	Stop(ctx context.Context) error
}
