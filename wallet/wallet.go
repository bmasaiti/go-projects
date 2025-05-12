package wallet

import (
	"errors"
	"fmt"
)

// lets you create new types from existing ones.

// The syntax is type MyName OriginalType
type Bitcoin int

//	type Wallet struct {
//		balance int
//	}
type Wallet struct {
	balance Bitcoin
}

type Stringer interface {
	String() string
}

// func (w Wallet) Deposit(amount int) {
// 	w.balance += amount
// }

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance //automatically direferenced , They're  called struct pointers and they are direferenced automatically hence referencing the value directly
}

// func (w *Wallet) Withdraw(amt Bitcoin) Bitcoin {
// 	w.balance -= Bitcoin(amt)
// 	return nil

// }
func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return errors.New("oh no")
	}

	w.balance -= amount
	return nil
}
