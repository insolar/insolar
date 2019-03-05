
package sample

import (
	"github.com/insolar/insolar/conveyor/generator/common"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
    "github.com/insolar/insolar/conveyor/interfaces/adapter"
	"errors"
)

type SMFIDTestStateMachine struct {}

func (*SMFIDTestStateMachine) TID() common.ElType {
    return 1
}

func (*SMFIDTestStateMachine) s_First() common.ElState {
    return 1
}
func (*SMFIDTestStateMachine) s_Second() common.ElState {
    return 2
}

type SMRHTestStateMachine struct {
	cleanHandlers TestStateMachine
}

func SMRHTestStateMachineFactory() *common.StateMachine {
    m := SMRHTestStateMachine{
        cleanHandlers: &TestStateMachineImplementation{},
    }

    var x []common.State
    x = append(x, common.State{
            Transition: m.i_Init,
            TransitionFuture: m.if_Init,
            TransitionPast: m.ip_Init,
            ErrorState: m.es_Init,
            ErrorStateFuture: m.esf_Init,
            ErrorStatePast: m.esp_Init,
        },
        common.State{
            Migration: m.m_FirstSecond,
            MigrationFuturePresent: m.mfp_FirstSecond,
            Transition: m.t_First,
            TransitionFuture: m.tf_First,
            TransitionPast: m.tp_First,
            AdapterResponse: m.a_First,
            AdapterResponseFuture: m.af_First,
            AdapterResponsePast: m.ap_First,
            ErrorState: m.es_First,
            ErrorStateFuture: m.esf_First,
            ErrorStatePast: m.esp_First,
            AdapterResponseError: m.ea_First,
            AdapterResponseErrorFuture: m.eaf_First,
            AdapterResponseErrorPast: m.eap_First,
        },
        common.State{
            Migration: m.m_SecondThird,
            MigrationFuturePresent: m.mfp_SecondThird,
            Transition: m.t_Second,
            TransitionFuture: m.tf_Second,
            TransitionPast: m.tp_Second,
            AdapterResponse: m.a_Second,
            AdapterResponseFuture: m.af_Second,
            AdapterResponsePast: m.ap_Second,
            ErrorState: m.es_Second,
            ErrorStateFuture: m.esf_Second,
            ErrorStatePast: m.esp_Second,
            AdapterResponseError: m.ea_Second,
            AdapterResponseErrorFuture: m.eaf_Second,
            AdapterResponseErrorPast: m.eap_Second,
        },)

    return &common.StateMachine{
        Id: int(m.cleanHandlers.(TestStateMachine).TID()),
        States: x,
    }
}

// (index .States 0).GetTransitionName
func (s *SMRHTestStateMachine) i_Init(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanHandlers.i_Init(aInput, element.GetPayload())
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) if_Init(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanHandlers.if_Init(aInput, element.GetPayload())
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) ip_Init(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanHandlers.ip_Init(aInput, element.GetPayload())
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) es_Init(element slot.SlotElementHelper, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.es_Init(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state.ToInt()
}
func (s *SMRHTestStateMachine) esf_Init(element slot.SlotElementHelper, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.esf_Init(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state.ToInt()
}
func (s *SMRHTestStateMachine) esp_Init(element slot.SlotElementHelper, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.esp_Init(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state.ToInt()
}


func (s *SMRHTestStateMachine) t_First(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanHandlers.t_First(aInput, aPayload)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) tf_First(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanHandlers.tf_First(aInput, aPayload)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) tp_First(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanHandlers.tp_First(aInput, aPayload)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) m_FirstSecond(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanHandlers.m_FirstSecond(aInput, aPayload)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) mfp_FirstSecond(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanHandlers.mfp_FirstSecond(aInput, aPayload)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) es_First(element slot.SlotElementHelper, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.es_First(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state.ToInt()
}
func (s *SMRHTestStateMachine) esf_First(element slot.SlotElementHelper, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.esf_First(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state.ToInt()
}
func (s *SMRHTestStateMachine) esp_First(element slot.SlotElementHelper, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.esp_First(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state.ToInt()
}
func (s *SMRHTestStateMachine) a_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanHandlers.a_First(aInput, aPayload, aResponse)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) af_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanHandlers.af_First(aInput, aPayload, aResponse)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) ap_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanHandlers.ap_First(aInput, aPayload, aResponse)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) ea_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.ea_First(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state.ToInt()
}
func (s *SMRHTestStateMachine) eaf_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.eaf_First(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state.ToInt()
}
func (s *SMRHTestStateMachine) eap_First(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.eap_First(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state.ToInt()
}

func (s *SMRHTestStateMachine) t_Second(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanHandlers.t_Second(aInput, aPayload)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) tf_Second(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanHandlers.tf_Second(aInput, aPayload)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) tp_Second(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanHandlers.tp_Second(aInput, aPayload)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) m_SecondThird(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanHandlers.m_SecondThird(aInput, aPayload)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) mfp_SecondThird(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanHandlers.mfp_SecondThird(aInput, aPayload)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) es_Second(element slot.SlotElementHelper, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.es_Second(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state.ToInt()
}
func (s *SMRHTestStateMachine) esf_Second(element slot.SlotElementHelper, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.esf_Second(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state.ToInt()
}
func (s *SMRHTestStateMachine) esp_Second(element slot.SlotElementHelper, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.esp_Second(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state.ToInt()
}
func (s *SMRHTestStateMachine) a_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanHandlers.a_Second(aInput, aPayload, aResponse)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) af_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanHandlers.af_Second(aInput, aPayload, aResponse)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) ap_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().(TAR)
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanHandlers.ap_Second(aInput, aPayload, aResponse)
    return payload, state.ToInt(), err
}
func (s *SMRHTestStateMachine) ea_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.ea_Second(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state.ToInt()
}
func (s *SMRHTestStateMachine) eaf_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.eaf_Second(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state.ToInt()
}
func (s *SMRHTestStateMachine) eap_Second(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, uint32) {
    payload, state := s.cleanHandlers.eap_Second(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state.ToInt()
}









