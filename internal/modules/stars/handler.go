package stars

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

/*
	func (h *Handler) GetAll(c echo.Context) error {
		genres, err := h.service.GetAll(c.Request().Context())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, genres)
	}
*/
func (h *Handler) Get(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetStarRequest](c)
	if err != nil {
		return err
	}
	star, err := h.service.GetByID(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, star)
}

func (h *Handler) Create(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.CreateStarRequest](c)
	if err != nil {
		return err
	}
	star := &Star{
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		LastName:   req.LastName,
		BirthDate:  req.BirthDate,
		BirthPlace: req.BirthPlace,
		DeathDate:  req.DeathDate,
		Bio:        req.Bio,
	}
	err = h.service.Create(c.Request().Context(), star)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, star)
}

/*
func (h *Handler) Delete(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.DeleteGenreRequest](c)
	if err != nil {
		return err
	}
	if err = h.service.Delete(c.Request().Context(), req.GenreID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) Update(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateGenreRequest](c)
	if err != nil {
		return err
	}

	if err = h.service.Update(c.Request().Context(), req.GenreID, req.Name); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

*/
