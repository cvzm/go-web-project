package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
