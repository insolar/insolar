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

package migrationdaemon

import (
	"fmt"
	"math/big"

	"github.com/insolar/insolar/application/appfoundation"
	"github.com/insolar/insolar/application/builtin/proxy/deposit"
	"github.com/insolar/insolar/application/builtin/proxy/member"
	"github.com/insolar/insolar/application/builtin/proxy/migrationadmin"
	"github.com/insolar/insolar/application/builtin/proxy/wallet"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

const CONVERSION = "10"

// MigrationDaemon make migration procedure.
type MigrationDaemon struct {
	foundation.BaseContract
	IsActive              bool
	MigrationDaemonMember insolar.Reference
}

// DepositMigrationResult struct for  return reference from contract.
type DepositMigrationResult struct {
	Reference string `json:"memberReference"`
}

// DepositMigrationCall internal function migration admin from api.
func (md *MigrationDaemon) DepositMigrationCall(
	params map[string]interface{}, caller insolar.Reference, request insolar.Reference,
) (*DepositMigrationResult, error) {

	amount, err := getAmountFromParam(params)
	if err != nil {
		return nil, err
	}

	txId, ok := params["ethTxHash"].(string)
	if !ok {
		return nil, fmt.Errorf("incorrect input: failed to get 'ethTxHash' param")
	}

	migrationAddress, ok := params["migrationAddress"].(string)
	if !ok {
		return nil, fmt.Errorf("incorrect input: failed to get 'migrationAddress' param")
	}
	base, _ := new(big.Int).SetString(CONVERSION, 10)
	amountXns := new(big.Int).Mul(amount, base)

	return md.depositMigration(txId, migrationAddress, amountXns, caller, request)
}

// Set status Migration daemon.
func (md *MigrationDaemon) SetActivationStatus(status bool) error {
	md.IsActive = status
	return nil
}

// Return status migration daemon.
// ins:immutable
func (md *MigrationDaemon) GetActivationStatus() (bool, error) {
	return md.IsActive, nil
}

// Return reference on migration daemon.
// ins:immutable
func (md *MigrationDaemon) GetMigrationDaemonMember() (insolar.Reference, error) {
	return md.MigrationDaemonMember, nil
}

func (md *MigrationDaemon) depositMigration(
	txHash string, migrationAddress string, amount *big.Int, caller insolar.Reference, request insolar.Reference,
) (*DepositMigrationResult, error) {

	if !caller.Equal(md.MigrationDaemonMember) {
		return nil, fmt.Errorf(" the migration daemon member is not related migration daemon contract, %s ", caller)
	}

	result, err := md.GetActivationStatus()
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, fmt.Errorf("this migration daemon is not active daemons: %s", caller)
	}

	migrationAdminContract := migrationadmin.GetObject(appfoundation.GetMigrationAdmin())
	// Get member by migration address
	tokenHolderRef, err := migrationAdminContract.GetMemberByMigrationAddress(migrationAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get member by migration address")
	}
	tokenHolder := member.GetObject(*tokenHolderRef)
	tokenHolderWallet, err := tokenHolder.GetWallet()
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %s", err.Error())
	}

	w := wallet.GetObject(*tokenHolderWallet)
	vestingParams, _ := migrationAdminContract.GetDepositParameters()
	depositRef, err := w.FindOrCreateDeposit(txHash, vestingParams.Lockup, vestingParams.Vesting, vestingParams.VestingStep)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create deposit: %s", err.Error())
	}

	return addConfirmToDeposit(*tokenHolderRef, *depositRef, txHash, amount.String(), caller, request)
}

func getAmountFromParam(params map[string]interface{}) (*big.Int, error) {
	amountStr, ok := params["amount"].(string)
	if !ok {
		return nil, fmt.Errorf("incorrect input: failed to get 'amount' param")
	}

	amount := new(big.Int)
	amount, ok = amount.SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse amount")
	}
	if amount.Sign() <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	return amount, nil
}

func addConfirmToDeposit(
	tokenHolderRef insolar.Reference,
	txDepositRef insolar.Reference,
	txHash string,
	amount string,
	caller insolar.Reference,
	request insolar.Reference,
) (*DepositMigrationResult, error) {
	txDeposit := deposit.GetObject(txDepositRef)

	err := txDeposit.Confirm(txHash, amount, caller, request, tokenHolderRef)
	if err != nil {
		return nil, fmt.Errorf("confirmed failed: %s", err.Error())
	}

	return &DepositMigrationResult{Reference: tokenHolderRef.String()}, nil
}
