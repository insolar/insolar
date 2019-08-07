package account

import (
	"fmt"
	"math/big"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/builtin/foundation/safemath"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/wallet"
)

type Account struct {
	foundation.BaseContract
	Balance string
	Name    string
}

func New(name string, balance string) (*Account, error) {
	return &Account{Name: name, Balance: balance}, nil
}

// Transfer transfers funds to giver reference.
func (a *Account) Transfer(amountStr string, toWallet *insolar.Reference) (err error) {
	amount, _ := new(big.Int).SetString(amountStr, 10)
	balance, ok := new(big.Int).SetString(a.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse wallet balance")
	}

	newBalance, err := safemath.Sub(balance, amount)
	if err != nil {
		return fmt.Errorf("not enough balance for transfer: %s", err.Error())
	}
	a.Balance = newBalance.String()
	destWallet := wallet.GetObject(*toWallet)
	return destWallet.Accept(amountStr, a.Name)
}

// Accept accepts transfer to balance.
func (a *Account) Accept(amountStr string) (err error) {

	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("can't parse input amount")
	}

	balance := new(big.Int)
	balance, ok = balance.SetString(a.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse wallet balance")
	}

	b, err := safemath.Add(balance, amount)
	if err != nil {
		return fmt.Errorf("failed to add amount to balance: %s", err.Error())
	}
	a.Balance = b.String()

	return nil
}

// RollBack rolls back transfer to balance.
func (a *Account) RollBack(amountStr string) (err error) {

	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("can't parse input amount")
	}

	balance := new(big.Int)
	balance, ok = balance.SetString(a.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse wallet balance")
	}

	b, err := safemath.Sub(balance, amount)
	if err != nil {
		return fmt.Errorf("failed to sub amount from balance: %s", err.Error())
	}
	a.Balance = b.String()

	return nil
}

// GetBalance gets total balance.
func (a *Account) GetBalance() (string, error) {
	return a.Balance, nil
}
