package api

import (
	"github.com/cvzm/go-web-project/domain"

	"github.com/labstack/echo/v4"
)

type EventController struct {
	eventUsecase domain.EventUsecase
}

func NewEventController(usecase domain.EventUsecase) *EventController {
	return &EventController{
		eventUsecase: usecase,
	}
}

func (c *EventController) CreateAWSEvent(ctx echo.Context) error {
	return HandleRequest(ctx, func(param domain.AWSEvent) (any, error) {
		return nil, c.eventUsecase.Save(param)
	})
}

func (c *EventController) CreateGCPEvent(ctx echo.Context) error {
	return HandleRequest(ctx, func(param domain.GCPEvent) (any, error) {
		return nil, c.eventUsecase.Save(param)
	})
}

func SetupEventRoutes(e *echo.Echo, controller *EventController) {
	e.POST("/events/aws", controller.CreateAWSEvent)
	e.POST("/events/gcp", controller.CreateGCPEvent)
}
