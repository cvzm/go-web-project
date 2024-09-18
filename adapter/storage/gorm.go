package storage

import (
	"gorm.io/gorm"
)

// DBConfig is db related configuration
type DBConfig struct {
	MaxIdleConns int // maximum number of connections in the idle connection pool
	MaxOpenConns int // maximum number of open connections to the database
}

// DBConnector is an interface that represents a database connector.
type DBConnector interface {
	// Connect connects to the database and returns a *gorm.DB instance and an error.
	Connect() (*gorm.DB, error)
}

// NewDB initializes the database connection based on the provided DBConfig.
// If there is an error connecting to the database, it returns an error.
func NewDB(connector DBConnector, conf DBConfig) (*gorm.DB, error) {
	// Connect to the database using the specified DBConnector.
	gormDB, err := connector.Connect()
	if err != nil {
		return nil, err
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings.
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetMaxOpenConns(conf.MaxOpenConns)

	// Check whether the connection is successful
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return gormDB, nil
}
