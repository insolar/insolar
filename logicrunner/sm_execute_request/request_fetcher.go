package sm_execute_request

// import (
// 	"github.com/insolar/insolar/conveyor/injector"
// 	"github.com/insolar/insolar/conveyor/smachine"
// 	"github.com/insolar/insolar/insolar"
// 	"github.com/insolar/insolar/logicrunner/s_artifact"
// )
//
// type RequestFetcher struct {
// 	smachine.StateMachineDeclTemplate
//
// 	ArtifactManager *s_artifact.ArtifactClientServiceAdapter
//
// 	// input arguments
// 	Object insolar.Reference
// 	Count  int
//
// 	// to pass between stages
// 	externalError error
// 	ignoredIDs    []insolar.ID
// 	requestIDs    []insolar.Reference
// }
//
// /* -------- Declaration ------------- */
//
// func (s *RequestFetcher) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
// 	return s.Init
// }
//
// func (s *RequestFetcher) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
// 	injector.MustInject(&s.ArtifactManager)
// }
//
// /* -------- Instance ------------- */
//
// func (s *RequestFetcher) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
// 	return s
// }
//
// func (s *RequestFetcher) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
// 	return ctx.Jump(s.stepFetchRequestIDs)
// }
//
// func (s *RequestFetcher) stepFetchRequestIDs(ctx smachine.ExecutionContext) smachine.StateUpdate {
// 	if s.Count == 0 {
// 		return ctx.Jump(s.stepStop)
// 	}
//
// 	var (
// 		goCtx   = ctx.GetContext()
// 		object  = s.Object
// 		ignored = s.ignoredIDs
// 	)
//
// 	return s.ArtifactManager.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
// 		requests, err := svc.GetPendings(goCtx, object, ignored)
// 		return func(ctx smachine.AsyncResultContext) {
// 			s.requestIDs = requests
// 			s.externalError = err
// 		}
// 	}).DelayedStart().Sleep().ThenJump(s.stepFetchRequest)
// }
//
// func (s *RequestFetcher) stepFetchRequest(ctx smachine.ExecutionContext) smachine.StateUpdate {
// 	if s.externalError != nil {
// 		return ctx.Jump(s.stepError)
// 	}
//
// }
//
// func (s *RequestFetcher) stepStop(ctx smachine.ExecutionContext) smachine.StateUpdate {
// 	return ctx.Stop()
// }
//
// func (s *RequestFetcher) stepError(ctx smachine.ExecutionContext) smachine.StateUpdate {
// 	return ctx.Error(s.externalError)
// }
