package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// newTestContext creates a new test context for Echo framework
// It takes the HTTP method, URL, request body, and Echo instance as parameters
// Returns an Echo context and httptest.ResponseRecorder for testing
func newTestContext(method, url string, body interface{}, e *echo.Echo) (echo.Context, *httptest.ResponseRecorder) {
	var reqBody string
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = string(jsonBody)
	}

	// Create a new HTTP request with the given method, URL and body
	req := httptest.NewRequest(method, url, strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)

	return c, resp
}

func TestNewServer(t *testing.T) {
	e := NewServer()
	assert.NotNil(t, e)
	assert.IsType(t, &echo.Echo{}, e)
}

func TestHandleRequest(t *testing.T) {
	e := echo.New()

	type TestParam struct {
		Name string `json:"name"`
	}

	handler := func(param TestParam) (any, error) {
		return map[string]string{"greeting": "Hello, " + param.Name}, nil
	}

	e.POST("/test", func(c echo.Context) error {
		return HandleRequest(c, handler)
	})

	t.Run("Successfully handle request", func(t *testing.T) {
		reqBody := `{"name":"John"}`
		req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := HandleRequest(c, handler)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response StandardResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{"greeting": "Hello, John"}, response.Data)
	})

	t.Run("Invalid request format", func(t *testing.T) {
		reqBody := `{"name":123}` // Invalid JSON
		req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := HandleRequest(c, handler)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response StandardResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request format", response.Message)
	})
}
