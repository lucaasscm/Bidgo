package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lucaasscm/bidgo/internal/api"
	"github.com/lucaasscm/bidgo/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("BIDGO_DATABASE_USER"),
		os.Getenv("BIDGO_DATABASE_PASSWORD"),
		os.Getenv("BIDGO_DATABASE_HOST"),
		os.Getenv("BIDGO_DATABASE_PORT"),
		os.Getenv("BIDGO_DATABASE_NAME"),
	))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	api := api.Api{
		Router:      chi.NewMux(),
		UserService: services.NewUserService(pool),
	}

	api.BindRoutes()

	fmt.Println("Server running on port :3080")
	if err := http.ListenAndServe(":3080", api.Router); err != nil {
		panic(err)
	}
}
