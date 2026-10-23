package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	config "github.com/kostas-111/share_trip/configs"
	"github.com/kostas-111/share_trip/internal/api"
	"github.com/kostas-111/share_trip/internal/repository"
	"github.com/kostas-111/share_trip/internal/service"
)

func main() {
	ctx := context.Background()

	cfg := repository.Config{
		Host:     config.Env("DB_HOST", "localhost"),
		Port:     config.EnvInt("DB_PORT", 6544),
		User:     config.Env("DB_USER", "postgres"),
		Password: config.Env("DB_PASSWORD", "password"),
		DBName:   config.Env("DB_NAME", "share_trip"),
		SSLMode:  config.Env("DB_SSLMODE", "disable"),
	}

	pool, err := repository.NewPool(ctx, cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := repository.NewRepoPg(pool)
	svc := service.NewTripService(repo)
	tripHandler := api.NewTripHandler(svc)
	server := api.NewServer(repo, tripHandler)

	app := fiber.New()
	server.Route(app.Group("/api"))

	err = app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
