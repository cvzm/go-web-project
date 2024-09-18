package storage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MockDBConnector is a mock implementation of a DBConnector interface that returns a pre-configured GORM database instance.
type MockDBConnector struct {
	DB *gorm.DB
}

// Connect returns the pre-configured GORM database instance.
func (c *MockDBConnector) Connect() (*gorm.DB, error) {
	return c.DB, nil
}

// GetMockDB returns a GORM database instance and a sqlmock.Sqlmock for testing purposes.
func GetMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	// Create a new mock database connection.
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)

	// Set up a mock DBConnector instance to return the mock database connection.
	mock.ExpectPing()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	return gormDB, mock
}
