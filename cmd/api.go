package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alepaez-dev/ecommerce/internal/adapters/postgresql"
	repo "github.com/alepaez-dev/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/alepaez-dev/ecommerce/internal/orders"
	"github.com/alepaez-dev/ecommerce/internal/products"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // important for rate limiting and analytics.
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	productService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productService)

	// Products
	r.Get("/products", productHandler.ListProducts)
	r.Get("/products/{id}", productHandler.FindProduct)

	// Orders
	txManager := postgresql.NewTxManager(app.db)

	productFactory := func(q repo.Querier) orders.ProductStore {
		return products.NewService(q)
	}
	orderService := orders.NewService(txManager, productFactory)
	orderHandler := orders.NewHandler(orderService)
	r.Post("/orders", orderHandler.PlaceOrder)

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Println("starting server at address " + app.config.addr)

	return srv.ListenAndServe()

}

type application struct {
	config config
	db     *pgx.Conn
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
