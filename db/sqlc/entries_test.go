package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomEntry(t *testing.T, accountID int64) Entry {
	arg := CreateEntriesParams{
		AccountID: accountID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntries(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account.ID)

	// Retrieve the entry from the database
	entryFromDB, err := testQueries.GetEntries(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryFromDB)

	// Compare the retrieved entry with the original entry
	require.Equal(t, entry.ID, entryFromDB.ID)
	require.Equal(t, entry.AccountID, entryFromDB.AccountID)
	require.Equal(t, entry.Amount, entryFromDB.Amount)
	require.WithinDuration(t, entry.CreatedAt, entryFromDB.CreatedAt, time.Second)
}

func TestDeleteEntries(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account.ID)

	err := testQueries.DeleteEntries(context.Background(), entry.ID)
	require.NoError(t, err)

	// Retrieve the entry from the database
	entryFromDB, err := testQueries.GetEntries(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entryFromDB)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, account.ID)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entryFromDB, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entryFromDB, 5)

	for _, entry := range entryFromDB {
		require.NotEmpty(t, entry)
	}
}

func TestUpdateEntries(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account.ID)

	arg := UpdateEntriesParams{
		ID:     entry.ID,
		Amount: util.RandomMoney(),
	}

	entryFromDB, err := testQueries.UpdateEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entryFromDB)
	require.Equal(t, arg.Amount, entryFromDB.Amount)
}
