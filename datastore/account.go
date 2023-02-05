package datastore

import (
	"context"
	"fmt"
)

type MemoryDS struct {
	dailyUsage map[int]dailyAccountUsage
}

func NewMemoryDS() MemoryDS {
	return MemoryDS{map[int]dailyAccountUsage{}}
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
	return "You may only withdraw up to $1000 per day from our ATM and are at the limit for today."
}

type ErrTooLarge struct{}

func (e ErrTooLarge) Error() string {
	return "You may only withdraw up to $500 in a single transaction."
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
	ds.dailyUsage[info.PersonalAccountNumber] = dailyAccountUsage{
		amount:      usage.amount + info.Amount,
		withdrawals: usage.withdrawals + 1,
	}
	return nil
}
