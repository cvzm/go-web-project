package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewServer initializes the Echo server and sets up middleware
func NewServer() *echo.Echo {
	e := echo.New()

	// Register basic middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: e.Logger.Output(),
	}))

	return e
}

// StandardResponse defines a simplified standard response structure
type StandardResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type (
	// RequestHandler is a function type for handling standard requests
	RequestHandler[T any] func(T) (any, error)
)

// HandleRequest processes standard requests
func HandleRequest[T any](c echo.Context, handler RequestHandler[T]) error {
	var param T
	// Bind request body to parameter
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, StandardResponse{
			Message: "Invalid request format",
		})
	}

	// Call the handler function
	result, err := handler(param)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, StandardResponse{
			Message: "Internal error",
		})
	}

	// Return successful response
	return c.JSON(http.StatusOK, StandardResponse{
		Data: result,
	})
}
