package auth

import (
	"net/http"

	"github.com/AnatoliyRib1/movie-reviews/internal/apperrors"
	"github.com/AnatoliyRib1/movie-reviews/internal/echox"

	"github.com/AnatoliyRib1/movie-reviews/internal/modules/users"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	authService *Service
}

func (h *Handler) Register(c echo.Context) error {
	req, err := echox.BindAndValidate[RegisterRequest](c)
	if err != nil {
		return err
	}

	user := &users.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     users.UserRole,
	}

	if err := h.authService.Register(c.Request().Context(), user, req.Password); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c echo.Context) error {
	req, err := echox.BindAndValidate[LoginRequest](c)
	if err != nil {
		return err
	}

	token, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return apperrors.BadRequestHidden(err, "user login error")
	}
	return c.JSON(http.StatusOK, LoginResponse{AccessToken: token})
}

func NewHandler(authService *Service) *Handler {
	return &Handler{authService: authService}
}

type RegisterRequest struct {
	Username string `json:"username" validate:"min=5,max=16"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"password"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
