package echox

import (
	"errors"
	"net/http"

	"github.com/AnatoliyRib1/movie-reviews/contracts"
	"github.com/AnatoliyRib1/movie-reviews/internal/log"

	"github.com/AnatoliyRib1/movie-reviews/internal/apperrors"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	var appError *apperrors.Error
	if !errors.As(err, &appError) {
		appError = apperrors.InternalWithoutStackTrace(err)
	}
	httpError := contracts.HTTPError{
		Message:    appError.SafeError(),
		IncidentID: appError.IncidentID,
	}
	logger := log.FromContext(c.Request().Context())

	if appError.Code == apperrors.InternalCode {
		logger.Error("server error", "message", err.Error(), "incidentId", appError.IncidentID, "stacktrace", appError.StackTrace)
	} else {
		logger.Error("client error", "message", err.Error())
	}
	if err = c.JSON(toHTTPStatus(appError.Code), httpError); err != nil {
		c.Logger().Error(err)
	}
}

func toHTTPStatus(code apperrors.Code) int {
	switch code {
	case apperrors.InternalCode:
		return http.StatusInternalServerError
	case apperrors.BadRequestCode:
		return http.StatusBadRequest
	case apperrors.NotFoundCode:
		return http.StatusNotFound
	case apperrors.AlreadyExistsCode, apperrors.VersionMismatchCode:
		return http.StatusConflict
	case apperrors.UnauthorizedCode:
		return http.StatusUnauthorized
	case apperrors.ForbiddenCode:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError

	}
}
