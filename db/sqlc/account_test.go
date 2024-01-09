package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	accoutn1 := createRandomAccount(t)
	accoutn2, err := testQueries.GetAccount(context.Background(), accoutn1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accoutn2)

	require.Equal(t, accoutn1.ID, accoutn2.ID)
	require.Equal(t, accoutn1.Owner, accoutn2.Owner)
	require.Equal(t, accoutn1.Balance, accoutn2.Balance)
	require.Equal(t, accoutn1.Currency, accoutn2.Currency)
	require.WithinDuration(t, accoutn1.CreatedAt, accoutn2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	accoutn1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      accoutn1.ID,
		Balance: util.RandomMoney(),
	}

	accoutn2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accoutn2)

	require.Equal(t, accoutn1.ID, accoutn2.ID)
	require.Equal(t, accoutn1.Owner, accoutn2.Owner)
	require.Equal(t, arg.Balance, accoutn2.Balance)
	require.Equal(t, accoutn1.Currency, accoutn2.Currency)
	require.WithinDuration(t, accoutn1.CreatedAt, accoutn2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	accoutn1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), accoutn1.ID)
	require.NoError(t, err)

	accoutn2, err := testQueries.GetAccount(context.Background(), accoutn1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accoutn2)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, acaccount := range accounts {
		require.NotEmpty(t, acaccount)
	}
}

func TestListAccountWithInvalidLimit(t *testing.T) {
	arg := ListAccountsParams{
		Limit:  -5,
		Offset: 5,
	}

	_, err := testQueries.ListAccounts(context.Background(), arg)
	require.Error(t, err)
}
