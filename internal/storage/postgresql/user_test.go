package postgresql

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	testDB, err = Connect("postgres://postgres:root@localhost:5432/testdb?sslmode=disable")
	if err != nil {
		panic("failed to connect to test database")
	}
	defer testDB.Close()

	m.Run()
}

func TestUserStoreCreate(t *testing.T) {
	cleanUserDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	userStore := NewUserStore(testDB, logger)

	// Попытка создать пользователя
	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")
}

func TestUserStoreGetUserByName(t *testing.T) {
	cleanUserDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	userStore := NewUserStore(testDB, logger)

	_, err := userStore.Create("testuser0", "testpassword")
	require.NoError(t, err)

	user, err := userStore.GetUserByName("testuser0")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser0", user.Name)
}

func TetsUserStoreUpdateUserBalance(t *testing.T) {
	cleanUserDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	userStore := NewUserStore(testDB, logger)

	_, err := userStore.Create("balanceuser", "testpassword")
	require.NoError(t, err)

	user, err := userStore.GetUserByName("balanceuser")
	require.NoError(t, err)
	require.NotNil(t, user)

	newBalance := decimal.NewFromFloat(100.50)
	updated, err := userStore.UpdateUserBalance(int64(user.ID), newBalance)
	require.NoError(t, err)
	assert.True(t, updated, "balance should be updated successfully")

	updatedUser, err := userStore.GetUserByName("balanceuser")
	require.NoError(t, err)
	assert.Equal(t, newBalance, updatedUser.Account, "user balance should match the new balance")
}

func cleanUserDatabase(t *testing.T, db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), "DELETE FROM users")
	require.NoError(t, err)
}
