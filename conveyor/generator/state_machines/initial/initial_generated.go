/*
*    Copyright 2019 Insolar Technologies
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

package initial

import (
	"context"
	"errors"

	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/statemachine"
)

func RawInitialPresentFactory(helpers *adapter.HelperCatalog) *statemachine.StateMachine {
	return &statemachine.StateMachine{
		ID: 3,
		States: []statemachine.State{
			{
				
				Transition: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(interface {})
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    state, payload := ParseInputEvent(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					
				},
				
			},
		},
	}
}

func RawInitialPastFactory(helpers *adapter.HelperCatalog) *statemachine.StateMachine {
	return &statemachine.StateMachine{
		ID: 3,
		States: []statemachine.State{
			{
				Transition: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(interface {})
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    state, payload := ParseInputEvent(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					
				},
				
			},
		},
	}
}

func RawInitialFutureFactory(helpers *adapter.HelperCatalog) *statemachine.StateMachine {
	return &statemachine.StateMachine{
		ID: 3,
		States: []statemachine.State{
			{
				
				Transition: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(interface {})
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    state, payload := ParseInputEvent(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					
				},
				
			},
		},
	}
}
