package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AnatoliyRib1/movie-reviews/internal/config"
	"github.com/AnatoliyRib1/movie-reviews/internal/jwt"
	"github.com/AnatoliyRib1/movie-reviews/internal/modules/auth"
	"github.com/AnatoliyRib1/movie-reviews/internal/modules/users"
	"github.com/AnatoliyRib1/movie-reviews/internal/validation"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

const dbConnectTimeout = 10 * time.Second

func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "parse config")

	validation.SetupValidators()

	fmt.Printf("started with config: %w", cfg)

	db, err := getDB(context.Background(), cfg.DbUrl)
	failOnError(err, "connect to db")

	jwtService := jwt.NewService(cfg.Jwt.Secret, cfg.Jwt.AccessExpiration)
	usersModule := users.NewModule(db)
	authModule := auth.NewModule(usersModule.Service, jwtService)

	e := echo.New()
	api := e.Group("/api")

	authMiddleware := jwt.NewAuthMiddleware(cfg.Jwt.Secret)

	api.POST("/auth/register", authModule.Handler.Register)
	api.POST("/users/login", authModule.Handler.Login)

	api.GET("/users/:userId", usersModule.Handler.GetUsers)
	api.DELETE("/users/:userId", usersModule.Handler.Delete, authMiddleware, auth.Self)
	api.PUT("/users/:userId", usersModule.Handler.Put, authMiddleware, auth.Self)

	go func() {
		signalChanel := make(chan os.Signal, 1)
		signal.Notify(signalChanel, os.Interrupt)
		<-signalChanel
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			log.Printf("shutdown: %w", err)
		}
	}()

	err = e.Start(fmt.Sprintf(":%d", cfg.Port))
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func getDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	ctx, _ = context.WithTimeout(ctx, dbConnectTimeout)

	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}
	return db, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", msg, err)
	}
}
