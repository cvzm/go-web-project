package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cvzm/go-web-project/adapter/storage"
	"github.com/cvzm/go-web-project/api"
	"github.com/cvzm/go-web-project/doamin"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// App struct represents the main application
type App struct {
	config *Config
	db     *gorm.DB
	echo   *echo.Echo

	eventController *api.EventController
}

// NewApp creates and returns a new App instance
func NewApp(cfg *Config, db *gorm.DB, e *echo.Echo, eventController *api.EventController) *App {
	return &App{
		config:          cfg,
		db:              db,
		echo:            e,
		eventController: eventController,
	}
}

// SetupAndRun sets up and runs the application
func (a *App) SetupAndRun() error {
	// Set up routes
	api.SetupEventRoutes(a.echo, a.eventController)

	// Start the server in a goroutine
	go func() {
		serverAddr := fmt.Sprintf(":%d", a.config.ServerPort)
		if err := a.echo.Start(serverAddr); err != nil {
			a.echo.Logger.Info("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Gracefully shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.echo.Shutdown(ctx); err != nil {
		return err
	}

	return a.Shutdown()
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown() error {
	// Add any necessary cleanup operations here
	// For example, closing database connections
	sqlDB, err := a.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// initDatabase initializes the database connection
func initDatabase(cfg *Config) (*gorm.DB, error) {
	connector := &storage.PostgresDBConnector{
		DSN:        cfg.DBDSN,
		ReplicaDSN: cfg.DBReplicaDSN,
	}
	db, err := storage.NewDB(connector, storage.DBConfig{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&doamin.Event{}); err != nil {
		return nil, err
	}

	return db, nil
}
