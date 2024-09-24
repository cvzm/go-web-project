package api

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/cvzm/go-web-project/doamin"
	domain_mock "github.com/cvzm/go-web-project/doamin/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewEventController(t *testing.T) {
	mockUsecase := new(domain_mock.MockEventUsecase)
	controller := NewEventController(mockUsecase)
	assert.NotNil(t, controller)
	assert.Equal(t, mockUsecase, controller.eventUsecase)
}

func TestEventController_CreateAWSEvent(t *testing.T) {
	mockUsecase := new(domain_mock.MockEventUsecase)
	controller := NewEventController(mockUsecase)
	e := echo.New()

	t.Run("Successfully create AWS event", func(t *testing.T) {
		awsEvent := doamin.AWSEvent{
			AWSEventID:   "aws-123",
			AWSEventType: "EC2_STARTED",
			AWSMessage:   "EC2 instance has started",
			AWSTimestamp: time.Now(),
		}

		mockUsecase.On("Save", mock.AnythingOfType("doamin.AWSEvent")).Return(nil).Once()

		c, resp := newTestContext(http.MethodPost, "/events/aws", awsEvent, e)

		err := controller.CreateAWSEvent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("Fail to create AWS event", func(t *testing.T) {
		awsEvent := doamin.AWSEvent{
			AWSEventID:   "aws-456",
			AWSEventType: "EC2_STOPPED",
			AWSMessage:   "EC2 instance has stopped",
			AWSTimestamp: time.Now(),
		}

		mockUsecase.On("Save", mock.AnythingOfType("doamin.AWSEvent")).Return(errors.New("creation failed")).Once()

		c, resp := newTestContext(http.MethodPost, "/events/aws", awsEvent, e)

		err := controller.CreateAWSEvent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)

		mockUsecase.AssertExpectations(t)
	})
}

func TestEventController_CreateGCPEvent(t *testing.T) {
	mockUsecase := new(domain_mock.MockEventUsecase)
	controller := NewEventController(mockUsecase)
	e := echo.New()

	t.Run("Successfully create GCP event", func(t *testing.T) {
		gcpEvent := doamin.GCPEvent{
			GCPEventID:   "gcp-123",
			GCPEventType: "VM_STARTED",
			GCPMessage:   "VM instance has started",
			GCPTimestamp: time.Now(),
		}

		mockUsecase.On("Save", mock.AnythingOfType("doamin.GCPEvent")).Return(nil).Once()

		c, resp := newTestContext(http.MethodPost, "/events/gcp", gcpEvent, e)

		err := controller.CreateGCPEvent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("Fail to create GCP event", func(t *testing.T) {
		gcpEvent := doamin.GCPEvent{
			GCPEventID:   "gcp-456",
			GCPEventType: "VM_STOPPED",
			GCPMessage:   "VM instance has stopped",
			GCPTimestamp: time.Now(),
		}

		mockUsecase.On("Save", mock.AnythingOfType("doamin.GCPEvent")).Return(errors.New("creation failed")).Once()

		c, resp := newTestContext(http.MethodPost, "/events/gcp", gcpEvent, e)

		err := controller.CreateGCPEvent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)

		mockUsecase.AssertExpectations(t)
	})
}

func TestSetupEventRoutes(t *testing.T) {
	e := echo.New()
	mockUsecase := new(domain_mock.MockEventUsecase)
	controller := NewEventController(mockUsecase)

	SetupEventRoutes(e, controller)

	assert.NotNil(t, e.Router().Routes())
	assert.Len(t, e.Router().Routes(), 2)

	routes := e.Router().Routes()
	awsRouteFound := false
	gcpRouteFound := false

	for _, route := range routes {
		switch route.Path {
		case "/events/aws":
			assert.Equal(t, http.MethodPost, route.Method)
			awsRouteFound = true
		case "/events/gcp":
			assert.Equal(t, http.MethodPost, route.Method)
			gcpRouteFound = true
		}
	}

	assert.True(t, awsRouteFound, "AWS route not found")
	assert.True(t, gcpRouteFound, "GCP route not found")
}
