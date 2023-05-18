package main

import (
	"context"
	"fmt"
	"github.com/AnatoliyRib1/movie-reviews/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const dbConnectTimeout = 10 * time.Second

func main() {
	e := echo.New()
	cfg, err := config.NewConfig()
	failOnError(err, "parse config")

	fmt.Printf("started with config: %w", cfg)

	db, err := getDB(context.Background(), cfg.DbUrl)
	failOnError(err, "connect to db")

	err = db.Ping(context.Background())
	failOnError(err, "ping db")

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, os.Interrupt)

	go func() {
		<-signalChanel
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			log.Printf("shutdown: %w", err)
		}
	}()

	err = e.Start(fmt.Sprintf(":%d", &cfg.Port))
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
