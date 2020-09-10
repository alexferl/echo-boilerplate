package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"echo-boilerplate/internal/app/httpd/handlers"
)

func NewTestRequest(e *echo.Echo, handlerName, target string, body interface{}) (*handlers.Handler, *httptest.ResponseRecorder, *http.Request) {
	route := handlers.Router.FindRouteByName(handlerName)
	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(route.Method, target, bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	h := handlers.Register(e)

	return h, rec, req
}
