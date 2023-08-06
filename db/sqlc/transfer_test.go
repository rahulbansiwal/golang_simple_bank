package sqlc

import (
	"context"
	"database/sql"
	"simple_bank/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestGetTransferFromId(t *testing.T) {
	transfer := CreateRandomTransfer(t)
	transfer1, err := testQueries.GetTransferFromId(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.Equal(t, transfer.ID, transfer1.ID)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := CreateRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.Empty(t, err)
	transfer1, err := testQueries.GetTransferFromId(context.Background(), transfer.ID)
	require.Error(t, err)
	require.Empty(t, transfer1)
}

func CreateRandomTransfer(t *testing.T) Transfer {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	args := CreateTransferParams{
		FromAccountID: sql.NullInt64{Int64: account1.ID, Valid: true},
		ToAccountID:   sql.NullInt64{Int64: account2.ID, Valid: true},
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.Empty(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, transfer.FromAccountID.Int64, account1.ID)
	require.Equal(t, transfer.ToAccountID.Int64, account2.ID)
	return transfer
}
