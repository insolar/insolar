// Code generated by "stringer -type=Type"; DO NOT EDIT.

package payload

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TypeUnknown-0]
	_ = x[TypeMeta-1]
	_ = x[TypeError-2]
	_ = x[TypeID-3]
	_ = x[TypeIDs-4]
	_ = x[TypeJet-5]
	_ = x[TypeState-6]
	_ = x[TypeGetObject-7]
	_ = x[TypePassState-8]
	_ = x[TypeIndex-9]
	_ = x[TypePass-10]
	_ = x[TypeGetCode-11]
	_ = x[TypeCode-12]
	_ = x[TypeSetCode-13]
	_ = x[TypeSetIncomingRequest-14]
	_ = x[TypeSetOutgoingRequest-15]
	_ = x[TypeSagaCallAcceptNotification-16]
	_ = x[TypeGetFilament-17]
	_ = x[TypeGetRequest-18]
	_ = x[TypeRequest-19]
	_ = x[TypeFilamentSegment-20]
	_ = x[TypeSetResult-21]
	_ = x[TypeActivate-22]
	_ = x[TypeRequestInfo-23]
	_ = x[TypeGetRequestInfo-24]
	_ = x[TypeGotHotConfirmation-25]
	_ = x[TypeDeactivate-26]
	_ = x[TypeUpdate-27]
	_ = x[TypeHotObjects-28]
	_ = x[TypeResultInfo-29]
	_ = x[TypeGetPendings-30]
	_ = x[TypeHasPendings-31]
	_ = x[TypePendingsInfo-32]
	_ = x[TypeReplication-33]
	_ = x[TypeGetJet-34]
	_ = x[TypeAbandonedRequestsNotification-35]
	_ = x[TypeGetLightInitialState-36]
	_ = x[TypeLightInitialState-37]
	_ = x[TypeGetIndex-38]
	_ = x[TypeSearchIndex-39]
	_ = x[TypeUpdateJet-40]
	_ = x[TypeReturnResults-41]
	_ = x[TypeCallMethod-42]
	_ = x[TypeExecutorResults-43]
	_ = x[TypePendingFinished-44]
	_ = x[TypeAdditionalCallFromPreviousExecutor-45]
	_ = x[TypeStillExecuting-46]
	_ = x[TypeErrorResultExitsts-47]
	_ = x[_latestType-48]
}

const _Type_name = "TypeUnknownTypeMetaTypeErrorTypeIDTypeIDsTypeJetTypeStateTypeGetObjectTypePassStateTypeIndexTypePassTypeGetCodeTypeCodeTypeSetCodeTypeSetIncomingRequestTypeSetOutgoingRequestTypeSagaCallAcceptNotificationTypeGetFilamentTypeGetRequestTypeRequestTypeFilamentSegmentTypeSetResultTypeActivateTypeRequestInfoTypeGetRequestInfoTypeGotHotConfirmationTypeDeactivateTypeUpdateTypeHotObjectsTypeResultInfoTypeGetPendingsTypeHasPendingsTypePendingsInfoTypeReplicationTypeGetJetTypeAbandonedRequestsNotificationTypeGetLightInitialStateTypeLightInitialStateTypeGetIndexTypeSearchIndexTypeUpdateJetTypeReturnResultsTypeCallMethodTypeExecutorResultsTypePendingFinishedTypeAdditionalCallFromPreviousExecutorTypeStillExecutingTypeErrorResultExitsts_latestType"

var _Type_index = [...]uint16{0, 11, 19, 28, 34, 41, 48, 57, 70, 83, 92, 100, 111, 119, 130, 152, 174, 204, 219, 233, 244, 263, 276, 288, 303, 321, 343, 357, 367, 381, 395, 410, 425, 441, 456, 466, 499, 523, 544, 556, 571, 584, 601, 615, 634, 653, 691, 709, 731, 742}

func (i Type) String() string {
	if i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
