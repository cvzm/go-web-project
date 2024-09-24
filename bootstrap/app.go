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
	config      *Config
	db          *gorm.DB
	echo        *echo.Echo
	sqsConsumer *SQSConsumer

	eventController *api.EventController
}

// NewApp creates and returns a new App instance
func NewApp(cfg *Config, db *gorm.DB, e *echo.Echo, sqsConsumer *SQSConsumer, eventController *api.EventController) *App {
	return &App{
		config:          cfg,
		db:              db,
		echo:            e,
		sqsConsumer:     sqsConsumer,
		eventController: eventController,
	}
}

// SetupAndRun sets up and runs the application
func (a *App) SetupAndRun() error {
	a.setupRoutes()

	go a.startServer()
	go a.sqsConsumer.Start()

	return a.gracefulShutdown()
}

// setupRoutes sets up all routes
func (a *App) setupRoutes() {
	api.SetupEventRoutes(a.echo, a.eventController)
}

// startServer starts the server in the background
func (a *App) startServer() {
	serverAddr := fmt.Sprintf(":%d", a.config.ServerPort)
	if err := a.echo.Start(serverAddr); err != nil {
		a.echo.Logger.Info("Server is shutting down")
	}
}

// gracefulShutdown gracefully shuts down the server
func (a *App) gracefulShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.echo.Shutdown(ctx); err != nil {
		return err
	}

	return a.closeDB()
}

// closeDB closes the database connection
func (a *App) closeDB() error {
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

	if err := db.AutoMigrate(getModelsToMigrate()...); err != nil {
		return nil, err
	}

	return db, nil
}

func getModelsToMigrate() []any {
	return []any{
		&doamin.Event{},
		// Add other models here
	}
}
