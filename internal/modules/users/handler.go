package users

import (
	"net/http"

	"github.com/AnatoliyRib1/movie-reviews/contracts"

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
	req, err := echox.BindAndValidate[contracts.DeleteUserRequest](c)
	if err != nil {
		return err
	}
	return h.service.Delete(c.Request().Context(), req.UserID)
}

func (h Handler) Update(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateUserRequest](c)
	if err != nil {
		return err
	}

	return h.service.Update(c.Request().Context(), req.UserID, *req.Bio)
}

func (h Handler) Get(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetUserRequest](c)
	if err != nil {
		return err
	}
	user, err := h.service.GetExistingUserByID(c.Request().Context(), req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func (h Handler) GetByUserName(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetUserByUserNameRequest](c)
	if err != nil {
		return err
	}
	user, err := h.service.GetExistingUserByUserName(c.Request().Context(), req.UserName)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func (h Handler) GetByUserEmail(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.LoginUserRequest](c)
	if err != nil {
		return err
	}
	user, err := h.service.GetExistingUserWithPasswordByEmail(c.Request().Context(), req.Email)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func (h Handler) SetRole(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.SetUserRoleRequest](c)
	if err != nil {
		return err
	}

	return h.service.SetRole(c.Request().Context(), req.UserID, req.Role)
}
