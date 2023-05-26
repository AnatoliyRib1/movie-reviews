package users

import (
	"net/http"

	"github.com/AnatoliyRib1/movie-reviews/internal/apperrors"
	"github.com/AnatoliyRib1/movie-reviews/internal/echox"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h Handler) Delete(c echo.Context) error {
	req, err := echox.BindAndValidate[DeleteOrGetRequest](c)
	if err != nil {
		return err
	}
	return h.service.Delete(c.Request().Context(), req.UserId)
}

func (h Handler) Update(c echo.Context) error {
	req, err := echox.BindAndValidate[UpdateUserRequest](c)
	if err != nil {
		return err
	}

	return h.service.UpdateBio(c.Request().Context(), req.UserId, *req.Bio)
}

func (h Handler) Get(c echo.Context) error {
	req, err := echox.BindAndValidate[DeleteOrGetRequest](c)
	if err != nil {
		return err
	}
	user, err := h.service.Get(c.Request().Context(), req.UserId)
	if err != nil {
		return apperrors.BadRequest(err)
	}
	return c.JSON(http.StatusOK, user)
}

func (h Handler) SetRole(c echo.Context) error {
	req, err := echox.BindAndValidate[SetUserRoleRequest](c)
	if err != nil {
		return err
	}

	return h.service.SetRole(c.Request().Context(), req.UserId, req.Role)
}

type DeleteOrGetRequest struct {
	UserId int `param:"userId" validate:"nonzero"`
}
type UpdateUserRequest struct {
	UserId int     `param:"userId" validate:"nonzero"`
	Bio    *string `json:"bio"`
}

type SetUserRoleRequest struct {
	UserId int    `param:"userId" validate:"nonzero"`
	Role   string `param:"role" validate:"role"`
}
