package config

import (
	"github.com/hadian90/ping-service/obj"
	"github.com/labstack/echo"
)

// UserMiddleware ...
func UserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// check if header has user id if not return 403
		uID := c.Request().Header.Get("user_id")

		if uID != "" {
			return next(c)
		}

		return c.JSON(403, obj.Response{
			Success: false,
			Message: "Missing user data",
		})
	}
}
