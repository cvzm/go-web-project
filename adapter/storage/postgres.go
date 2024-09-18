package storage

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// PostgresDBConnector is an implementation of the DBConnector interface that connects
// to a real database for PostgreSQL.
type PostgresDBConnector struct {
	DSN        string
	ReplicaDSN string
}

// Connect connects to the database using the DSN provided in the PostgresDBConnector
// instance and returns a *gorm.DB instance and an error
func (c *PostgresDBConnector) Connect() (*gorm.DB, error) {
	// Open a new GORM database connection using the PostgreSQL dialector and configuration options.
	// The PrepareStmt option enables automatic prepared statement caching, which can improve performance.
	gormDB, err := gorm.Open(postgres.Open(c.DSN), &gorm.Config{
		PrepareStmt:     true,
		CreateBatchSize: 1000,
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	})
	if err != nil {
		return gormDB, err
	}
	if c.ReplicaDSN != "" {
		err = gormDB.Use(dbresolver.Register(dbresolver.Config{

			// Read-only sql will use this replicas
			Replicas: []gorm.Dialector{postgres.Open(c.ReplicaDSN)},

			// print sources/replicas mode in logger
			TraceResolverMode: true,
		}))
	}
	return gormDB, err
}
