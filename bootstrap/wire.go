//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/cvzm/go-web-project/api"
	"github.com/cvzm/go-web-project/repository"
	"github.com/cvzm/go-web-project/usecase"

	"github.com/google/wire"
)

// InitializeApp initializes the application using the Wire framework
func InitializeApp() (*App, error) {
	wire.Build(
		// Initialize configuration
		NewConfig,

		// Initialize database connection
		initDatabase,

		// Initialize SQS consumer
		initSQSConsumer,

		// Create API server instance
		api.NewServer,

		// Create event repository instance
		repository.NewEventRepository,

		// Create event usecase instance
		usecase.NewEventUsecase,

		// Create event controller instance
		api.NewEventController,

		// Create and return App instance
		NewApp,
	)

	// wire.NewSet(NewConfig, NewEventUsecase, SetupSQSConsumer)
	return &App{}, nil
}
