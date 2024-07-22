package application

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"orders-api/hanlder"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Route("/orders", loadRouters)
	return router
}

func loadRouters(router chi.Router) {
	orderHandler := &hanlder.Order{}

	router.Post("/", orderHandler.CreateOrder)
	router.Get("/", orderHandler.List)
	router.Get("/{id}", orderHandler.GetByID)
	router.Put("/{id}", orderHandler.UpdateByID)
	router.Delete("/{id}", orderHandler.UpdateByID)

}
