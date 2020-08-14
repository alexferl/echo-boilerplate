package healthz

import (
	"github.com/labstack/echo/v4"
)

// Healthz is for load-balancer/Kubernetes health checks
func Healthz(c echo.Context) error {
	return c.String(200, "ok")
}
