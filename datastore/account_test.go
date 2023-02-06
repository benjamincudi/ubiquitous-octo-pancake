package datastore

import (
	"context"
	"errors"
	"testing"
)

func Test_MemoryDS_Withdraw(t *testing.T) {
	ctx := context.Background()

	t.Run("happy path", func(t *testing.T) {
		ds := NewMemoryDS()
		if ds.Withdraw(ctx, WithdrawalInfo{12345, 500}) != nil {
			t.Fail()
		}
		if ds.Withdraw(ctx, WithdrawalInfo{12345, 499}) != nil {
			t.Fail()
		}
		if ds.Withdraw(ctx, WithdrawalInfo{12345, 1}) != nil {
			t.Fail()
		}
	})

	t.Run("sad path - too many withdrawals", func(t *testing.T) {
		ds := NewMemoryDS()
		withdrawalAttempts := []int{1, 2, 3}
		for i := range withdrawalAttempts {
			if ds.Withdraw(ctx, WithdrawalInfo{54321, 100}) != nil {
				t.Errorf("expected to be able to make %d withdrawals, failed after %d", len(withdrawalAttempts), i)
				return
			}
		}
		if err := ds.Withdraw(ctx, WithdrawalInfo{54321, 1}); err == nil || !errors.Is(err, ErrDailyUsage{}) {
			t.Errorf("expected ErrDailyUsage, got: (%v)", err)
		}
	})

	t.Run("sad path - withdraw too much", func(t *testing.T) {
		ds := NewMemoryDS()
		if err := ds.Withdraw(ctx, WithdrawalInfo{999, 500}); err != nil {
			t.Error(err)
		}
		if err := ds.Withdraw(ctx, WithdrawalInfo{999, 400}); err != nil {
			t.Error(err)
		}

		if err := ds.Withdraw(ctx, WithdrawalInfo{999, 200}); err == nil || !errors.Is(err, ErrDailyAmount{remainingAmount: 100}) {
			t.Errorf("expected ErrDailyAmount{100}, got: (%v)", err)
		}
	})

	t.Run("sad path - transaction too large", func(t *testing.T) {
		ds := NewMemoryDS()
		if err := ds.Withdraw(ctx, WithdrawalInfo{1, 501}); err == nil || !errors.Is(err, ErrTooLarge{}) {
			t.Errorf("expected ErrTooLarge, got: (%v)", err)
		}
	})
}
