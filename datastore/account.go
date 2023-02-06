package datastore

import (
	"context"
	"fmt"
)

type MemoryDS struct {
	dailyUsage map[int]dailyAccountUsage
}

func NewMemoryDS() *MemoryDS {
	return &MemoryDS{map[int]dailyAccountUsage{}}
}

func (ds MemoryDS) getRemainingSupply() int {
	used := 0
	for _, usage := range ds.dailyUsage {
		used += usage.amount
	}
	return 5000 - used
}

type WithdrawalInfo struct {
	PersonalAccountNumber int `form:"pan"`
	Amount                int `form:"amount"`
}

type dailyAccountUsage struct {
	withdrawals, amount int
}

type ErrDailyUsage struct{}

func (e ErrDailyUsage) Error() string {
	return "You may only make up to 3 withdrawals from our ATM each day."
}

type ErrDailyAmount struct {
	remainingAmount int
}

func (e ErrDailyAmount) Error() string {
	if e.remainingAmount > 0 {
		return fmt.Sprintf("You may only withdraw up to $%d more today from our ATM.", e.remainingAmount)
	}
	return "You may only withdraw up to $1000 per day from our ATM."
}

type ErrTooLarge struct{}

func (e ErrTooLarge) Error() string {
	return "You may only withdraw up to $500 in a single transaction."
}

type ErrATMSupply struct{}

func (e ErrATMSupply) Error() string {
	return "Sorry, this ATM does not have sufficient cash to complete your request. Please visit another location."
}

func (ds MemoryDS) Withdraw(_ context.Context, info WithdrawalInfo) error {
	usage, exists := ds.dailyUsage[info.PersonalAccountNumber]
	if exists {
		if usage.withdrawals == 3 {
			return ErrDailyUsage{}
		}
		if usage.amount+info.Amount > 1000 {
			return ErrDailyAmount{1000 - usage.amount}
		}
	}
	if info.Amount > 500 {
		return ErrTooLarge{}
	}
	if info.Amount > ds.getRemainingSupply() {
		return ErrATMSupply{}
	}
	ds.dailyUsage[info.PersonalAccountNumber] = dailyAccountUsage{
		amount:      usage.amount + info.Amount,
		withdrawals: usage.withdrawals + 1,
	}
	return nil
}

func (ds *MemoryDS) ResetDay(_ context.Context) error {
	ds.dailyUsage = map[int]dailyAccountUsage{}

	return nil
}
