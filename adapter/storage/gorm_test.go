package storage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MockDBConnector struct {
	DB *gorm.DB
}

// Connect returns the pre-configured GORM database instance.
func (c *MockDBConnector) Connect() (*gorm.DB, error) {
	return c.DB, nil
}

func TestNewDB(t *testing.T) {
	db, mock := GetMockDB(t)
	connector := &MockDBConnector{DB: db}

	// Set up the expected configuration options.
	conf := DBConfig{
		MaxIdleConns: 10,
		MaxOpenConns: 20,
	}

	// Initialize the GORM database connection using the mock DBConnector instance.
	mock.ExpectPing()
	gormDB, err := NewDB(connector, conf)
	assert.NoError(t, err)

	actualDB, err := gormDB.DB()
	assert.NoError(t, err)
	assert.Equal(t, conf.MaxOpenConns, actualDB.Stats().MaxOpenConnections)

	assert.NoError(t, mock.ExpectationsWereMet())
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
