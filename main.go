package main

import (
	"context"
	"fmt"
	"github.com/AnatoliyRib1/movie-reviews/internal/log"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AnatoliyRib1/movie-reviews/internal/apperrors"
	"github.com/AnatoliyRib1/movie-reviews/internal/echox"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/validator.v2"

	"github.com/AnatoliyRib1/movie-reviews/internal/config"
	"github.com/AnatoliyRib1/movie-reviews/internal/jwt"
	"github.com/AnatoliyRib1/movie-reviews/internal/modules/auth"
	"github.com/AnatoliyRib1/movie-reviews/internal/modules/users"
	"github.com/AnatoliyRib1/movie-reviews/internal/validation"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

const (
	dbConnectTimeout     = 10 * time.Second
	adminCreationTimeout = 5 * time.Second
)

func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "parse config")
	validation.SetupValidators()

	logger, err := log.SetupLogger(cfg.Local, cfg.LogLevel)
	failOnError(err, "setup logger")
	slog.SetDefault(logger)

	slog.Info("started", "config", cfg)

	db, err := getDB(context.Background(), cfg.DbUrl)
	failOnError(err, "connect to db")

	jwtService := jwt.NewService(cfg.Jwt.Secret, cfg.Jwt.AccessExpiration)
	usersModule := users.NewModule(db)
	authModule := auth.NewModule(usersModule.Service, jwtService)

	err = createAdmin(cfg.Admin, authModule.Service)
	failOnError(err, "create initial admin")

	e := echo.New()
	e.HTTPErrorHandler = echox.ErrorHandler

	e.Use(middleware.Recover())
	api := e.Group("/api")

	api.Use(jwt.NewAuthMiddleware(cfg.Jwt.Secret))
	api.Use(echox.Logger)

	api.POST("/auth/register", authModule.Handler.Register)
	api.POST("/users/login", authModule.Handler.Login)

	api.DELETE("/users/:userId", usersModule.Handler.Delete, auth.Self)
	api.PUT("/users/:userId", usersModule.Handler.Update, auth.Self)
	api.GET("/users/:userId", usersModule.Handler.Get)
	api.PUT("/users/:userId/role/:role", usersModule.Handler.SetRole, auth.Admin)

	go func() {
		signalChanel := make(chan os.Signal, 1)
		signal.Notify(signalChanel, os.Interrupt)
		<-signalChanel
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			slog.Error("server shutdown", "error", err)
		}
	}()

	slog.Info("server started", "port", cfg.Port)
	err = e.Start(fmt.Sprintf(":%d", cfg.Port))
	if err != nil && err != http.ErrServerClosed {
		slog.Error("server stopped", "error", err)
	}
	slog.Info("server stopped")
}

func createAdmin(cfg config.AdminConfig, authService *auth.Service) error {
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
		return fmt.Errorf("register admin :%w", err)
	case err != nil:
		return nil
	default:
		slog.Info("admin user created", "username", cfg.Username, "email", cfg.Email)
		return nil

	}
}

func getDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, dbConnectTimeout)
	defer cancel()

	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}
	return db, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		slog.Error(msg, "error", err)
		os.Exit(1)
	}
}
