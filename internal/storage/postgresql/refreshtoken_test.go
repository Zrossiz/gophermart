package postgresql

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRefreshTokenStoreCreate(t *testing.T) {
	clenRefreshTokenDatabase(t, testDB)
	cleanUserDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	refreshTokenStore := NewTokenStore(testDB, logger)
	userStore := NewUserStore(testDB, logger)

	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")

	user, err := userStore.GetUserByName("testuser1")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser1", user.Name)

	createdToken, err := refreshTokenStore.Create(int64(user.ID), "123")
	require.NoError(t, err)
	assert.True(t, createdToken, "refresh token should be created successfully")
}

func TestRefreshTokenStoreGetTokenByToken(t *testing.T) {
	clenRefreshTokenDatabase(t, testDB)
	cleanUserDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	refreshTokenStore := NewTokenStore(testDB, logger)
	userStore := NewUserStore(testDB, logger)

	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")

	user, err := userStore.GetUserByName("testuser1")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser1", user.Name)

	createdToken, err := refreshTokenStore.Create(int64(user.ID), "123")
	require.NoError(t, err)
	assert.True(t, createdToken, "refresh token should be created successfully")

	token, err := refreshTokenStore.GetTokenByToken("123")
	require.NoError(t, err)
	require.NotNil(t, token)
}

func TestRefreshTokenStoreGetTokenByUser(t *testing.T) {
	clenRefreshTokenDatabase(t, testDB)
	cleanUserDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	refreshTokenStore := NewTokenStore(testDB, logger)
	userStore := NewUserStore(testDB, logger)

	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")

	user, err := userStore.GetUserByName("testuser1")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser1", user.Name)

	createdToken, err := refreshTokenStore.Create(int64(user.ID), "123")
	require.NoError(t, err)
	assert.True(t, createdToken, "refresh token should be created successfully")

	token, err := refreshTokenStore.GetTokenByUser(int64(user.ID))
	require.NoError(t, err)
	require.NotNil(t, token)
}

func TestRefreshTokenStoreDeleteByToken(t *testing.T) {
	clenRefreshTokenDatabase(t, testDB)
	cleanUserDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	refreshTokenStore := NewTokenStore(testDB, logger)
	userStore := NewUserStore(testDB, logger)

	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")

	user, err := userStore.GetUserByName("testuser1")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser1", user.Name)

	createdToken, err := refreshTokenStore.Create(int64(user.ID), "123")
	require.NoError(t, err)
	assert.True(t, createdToken, "refresh token should be created successfully")

	deleted, err := refreshTokenStore.DeleteByToken("123")
	require.NoError(t, err)
	assert.True(t, deleted, "refresh token delete should be true")
}

func TestRefreshTokenStoreDeleteteByUser(t *testing.T) {
	clenRefreshTokenDatabase(t, testDB)
	cleanUserDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	refreshTokenStore := NewTokenStore(testDB, logger)
	userStore := NewUserStore(testDB, logger)

	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")

	user, err := userStore.GetUserByName("testuser1")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser1", user.Name)

	createdToken, err := refreshTokenStore.Create(int64(user.ID), "123")
	require.NoError(t, err)
	assert.True(t, createdToken, "refresh token should be created successfully")

	deleted, err := refreshTokenStore.DeleteTokenByUser(int64(user.ID))
	require.NoError(t, err)
	assert.True(t, deleted, "refresh token delete should be true")
}

func clenRefreshTokenDatabase(t *testing.T, db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), "DELETE FROM refresh_tokens")
	require.NoError(t, err)
}
