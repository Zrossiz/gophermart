package postgresql

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestOrderStoreCreate(t *testing.T) {
	cleanOrderDatabase(t, testDB)
	cleanUserDatabase(t, testDB)
	cleanStatusDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	userStore := NewUserStore(testDB, logger)
	orderStore := NewOrderStore(testDB, logger)
	statusStore := NewStatusStore(testDB, logger)

	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")

	user, err := userStore.GetUserByName("testuser1")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser1", user.Name)

	createdStatus, err := statusStore.Create("new")
	require.NoError(t, err)
	assert.True(t, createdStatus, "status should be created successfully")

	createdOrder, err := orderStore.CreateOrder("123", user.ID)
	require.NoError(t, err)
	require.True(t, createdOrder)
}

func TestOrderStoreGetOrderByID(t *testing.T) {
	cleanOrderDatabase(t, testDB)
	cleanUserDatabase(t, testDB)
	cleanStatusDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	userStore := NewUserStore(testDB, logger)
	orderStore := NewOrderStore(testDB, logger)
	statusStore := NewStatusStore(testDB, logger)

	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")

	user, err := userStore.GetUserByName("testuser1")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser1", user.Name)

	createdStatus, err := statusStore.Create("new")
	require.NoError(t, err)
	assert.True(t, createdStatus, "status should be created successfully")

	createdOrder, err := orderStore.CreateOrder("123", user.ID)
	require.NoError(t, err)
	require.True(t, createdOrder)

	order, err := orderStore.GetOrderByID("123")
	require.NoError(t, err)
	require.NotNil(t, order)
}

func TestOrderStoreUpdateSumAndStatusOrder(t *testing.T) {
	cleanOrderDatabase(t, testDB)
	cleanUserDatabase(t, testDB)
	cleanStatusDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	userStore := NewUserStore(testDB, logger)
	orderStore := NewOrderStore(testDB, logger)
	statusStore := NewStatusStore(testDB, logger)

	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")

	user, err := userStore.GetUserByName("testuser1")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser1", user.Name)

	createdStatus, err := statusStore.Create("new")
	require.NoError(t, err)
	assert.True(t, createdStatus, "status should be created successfully")

	createdOrder, err := orderStore.CreateOrder("123", user.ID)
	require.NoError(t, err)
	require.True(t, createdOrder)

	updatedOrder, err := orderStore.UpdateSumAndStatusOrder("123", "new", 1.2)
	require.NoError(t, err)
	assert.True(t, updatedOrder)
}

func TestOrderStoreGetAllOrdersByUser(t *testing.T) {
	cleanOrderDatabase(t, testDB)
	cleanUserDatabase(t, testDB)
	cleanStatusDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	userStore := NewUserStore(testDB, logger)
	orderStore := NewOrderStore(testDB, logger)
	statusStore := NewStatusStore(testDB, logger)

	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")

	user, err := userStore.GetUserByName("testuser1")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser1", user.Name)

	createdStatus, err := statusStore.Create("new")
	require.NoError(t, err)
	assert.True(t, createdStatus, "status should be created successfully")

	createdOrder1, err := orderStore.CreateOrder("123", user.ID)
	require.NoError(t, err)
	require.True(t, createdOrder1)

	createdOrder2, err := orderStore.CreateOrder("1234", user.ID)
	require.NoError(t, err)
	require.True(t, createdOrder2)

	orders, err := orderStore.GetAllOrdersByUser(int64(user.ID))
	require.NoError(t, err)
	require.NotNil(t, orders)
	require.Len(t, orders, 2)
}

func TestOrderStoreGetAllUnhandlerOrders(t *testing.T) {
	cleanOrderDatabase(t, testDB)
	cleanUserDatabase(t, testDB)
	cleanStatusDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	userStore := NewUserStore(testDB, logger)
	orderStore := NewOrderStore(testDB, logger)
	statusStore := NewStatusStore(testDB, logger)

	created, err := userStore.Create("testuser1", "testpassword")
	require.NoError(t, err)
	assert.True(t, created, "user should be created successfully")

	user, err := userStore.GetUserByName("testuser1")
	require.NoError(t, err)
	assert.NotNil(t, user, "user should be found")
	assert.Equal(t, "testuser1", user.Name)

	createdStatus1, err := statusStore.Create("new")
	require.NoError(t, err)
	assert.True(t, createdStatus1, "status should be created successfully")

	createdStatus2, err := statusStore.Create("processing")
	require.NoError(t, err)
	assert.True(t, createdStatus2, "status should be created successfully")

	createdStatus3, err := statusStore.Create("invalid")
	require.NoError(t, err)
	assert.True(t, createdStatus3, "status should be created successfully")

	createdOrder1, err := orderStore.CreateOrder("123", user.ID)
	require.NoError(t, err)
	require.True(t, createdOrder1)

	createdOrder2, err := orderStore.CreateOrder("1234", user.ID)
	require.NoError(t, err)
	require.True(t, createdOrder2)

	createdOrder3, err := orderStore.CreateOrder("12345", user.ID)
	require.NoError(t, err)
	require.True(t, createdOrder3)

	updatedOrder1, err := orderStore.UpdateSumAndStatusOrder("123", "processing", 0)
	require.NoError(t, err)
	assert.True(t, updatedOrder1)

	updatedOrder2, err := orderStore.UpdateSumAndStatusOrder("12345", "invalid", 0)
	require.NoError(t, err)
	assert.True(t, updatedOrder2)

	statuses, err := statusStore.GetAll()
	require.NoError(t, err)
	require.Len(t, statuses, 3)

	var unhandledStatuses []int

	for _, status := range statuses {
		if status.Status == "new" || status.Status == "processing" {
			unhandledStatuses = append(unhandledStatuses, status.ID)
		}
	}

	unhandledOrders, err := orderStore.GetAllUnhandlerOrders(unhandledStatuses[0], unhandledStatuses[1])
	require.NoError(t, err)
	require.Len(t, unhandledOrders, 2)
}

func cleanOrderDatabase(t *testing.T, db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), "DELETE FROM orders")
	require.NoError(t, err)
}
