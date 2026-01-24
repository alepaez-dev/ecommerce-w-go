package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alepaez-dev/ecommerce/internal/products"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	productService := products.NewService()
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.ListProducts)

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
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
