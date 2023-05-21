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

func (h Handler) Delete(c echo.Context) error {
	var req DeleteOrGetRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return h.service.Delete(c.Request().Context(), req.UserId)
}

func (h Handler) Update(c echo.Context) error {
	var req PutRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return h.service.Update(c.Request().Context(), req.UserId, req.Bio)
}

func (h Handler) Get(c echo.Context) error {
	var req DeleteOrGetRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user, err := h.service.Get(c.Request().Context(), req.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

type DeleteOrGetRequest struct {
	UserId int `param:"userId"`
}
type PutRequest struct {
	UserId int    `param:"userId" validate:"nonzero"`
	Bio    string `json:"bio"`
}
