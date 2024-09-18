package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func TestPostgresDBConnector_Connect(t *testing.T) {
	connector := &PostgresDBConnector{
		DSN:        "dsn",
		ReplicaDSN: "rep DNS",
	}
	db, err := connector.Connect()
	assert.ErrorContains(t, err, "cannot parse `dsn`")

	// check dialector name
	assert.Equal(t, db.Dialector.Name(), "postgres")
	assert.True(t, hasReplica(db))
}

func hasReplica(db *gorm.DB) bool {
	return db.Use(dbresolver.Register(dbresolver.Config{})) == nil
}
