package auth

import (
	"github.com/AnatoliyRib1/movie-reviews/internal/modules/users"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	authService *Service
}

func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	user := &users.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := h.authService.Register(c.Request().Context(), user, req.Password); err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c echo.Context) error {
	var logReq LoginRequest

	if err := c.Bind(&logReq); err != nil {
		return err
	}

	token, err := h.authService.Login(c.Request().Context(), logReq.Email, logReq.Password)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
	}
	return c.JSON(http.StatusOK, token)
}

func NewHandler(authService *Service) *Handler {
	return &Handler{authService: authService}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
