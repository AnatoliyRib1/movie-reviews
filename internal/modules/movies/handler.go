package movies

import (
	"net/http"

	"github.com/AnatoliyRib1/movie-reviews/contracts"
	"github.com/AnatoliyRib1/movie-reviews/internal/config"
	"github.com/AnatoliyRib1/movie-reviews/internal/echox"
	"github.com/AnatoliyRib1/movie-reviews/internal/pagination"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service          *Service
	paginationConfig config.PaginationConfig
}

func NewHandler(service *Service, paginationConfig config.PaginationConfig) *Handler {
	return &Handler{
		service:          service,
		paginationConfig: paginationConfig,
	}
}

func (h *Handler) GetAll(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetMoviesRequest](c)
	if err != nil {
		return err
	}
	pagination.SetDefaults(&req.PaginatedRequest, h.paginationConfig)
	offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)

	movies, total, err := h.service.GetAllPaginated(c.Request().Context(), offset, limit)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, pagination.Response(&req.PaginatedRequest, total, movies))
}

func (h *Handler) Get(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetMovieRequest](c)
	if err != nil {
		return err
	}
	movie, err := h.service.GetByID(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, movie)
}

func (h *Handler) Create(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.CreateMovieRequest](c)
	if err != nil {
		return err
	}
	movie := &Movie{
		Title:       req.Title,
		Description: req.Description,
		ReleaseDate: req.ReleaseDate,
	}
	err = h.service.Create(c.Request().Context(), movie)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, movie)
}

func (h *Handler) Delete(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.DeleteMovieRequest](c)
	if err != nil {
		return err
	}
	if err = h.service.Delete(c.Request().Context(), req.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) Update(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateMovieRequest](c)
	if err != nil {
		return err
	}
	movie := &Movie{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		ReleaseDate: req.ReleaseDate,
	}

	if err = h.service.Update(c.Request().Context(), movie); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
