package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/AnatoliyRib1/movie-reviews/internal/reviews"

	"github.com/AnatoliyRib1/movie-reviews/internal/modules/movies"

	"github.com/AnatoliyRib1/movie-reviews/internal/modules/stars"

	"github.com/AnatoliyRib1/movie-reviews/internal/modules/genres"

	"github.com/AnatoliyRib1/movie-reviews/internal/apperrors"
	"github.com/AnatoliyRib1/movie-reviews/internal/config"
	"github.com/AnatoliyRib1/movie-reviews/internal/echox"
	"github.com/AnatoliyRib1/movie-reviews/internal/jwt"
	"github.com/AnatoliyRib1/movie-reviews/internal/log"
	"github.com/AnatoliyRib1/movie-reviews/internal/modules/auth"
	"github.com/AnatoliyRib1/movie-reviews/internal/modules/users"
	"github.com/AnatoliyRib1/movie-reviews/internal/validation"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
	"gopkg.in/validator.v2"
)

const (
	dbConnectionTimeout  = 10 * time.Second
	adminCreationTimeout = 5 * time.Second
)

type Server struct {
	e       *echo.Echo
	cfg     *config.Config
	closers []func() error
}

func New(ctx context.Context, cfg *config.Config) (*Server, error) {
	logger, err := log.SetupLogger(cfg.Local, cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("setup logger: %w", err)
	}
	slog.SetDefault(logger)

	validation.SetupValidators()

	var closers []func() error
	db, err := getDb(ctx, cfg.DbURL)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}
	closers = append(closers, func() error { db.Close(); return nil })

	jwtService := jwt.NewService(cfg.Jwt.Secret, cfg.Jwt.AccessExpiration)
	usersModule := users.NewModule(db)
	authModule := auth.NewModule(jwtService, usersModule.Service)
	genresModule := genres.NewModule(db)
	starsModule := stars.NewModule(db, cfg.Pagination)
	moviesModule := movies.NewModule(db, genresModule, starsModule, cfg.Pagination)
	reviewsModule := reviews.NewModule(db, moviesModule, cfg.Pagination)

	if err = createInitialAdminUser(cfg.Admin, authModule.Service); err != nil {
		return nil, withClosers(closers, fmt.Errorf("create initial admin user: %w", err))
	}

	e := echo.New()
	e.HTTPErrorHandler = echox.ErrorHandler

	e.Use(middleware.Recover())
	e.HideBanner = true
	e.HidePort = true

	api := e.Group("/api")
	api.Use(jwt.NewAuthMiddleware(cfg.Jwt.Secret))
	api.Use(echox.Logger)

	// Auth API routes
	api.POST("/auth/register", authModule.Handler.Register)
	api.POST("/auth/login", authModule.Handler.Login)

	// Users API routes
	api.GET("/users/:userId", usersModule.Handler.Get)
	api.GET("/users/username/:username", usersModule.Handler.GetByUserName)
	api.PUT("/users/:userId", usersModule.Handler.Update, auth.Self)
	api.DELETE("/users/:userId", usersModule.Handler.Delete, auth.Self)
	api.PUT("/users/:userId/role/:role", usersModule.Handler.SetRole, auth.Admin)

	// Genres API routes
	api.GET("/genres", genresModule.Handler.GetAll)
	api.GET("/genres/:genreId", genresModule.Handler.Get)
	api.POST("/genres", genresModule.Handler.Create, auth.Editor)
	api.PUT("/genres/:genreId", genresModule.Handler.Update, auth.Editor)
	api.DELETE("/genres/:genreId", genresModule.Handler.Delete, auth.Editor)

	// Stars API routes
	api.GET("/stars", starsModule.Handler.GetAll)
	api.GET("/stars/:starId", starsModule.Handler.Get)
	api.POST("/stars", starsModule.Handler.Create, auth.Editor)
	api.PUT("/stars/:starId", starsModule.Handler.Update, auth.Editor)
	api.DELETE("/stars/:starId", starsModule.Handler.Delete, auth.Editor)

	// Movies API routes
	api.GET("/movies", moviesModule.Handler.GetAll)
	api.GET("/movies/:movieId", moviesModule.Handler.Get)
	api.POST("/movies", moviesModule.Handler.Create, auth.Editor)
	api.PUT("/movies/:movieId", moviesModule.Handler.Update, auth.Editor)
	api.DELETE("/movies/:movieId", moviesModule.Handler.Delete, auth.Editor)

	// Reviews API routes
	api.GET("/reviews", reviewsModule.Handler.GetAll)
	api.GET("/reviews/:reviewId", reviewsModule.Handler.Get)
	api.POST("/users/:userId/reviews", reviewsModule.Handler.Create, auth.Self)
	api.PUT("/users/:userId/reviews/:reviewId", reviewsModule.Handler.Update, auth.Self)
	api.DELETE("/users/:userId/reviews/:reviewId", reviewsModule.Handler.Delete, auth.Self)

	return &Server{e: e, cfg: cfg, closers: closers}, nil
}

func (s *Server) Start() error {
	port := s.cfg.Port
	slog.Info("server started", "port", port)
	return s.e.Start(fmt.Sprintf(":%d", port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) Close() error {
	return withClosers(s.closers, nil)
}

func (s *Server) Port() (int, error) {
	listener := s.e.Listener
	if listener == nil {
		return 0, errors.New("server is not started")
	}

	addr := listener.Addr()
	if addr == nil {
		return 0, errors.New("server is not started")
	}

	return addr.(*net.TCPAddr).Port, nil
}

func getDb(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, dbConnectionTimeout)
	defer cancel()

	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}
	return db, nil
}

func createInitialAdminUser(cfg config.AdminConfig, authService *auth.Service) error {
	if !cfg.IsSet() {
		return nil
	}

	if err := validator.Validate(cfg); err != nil {
		return fmt.Errorf("validate admin config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), adminCreationTimeout)
	defer cancel()

	err := authService.Register(ctx, &users.User{
		Username: cfg.Username,
		Email:    cfg.Email,
		Role:     users.AdminRole,
	}, cfg.Password)

	switch {
	case apperrors.Is(err, apperrors.InternalCode):
		return fmt.Errorf("register admin: %w", err)
	case err != nil:
		// just ignore all other errors
		return nil
	default:
		slog.Info("admin user created", "username", cfg.Username, "email", cfg.Email)
		return nil
	}
}

func withClosers(closers []func() error, err error) error {
	errs := []error{err}

	for i := len(closers) - 1; i >= 0; i-- {
		if err = closers[i](); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
