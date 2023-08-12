package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		rberr := tx.Rollback()
		if rberr != nil {
			return fmt.Errorf("rollback error is %v, transaction error is %v", rberr, err)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: sql.NullInt64{Int64: arg.FromAccountId, Valid: true},
			ToAccountID:   sql.NullInt64{Int64: arg.ToAccountId, Valid: true},
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: sql.NullInt64{Int64: arg.FromAccountId, Valid: true},
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: sql.NullInt64{Int64: arg.ToAccountId, Valid: true},
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		//TODO : Update account balance for both
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountId)
		if err != nil {
			return err
		}
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountId,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountId)
		if err != nil {
			return err
		}
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountId,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}
