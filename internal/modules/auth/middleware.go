package auth

import (
	"github.com/AnatoliyRib1/movie-reviews/internal/jwt"
	"github.com/AnatoliyRib1/movie-reviews/internal/modules/users"
	"github.com/labstack/echo/v4"
)

func Self(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userId")

		claims := jwt.GetClaims(c)
		if claims.Role == users.AdminRole || claims.Subject == userId {
			return next(c)
		}

		return echo.ErrForbidden
	}
}

func Editor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)
		if claims.Role == users.AdminRole || claims.Role == users.EditorRole {
			return next(c)
		}

		return echo.ErrForbidden
	}
}

func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)
		if claims.Role == users.AdminRole {
			return next(c)
		}

		return echo.ErrForbidden
	}
}
