package database

import (
	"github.com/labstack/echo/v4"
)

// DBContext middleware makes db available in request context.
func DBContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("db", db)
		return next(c)
	}
}
