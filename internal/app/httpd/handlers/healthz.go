package handlers

import "github.com/labstack/echo/v4"

// Healthz is for load-balancer/Kubernetes health checks
func (h *Handler) Healthz(c echo.Context) error {
	return c.String(200, "ok")
}
