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

package deposit

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/insolar/insolar/application/appfoundation"
	"github.com/insolar/insolar/application/builtin/proxy/deposit"
	"github.com/insolar/insolar/application/builtin/proxy/member"
	"github.com/insolar/insolar/application/builtin/proxy/migrationdaemon"
	"github.com/insolar/insolar/application/builtin/proxy/wallet"
	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/builtin/foundation/safemath"
)

const numConfirmation = 2

// Deposit is like wallet. It holds migrated money.
type Deposit struct {
	foundation.BaseContract
	Balance                 string                    `json:"balance"`
	PulseDepositUnHold      insolar.PulseNumber       `json:"holdReleaseDate"`
	MigrationDaemonConfirms foundation.StableMap      `json:"confirmerReferences"`
	Amount                  string                    `json:"amount"`
	TxHash                  string                    `json:"ethTxHash"`
	VestingType             appfoundation.VestingType `json:"vestingType"`
	Lockup                  int64                     `json:"lockupInPulses"`
	Vesting                 int64                     `json:"vestingInPulses"`
	VestingStep             int64                     `json:"vestingStepInPulses"`
	IsConfirmed             bool                      `json:"isConfirmed"`
}

// New creates new deposit.
func New(txHash string, lockup int64, vesting int64, vestingStep int64) (*Deposit, error) {

	migrationDaemonConfirms := make(foundation.StableMap)

	return &Deposit{
		Balance:                 "0",
		MigrationDaemonConfirms: migrationDaemonConfirms,
		Amount:                  "0",
		TxHash:                  txHash,
		Lockup:                  lockup,
		Vesting:                 vesting,
		VestingStep:             vestingStep,
		VestingType:             appfoundation.DefaultVesting,
	}, nil
}

// Form of Deposit that is applied in API
type DepositOut struct {
	Balance                 string                    `json:"balance"`
	HoldStartDate           int64                     `json:"holdStartDate"`
	PulseDepositUnHold      int64                     `json:"holdReleaseDate"`
	MigrationDaemonConfirms []DaemonConfirm           `json:"confirmerReferences"`
	Amount                  string                    `json:"amount"`
	TxHash                  string                    `json:"ethTxHash"`
	VestingType             appfoundation.VestingType `json:"vestingType"`
	Lockup                  int64                     `json:"lockup"`
	Vesting                 int64                     `json:"vesting"`
	VestingStep             int64                     `json:"vestingStep"`
}

type DaemonConfirm struct {
	Reference string `json:"reference"`
	Amount    string `json:"amount"`
}

// GetTxHash gets transaction hash.
// ins:immutable
func (d *Deposit) GetTxHash() (string, error) {
	return d.TxHash, nil
}

// GetAmount gets amount.
// ins:immutable
func (d *Deposit) GetAmount() (string, error) {
	return d.Amount, nil
}

// Return pulse of unhold deposit.
// ins:immutable
func (d *Deposit) GetPulseUnHold() (insolar.PulseNumber, error) {
	return d.PulseDepositUnHold, nil
}

// Itself gets deposit information.
// ins:immutable
func (d *Deposit) Itself() (interface{}, error) {
	var daemonConfirms = make([]DaemonConfirm, 0, len(d.MigrationDaemonConfirms))
	var pulseDepositUnHold int64
	for k, v := range d.MigrationDaemonConfirms {
		daemonConfirms = append(daemonConfirms, DaemonConfirm{Reference: k, Amount: v})
	}
	t, err := d.PulseDepositUnHold.AsApproximateTime()
	if err == nil {
		pulseDepositUnHold = t.Unix()
	}
	return &DepositOut{
		Balance:                 d.Balance,
		HoldStartDate:           pulseDepositUnHold - d.Lockup,
		PulseDepositUnHold:      pulseDepositUnHold,
		MigrationDaemonConfirms: daemonConfirms,
		Amount:                  d.Amount,
		TxHash:                  d.TxHash,
		VestingType:             d.VestingType,
		Lockup:                  d.Lockup,
		Vesting:                 d.Vesting,
		VestingStep:             d.VestingStep,
	}, nil
}

// Confirm adds confirm for deposit by migration daemon.
func (d *Deposit) Confirm(
	txHash string, amountStr string, fromMember insolar.Reference, request insolar.Reference, toMember insolar.Reference,
) error {

	migrationDaemonRef := fromMember.String()
	if txHash != d.TxHash {
		return fmt.Errorf("transaction hash is incorrect")
	}
	if confirmedAmount, ok := d.MigrationDaemonConfirms[migrationDaemonRef]; ok {
		if amountStr != confirmedAmount {
			return fmt.Errorf(
				"confirm from this migration daemon %s already exists with different amount: was %s, now %s",
				migrationDaemonRef,
				confirmedAmount,
				amountStr,
			)
		}
		return nil
	}

	if d.IsConfirmed {
		d.MigrationDaemonConfirms[migrationDaemonRef] = amountStr
		if amountStr != d.Amount {
			return fmt.Errorf(
				"migration is done for this deposit %s, but with different amount: confirmed is %s, from this daemon %s",
				txHash,
				d.Amount,
				amountStr,
			)
		}
		return nil
	}

	if len(d.MigrationDaemonConfirms) > 0 {
		canConfirm, errFromConfirm := d.checkConfirm(migrationDaemonRef, amountStr)
		if canConfirm {
			currentPulse, err := foundation.GetPulseNumber()
			if err != nil {
				return fmt.Errorf("failed to get current pulse: %s", err.Error())
			}
			d.Amount = amountStr
			d.PulseDepositUnHold = currentPulse + insolar.PulseNumber(d.Lockup)

			ma := member.GetObject(appfoundation.GetMigrationAdminMember())
			walletRef, err := ma.GetWallet()
			if err != nil {
				return fmt.Errorf("failed to get wallet: %s", err.Error())
			}
			ok, maDeposit, _ := wallet.GetObject(*walletRef).FindDeposit(genesisrefs.FundsDepositName)
			if !ok {
				return fmt.Errorf("failed to find source deposit - %s", walletRef.String())
			}

			err = deposit.GetObject(*maDeposit).TransferToDeposit(
				amountStr, d.GetReference(), appfoundation.GetMigrationAdminMember(), request, toMember,
			)
			if err != nil {
				return fmt.Errorf("failed to transfer from migration deposit to deposit: %s", err.Error())
			}
			d.IsConfirmed = true
		}
		if errFromConfirm != nil {
			return errFromConfirm
		}
		return nil
	}
	d.MigrationDaemonConfirms[migrationDaemonRef] = amountStr
	return nil
}

// TransferToDeposit transfers funds to deposit.
func (d *Deposit) TransferToDeposit(
	amountStr string,
	toDeposit insolar.Reference,
	fromMember insolar.Reference,
	request insolar.Reference,
	toMember insolar.Reference,
) error {
	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("can't parse input amount")
	}
	balance, ok := new(big.Int).SetString(d.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse deposit balance")
	}
	if balance.Sign() <= 0 {
		return fmt.Errorf("not enough balance for transfer")
	}
	newBalance, err := safemath.Sub(balance, amount)
	if err != nil {
		return fmt.Errorf("not enough balance for transfer: %s", err.Error())
	}
	d.Balance = newBalance.String()
	destination := deposit.GetObject(toDeposit)
	acceptDepositErr := destination.Accept(appfoundation.SagaAcceptInfo{
		Amount:     amountStr,
		FromMember: fromMember,
		Request:    request,
	})
	if acceptDepositErr == nil {
		return nil
	}
	d.Balance = balance.String()
	return fmt.Errorf("failed to transfer amount: %s", acceptDepositErr.Error())

}

// Check amount field in confirmation from migration daemons.
func (d *Deposit) checkAmount(activeDaemons map[string]string, migrationDaemonRef string, amountStr string) (bool, error) {
	confirmed := false
	if activeDaemons == nil || len(activeDaemons) == 0 {
		return false, fmt.Errorf("list with migration daemons member is empty")
	}
	amountConfirms := make(map[string]int) // amount: num of confirms
	var errDaemon []string
	for migrationRef, amount := range activeDaemons {
		if amount != amountStr {
			errDaemon = append(errDaemon, fmt.Sprintf("%s send amount %s", migrationRef, amount))
		}
		amountConfirms[amount] = amountConfirms[amount] + 1
	}
	amountConfirms[amountStr] = amountConfirms[amountStr] + 1
	if amountConfirms[amountStr] >= numConfirmation {
		confirmed = true
	}
	var err error
	if len(errDaemon) > 0 {
		if !confirmed {
			errDaemon = append(errDaemon, fmt.Sprintf("%s send amount %s", migrationDaemonRef, amountStr))
		}
		err = fmt.Errorf("several migration daemons send different amount: %s", strings.Join(errDaemon, ": "))
	}
	return confirmed, err
}

func (d *Deposit) checkConfirm(migrationDaemonRef string, amountStr string) (bool, error) {
	activateDaemons := make(map[string]string)

	for ref, a := range d.MigrationDaemonConfirms {
		migrationDaemonMemberRef, err := insolar.NewObjectReferenceFromString(ref)
		if err != nil {
			return false, fmt.Errorf(" failed to parse params.Reference")
		}

		migrationDaemonContractRef, err := appfoundation.GetMigrationDaemon(*migrationDaemonMemberRef)
		if err != nil || migrationDaemonContractRef.IsEmpty() {
			return false, fmt.Errorf(" get migration daemon contract from foundation failed, %s ", err)
		}

		migrationDaemonContract := migrationdaemon.GetObject(migrationDaemonContractRef)
		result, err := migrationDaemonContract.GetActivationStatus()

		if err != nil {
			return false, err
		}
		if result {
			activateDaemons[ref] = a
		}
	}
	d.MigrationDaemonConfirms[migrationDaemonRef] = amountStr
	if len(activateDaemons) > 0 {
		canConfirm, err := d.checkAmount(activateDaemons, migrationDaemonRef, amountStr)
		if err != nil {
			return canConfirm, fmt.Errorf("failed to check amount in confirmation from migration daemon: %s", err.Error())
		}
		return canConfirm, nil
	}
	return false, nil
}

func (d *Deposit) canTransfer(transferAmount *big.Int) error {
	c := 0
	for _, r := range d.MigrationDaemonConfirms {
		if r != "" {
			c++
		}
	}
	if d.VestingType == appfoundation.DefaultVesting && c < numConfirmation {
		return fmt.Errorf("number of confirms is less then 2")
	}

	currentPulse, err := foundation.GetPulseNumber()
	if err != nil {
		return fmt.Errorf("failed to get pulse number: %s", err.Error())
	}
	if d.PulseDepositUnHold > currentPulse {
		return fmt.Errorf("hold period didn't end")
	}

	spentPeriodInPulses := big.NewInt(int64(currentPulse-d.PulseDepositUnHold) / d.VestingStep)
	amount, ok := new(big.Int).SetString(d.Amount, 10)
	if !ok {
		return fmt.Errorf("can't parse derposit amount")
	}
	balance, ok := new(big.Int).SetString(d.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse derposit balance")
	}

	// How much can we transfer for this time
	availableForNow := new(big.Int).Div(
		new(big.Int).Mul(amount, spentPeriodInPulses),
		big.NewInt(d.Vesting/d.VestingStep),
	)

	if new(big.Int).Sub(amount, availableForNow).Cmp(
		new(big.Int).Sub(balance, transferAmount),
	) == 1 {
		return fmt.Errorf("not enough unholded balance for transfer")
	}

	return nil
}

// Transfer transfers money from deposit to wallet. It can be called only after deposit hold period.
func (d *Deposit) Transfer(
	amountStr string, memberRef insolar.Reference, request insolar.Reference,
) (interface{}, error) {

	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse input amount")
	}

	balance, ok := new(big.Int).SetString(d.Balance, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse deposit balance")
	}
	if balance.Sign() <= 0 {
		return nil, fmt.Errorf("not enough balance for transfer")
	}
	newBalance, err := safemath.Sub(balance, amount)
	if err != nil {
		return nil, fmt.Errorf("not enough balance for transfer: %s", err.Error())
	}
	err = d.canTransfer(amount)
	if err != nil {
		return nil, fmt.Errorf("can't start transfer: %s", err.Error())
	}
	d.Balance = newBalance.String()

	m := member.GetObject(memberRef)
	acceptMemberErr := m.Accept(appfoundation.SagaAcceptInfo{
		Amount:     amountStr,
		FromMember: memberRef,
		Request:    request,
	})
	if acceptMemberErr == nil {
		return nil, nil
	}
	d.Balance = balance.String()
	return nil, fmt.Errorf("failed to transfer amount: %s", acceptMemberErr.Error())
}

// Accept accepts transfer to balance.
// ins:saga(INS_FLAG_NO_ROLLBACK_METHOD)
func (d *Deposit) Accept(arg appfoundation.SagaAcceptInfo) error {

	amount := new(big.Int)
	amount, ok := amount.SetString(arg.Amount, 10)
	if !ok {
		return fmt.Errorf("can't parse input amount")
	}

	balance := new(big.Int)
	balance, ok = balance.SetString(d.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse deposit balance")
	}

	b, err := safemath.Add(balance, amount)
	if err != nil {
		return fmt.Errorf("failed to add amount to balance: %s", err.Error())
	}
	d.Balance = b.String()

	return nil
}
