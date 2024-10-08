// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package bootstrap

import (
	"github.com/cvzm/go-web-project/api"
	"github.com/cvzm/go-web-project/repository"
	"github.com/cvzm/go-web-project/usecase"
)

// Injectors from wire.go:

// InitializeApp initializes the application using the Wire framework
func InitializeApp() (*App, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}
	db, err := initDatabase(config)
	if err != nil {
		return nil, err
	}
	echo := api.NewServer()
	eventRepository := repository.NewEventRepository(db)
	eventUsecase := usecase.NewEventUsecase(eventRepository)
	sqsConsumer, err := initSQSConsumer(config, eventUsecase)
	if err != nil {
		return nil, err
	}
	eventController := api.NewEventController(eventUsecase)
	app := NewApp(config, db, echo, sqsConsumer, eventController)
	return app, nil
}
