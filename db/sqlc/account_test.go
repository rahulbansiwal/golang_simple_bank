package sqlc

import (
	"context"
	"database/sql"
	"simple_bank/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := CreateRandomAccount(t)
	getaccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getaccount)
	require.Equal(t, account.ID, getaccount.ID)
	require.Equal(t, account.Owner, getaccount.Owner)
	require.Equal(t, account.Balance, getaccount.Balance)
	require.Equal(t, account.Currency, getaccount.Currency)
	require.NotEmpty(t, getaccount.CreatedAt)
}

func CreateRandomAccount(t *testing.T) Account {
	user := CreateRandomUser(t)
	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestUpdateAccount(t *testing.T) {
	account := CreateRandomAccount(t)
	args := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}
	updateAccount, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount)
	require.Equal(t, updateAccount.Balance, args.Balance)
	require.Equal(t, updateAccount.ID, args.ID)
}

func TestDeleteAccount(t *testing.T) {
	account := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount =  CreateRandomAccount(t)
	}

	args := ListAccountsParams{
		Owner: lastAccount.Owner,
		Limit: 5,
		Offset: 0,
	}

	account,err:= testQueries.ListAccounts(context.Background(),args)
	require.NoError(t,err)
	require.NotEmpty(t,account)

}
