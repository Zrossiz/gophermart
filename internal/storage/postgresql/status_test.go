package postgresql

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestStatusStoreCreate(t *testing.T) {
	cleanStatusDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	statusStore := NewStatusStore(testDB, logger)

	created, err := statusStore.Create("new")
	require.NoError(t, err)
	assert.True(t, created, "status should be created successfully")
}

func TestStatusStoreGetAll(t *testing.T) {
	cleanStatusDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	statusStore := NewStatusStore(testDB, logger)

	created, err := statusStore.Create("new")
	require.NoError(t, err)
	assert.True(t, created, "status should be created successfully")

	statuses, err := statusStore.GetAll()
	require.NoError(t, err)
	require.NotNil(t, statuses)
}

func cleanStatusDatabase(t *testing.T, db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), "DELETE FROM statuses")
	require.NoError(t, err)
}
