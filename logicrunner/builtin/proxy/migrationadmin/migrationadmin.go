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

package migrationadmin

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/common"
)

type CheckDaemonResponse struct {
	Status string `json:"status"`
}
type VestingParams struct {
	Lokup       int64 `json:"lokupInPulses"`
	Vesting     int64 `json:"vestingInPulses"`
	VestingStep int64 `json:"vestingStepInPulses"`
}

// PrototypeReference to prototype of this contract
// error checking hides in generator
var PrototypeReference, _ = insolar.NewReferenceFromBase58("111A8DhUhw5pzyvzVg1qXomNEHXs7kDtJRQGSD1PUpc.11111111111111111111111111111111")

// MigrationAdmin holds proxy type
type MigrationAdmin struct {
	Reference insolar.Reference
	Prototype insolar.Reference
	Code      insolar.Reference
}

// ContractConstructorHolder holds logic with object construction
type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

// AsChild saves object as child
func (r *ContractConstructorHolder) AsChild(objRef insolar.Reference) (*MigrationAdmin, error) {
	ref, ret, err := common.CurrentProxyCtx.SaveAsChild(objRef, *PrototypeReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}

	var constructorError *foundation.Error
	resultContainer := foundation.Result{
		Returns: []interface{}{&constructorError},
	}
	err = common.CurrentProxyCtx.Deserialize(ret, &resultContainer)
	if err != nil {
		return nil, err
	}

	if resultContainer.Error != nil {
		return nil, resultContainer.Error
	}

	if constructorError != nil {
		return nil, constructorError
	}

	return &MigrationAdmin{Reference: *ref}, nil
}

// GetObject returns proxy object
func GetObject(ref insolar.Reference) (r *MigrationAdmin) {
	return &MigrationAdmin{Reference: ref}
}

// GetPrototype returns reference to the prototype
func GetPrototype() insolar.Reference {
	return *PrototypeReference
}

// GetReference returns reference of the object
func (r *MigrationAdmin) GetReference() insolar.Reference {
	return r.Reference
}

// GetPrototype returns reference to the code
func (r *MigrationAdmin) GetPrototype() (insolar.Reference, error) {
	if r.Prototype.IsEmpty() {
		ret := [2]interface{}{}
		var ret0 insolar.Reference
		ret[0] = &ret0
		var ret1 *foundation.Error
		ret[1] = &ret1

		res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetPrototype", make([]byte, 0), *PrototypeReference)
		if err != nil {
			return ret0, err
		}

		err = common.CurrentProxyCtx.Deserialize(res, &ret)
		if err != nil {
			return ret0, err
		}

		if ret1 != nil {
			return ret0, ret1
		}

		r.Prototype = ret0
	}

	return r.Prototype, nil

}

// GetCode returns reference to the code
func (r *MigrationAdmin) GetCode() (insolar.Reference, error) {
	if r.Code.IsEmpty() {
		ret := [2]interface{}{}
		var ret0 insolar.Reference
		ret[0] = &ret0
		var ret1 *foundation.Error
		ret[1] = &ret1

		res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetCode", make([]byte, 0), *PrototypeReference)
		if err != nil {
			return ret0, err
		}

		err = common.CurrentProxyCtx.Deserialize(res, &ret)
		if err != nil {
			return ret0, err
		}

		if ret1 != nil {
			return ret0, ret1
		}

		r.Code = ret0
	}

	return r.Code, nil
}

// MigrationAdminCall is proxy generated method
func (r *MigrationAdmin) MigrationAdminCall(params map[string]interface{}, nameMethod string, caller insolar.Reference) (interface{}, error) {
	var args [3]interface{}
	args[0] = params
	args[1] = nameMethod
	args[2] = caller

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 interface{}
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "MigrationAdminCall", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// MigrationAdminCallNoWait is proxy generated method
func (r *MigrationAdmin) MigrationAdminCallNoWait(params map[string]interface{}, nameMethod string, caller insolar.Reference) error {
	var args [3]interface{}
	args[0] = params
	args[1] = nameMethod
	args[2] = caller

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "MigrationAdminCall", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// MigrationAdminCallAsImmutable is proxy generated method
func (r *MigrationAdmin) MigrationAdminCallAsImmutable(params map[string]interface{}, nameMethod string, caller insolar.Reference) (interface{}, error) {
	var args [3]interface{}
	args[0] = params
	args[1] = nameMethod
	args[2] = caller

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 interface{}
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "MigrationAdminCall", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetAllMigrationDaemon is proxy generated method
func (r *MigrationAdmin) GetAllMigrationDaemonAsMutable() (foundation.StableMap, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 foundation.StableMap
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetAllMigrationDaemon", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetAllMigrationDaemonNoWait is proxy generated method
func (r *MigrationAdmin) GetAllMigrationDaemonNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetAllMigrationDaemon", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetAllMigrationDaemonAsImmutable is proxy generated method
func (r *MigrationAdmin) GetAllMigrationDaemon() (foundation.StableMap, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 foundation.StableMap
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetAllMigrationDaemon", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// ActivateDaemon is proxy generated method
func (r *MigrationAdmin) ActivateDaemon(daemonMember string, caller insolar.Reference) error {
	var args [2]interface{}
	args[0] = daemonMember
	args[1] = caller

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "ActivateDaemon", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}

// ActivateDaemonNoWait is proxy generated method
func (r *MigrationAdmin) ActivateDaemonNoWait(daemonMember string, caller insolar.Reference) error {
	var args [2]interface{}
	args[0] = daemonMember
	args[1] = caller

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "ActivateDaemon", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// ActivateDaemonAsImmutable is proxy generated method
func (r *MigrationAdmin) ActivateDaemonAsImmutable(daemonMember string, caller insolar.Reference) error {
	var args [2]interface{}
	args[0] = daemonMember
	args[1] = caller

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "ActivateDaemon", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}

// DeactivateDaemon is proxy generated method
func (r *MigrationAdmin) DeactivateDaemon(daemonMember string, caller insolar.Reference) error {
	var args [2]interface{}
	args[0] = daemonMember
	args[1] = caller

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "DeactivateDaemon", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}

// DeactivateDaemonNoWait is proxy generated method
func (r *MigrationAdmin) DeactivateDaemonNoWait(daemonMember string, caller insolar.Reference) error {
	var args [2]interface{}
	args[0] = daemonMember
	args[1] = caller

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "DeactivateDaemon", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// DeactivateDaemonAsImmutable is proxy generated method
func (r *MigrationAdmin) DeactivateDaemonAsImmutable(daemonMember string, caller insolar.Reference) error {
	var args [2]interface{}
	args[0] = daemonMember
	args[1] = caller

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "DeactivateDaemon", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}

// CheckDaemon is proxy generated method
func (r *MigrationAdmin) CheckDaemonAsMutable(daemonMember string) (bool, error) {
	var args [1]interface{}
	args[0] = daemonMember

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 bool
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "CheckDaemon", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// CheckDaemonNoWait is proxy generated method
func (r *MigrationAdmin) CheckDaemonNoWait(daemonMember string) error {
	var args [1]interface{}
	args[0] = daemonMember

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "CheckDaemon", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// CheckDaemonAsImmutable is proxy generated method
func (r *MigrationAdmin) CheckDaemon(daemonMember string) (bool, error) {
	var args [1]interface{}
	args[0] = daemonMember

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 bool
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "CheckDaemon", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetActiveDaemons is proxy generated method
func (r *MigrationAdmin) GetActiveDaemonsAsMutable() ([]string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 []string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetActiveDaemons", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetActiveDaemonsNoWait is proxy generated method
func (r *MigrationAdmin) GetActiveDaemonsNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetActiveDaemons", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetActiveDaemonsAsImmutable is proxy generated method
func (r *MigrationAdmin) GetActiveDaemons() ([]string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 []string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetActiveDaemons", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetDepositParameters is proxy generated method
func (r *MigrationAdmin) GetDepositParameters() (*VestingParams, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 *VestingParams
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetDepositParameters", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetDepositParametersNoWait is proxy generated method
func (r *MigrationAdmin) GetDepositParametersNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetDepositParameters", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetDepositParametersAsImmutable is proxy generated method
func (r *MigrationAdmin) GetDepositParametersAsImmutable() (*VestingParams, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 *VestingParams
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetDepositParameters", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetMemberByMigrationAddress is proxy generated method
func (r *MigrationAdmin) GetMemberByMigrationAddressAsMutable(migrationAddress string) (*insolar.Reference, error) {
	var args [1]interface{}
	args[0] = migrationAddress

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 *insolar.Reference
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetMemberByMigrationAddress", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetMemberByMigrationAddressNoWait is proxy generated method
func (r *MigrationAdmin) GetMemberByMigrationAddressNoWait(migrationAddress string) error {
	var args [1]interface{}
	args[0] = migrationAddress

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetMemberByMigrationAddress", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetMemberByMigrationAddressAsImmutable is proxy generated method
func (r *MigrationAdmin) GetMemberByMigrationAddress(migrationAddress string) (*insolar.Reference, error) {
	var args [1]interface{}
	args[0] = migrationAddress

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 *insolar.Reference
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetMemberByMigrationAddress", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetFreeMigrationAddress is proxy generated method
func (r *MigrationAdmin) GetFreeMigrationAddressAsMutable(publicKey string) (string, error) {
	var args [1]interface{}
	args[0] = publicKey

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetFreeMigrationAddress", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetFreeMigrationAddressNoWait is proxy generated method
func (r *MigrationAdmin) GetFreeMigrationAddressNoWait(publicKey string) error {
	var args [1]interface{}
	args[0] = publicKey

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetFreeMigrationAddress", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetFreeMigrationAddressAsImmutable is proxy generated method
func (r *MigrationAdmin) GetFreeMigrationAddress(publicKey string) (string, error) {
	var args [1]interface{}
	args[0] = publicKey

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetFreeMigrationAddress", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// AddNewMigrationAddressToMaps is proxy generated method
func (r *MigrationAdmin) AddNewMigrationAddressToMapsAsMutable(migrationAddress string, memberRef insolar.Reference) error {
	var args [2]interface{}
	args[0] = migrationAddress
	args[1] = memberRef

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "AddNewMigrationAddressToMaps", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}

// AddNewMigrationAddressToMapsNoWait is proxy generated method
func (r *MigrationAdmin) AddNewMigrationAddressToMapsNoWait(migrationAddress string, memberRef insolar.Reference) error {
	var args [2]interface{}
	args[0] = migrationAddress
	args[1] = memberRef

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "AddNewMigrationAddressToMaps", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// AddNewMigrationAddressToMapsAsImmutable is proxy generated method
func (r *MigrationAdmin) AddNewMigrationAddressToMaps(migrationAddress string, memberRef insolar.Reference) error {
	var args [2]interface{}
	args[0] = migrationAddress
	args[1] = memberRef

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "AddNewMigrationAddressToMaps", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}
