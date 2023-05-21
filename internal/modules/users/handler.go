package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h Handler) GetUsers(c echo.Context) error {
	return c.String(200, "not implemented")
}

func (h Handler) Delete(c echo.Context) error {
	var req DeleteRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return h.service.Delete(c.Request().Context(), req.UserId)
}

func (h Handler) Put(c echo.Context) error {
	var req PutRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return h.service.Put(c.Request().Context(), req.UserId, req.Bio)

}

type DeleteRequest struct {
	UserId int `param:"UserId"`
}
type PutRequest struct {
	UserId int    `param:"UserId"`
	Bio    string `param:"Bio"`
}
