package postgresql

import (
	"context"
	"testing"

	"github.com/Zrossiz/gophermart/internal/dto"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBalanceHistoryStoreCreate(t *testing.T) {
	cleanUserDatabase(t, testDB)
	cleanStatusDatabase(t, testDB)
	cleanBalanceHistoryDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	userStore := NewUserStore(testDB, logger)
	balanceHistoryStore := NewBalanceHistoryStore(testDB, logger)

	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")

	user, err := userStore.GetUserByName("testuser1")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser1", user.Name)

	created, err = balanceHistoryStore.Create(dto.CreateBalanceHistory{
		UserID:  int64(user.ID),
		OrderID: "123",
		Change:  1,
	})
	assert.NoError(t, err)
	assert.True(t, created)
}

func TestBalanceHistoryStoreGetAllDebits(t *testing.T) {
	cleanUserDatabase(t, testDB)
	cleanStatusDatabase(t, testDB)
	cleanBalanceHistoryDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	userStore := NewUserStore(testDB, logger)
	balanceHistoryStore := NewBalanceHistoryStore(testDB, logger)

	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")

	user, err := userStore.GetUserByName("testuser1")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser1", user.Name)

	created, err = balanceHistoryStore.Create(dto.CreateBalanceHistory{
		UserID:  int64(user.ID),
		OrderID: "123",
		Change:  1,
	})
	assert.NoError(t, err)
	assert.True(t, created)

	created, err = balanceHistoryStore.Create(dto.CreateBalanceHistory{
		UserID:  int64(user.ID),
		OrderID: "1234",
		Change:  1,
	})
	assert.NoError(t, err)
	assert.True(t, created)

	allDebits, err := balanceHistoryStore.GetAllDebits(int64(user.ID))
	assert.NoError(t, err)
	assert.NotNil(t, allDebits)
	assert.Len(t, allDebits, 2)
}

func TestBalanceHistoryStoreWithdraw(t *testing.T) {
	cleanUserDatabase(t, testDB)
	cleanStatusDatabase(t, testDB)
	cleanBalanceHistoryDatabase(t, testDB)
	cleanOrderDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	userStore := NewUserStore(testDB, logger)
	balanceHistoryStore := NewBalanceHistoryStore(testDB, logger)

	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")

	user, err := userStore.GetUserByName("testuser1")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser1", user.Name)

	updatedBalance, err := userStore.UpdateUserBalance(int64(user.ID), decimal.NewFromFloat(100.50))
	assert.NoError(t, err)
	assert.True(t, updatedBalance)

	err = balanceHistoryStore.Withdraw(user.ID, "123", 20.00)
	assert.NoError(t, err)
}

func cleanBalanceHistoryDatabase(t *testing.T, db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), "DELETE FROM balance_history")
	require.NoError(t, err)
}
