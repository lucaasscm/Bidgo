package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lucaasscm/bidgo/internal/api"
	"github.com/lucaasscm/bidgo/internal/services"
)

func main() {
	gob.Register(uuid.UUID{})

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

	sessions := scs.New()
	sessions.Store = pgxstore.New(pool)
	sessions.Lifetime = 24 * time.Hour
	sessions.Cookie.HttpOnly = true
	sessions.Cookie.SameSite = http.SameSiteLaxMode

	csrfKey := os.Getenv("BIDGO_CSRF_KEY")
	if csrfKey == "" {
		panic("BIDGO_CSRF_KEY is not set")
	}

	api := api.Api{
		Router:         chi.NewMux(),
		Sessions:       sessions,
		UserService:    services.NewUserService(pool),
		ProductService: services.NewProductService(pool),
		BidsService:    services.NewBidsService(pool),
		CSRFKey:        []byte(csrfKey),
	}

	api.BindRoutes()

	fmt.Println("Server running on port :3080")
	if err := http.ListenAndServe(":3080", api.Router); err != nil {
		panic(err)
	}
}
