package api

import (
	"github.com/cvzm/go-web-project/doamin"

	"github.com/labstack/echo/v4"
)

type EventController struct {
	eventUsecase doamin.EventUsecase
}

func NewEventController(usecase doamin.EventUsecase) *EventController {
	return &EventController{
		eventUsecase: usecase,
	}
}

func (c *EventController) CreateAWSEvent(ctx echo.Context) error {
	return HandleRequest(ctx, func(param doamin.AWSEvent) (any, error) {
		return nil, c.eventUsecase.Save(param)
	})
}

func (c *EventController) CreateGCPEvent(ctx echo.Context) error {
	return HandleRequest(ctx, func(param doamin.GCPEvent) (any, error) {
		return nil, c.eventUsecase.Save(param)
	})
}

func SetupEventRoutes(e *echo.Echo, controller *EventController) {
	e.POST("/events/aws", controller.CreateAWSEvent)
	e.POST("/events/gcp", controller.CreateGCPEvent)
}
