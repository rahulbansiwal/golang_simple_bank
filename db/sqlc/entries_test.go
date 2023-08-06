package sqlc

import (
	"context"
	"database/sql"
	"simple_bank/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestGetEntry(t *testing.T){
	entry := CreateRandomEntry(t)
	entry2,err:= testQueries.GetEntry(context.Background(),entry.ID)
	require.NoError(t,err)
	require.NotEmpty(t,entry2)
	require.Equal(t,entry2.ID,entry.ID)
	require.Equal(t,entry2.Amount,entry.Amount)
	require.Equal(t,entry2.AccountID,entry.AccountID)
}

func TestGetEntriesFromAccountId(t *testing.T){
	for i:=0;i<10;i++{
		 CreateRandomEntry(t)
	}
	params := GetEntriesFromAccountIdParams{
		Limit: 5,
		Offset: 5,
	}
	entry,err:= testQueries.GetEntriesFromAccountId(context.Background(),params)
	require.NoError(t,err)
	require.Len(t,entry,5)
	
}

func TestDeleteEntry(t *testing.T){
	account := CreateRandomEntry(t)
	err:= testQueries.DeleteEntry(context.Background(),account.ID)
	require.NoError(t,err)
	require.Empty(t,err)
	account2,err := testQueries.GetEntry(context.Background(),account.ID)
	require.Error(t,err)
	require.Empty(t,account2)
}


func CreateRandomEntry(t *testing.T) Entry {
	account := CreateRandomAccount(t)
	args := CreateEntryParams{
		AccountID: sql.NullInt64{Int64: account.ID, Valid: true},
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.Empty(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, entry.AccountID, args.AccountID)
	require.Equal(t, args.Amount, entry.Amount)
	require.NotEmpty(t, entry.ID)
	return entry
}


