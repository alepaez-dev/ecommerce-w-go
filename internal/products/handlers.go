package products

import (
	"log"
	"net/http"

	"github.com/alepaez-dev/ecommerce/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(svc Service) *handler {
	return &handler{service: svc}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Printf("error fetching products: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products := struct {
		Products []string `json:"products"`
	}{}

	json.Write(w, http.StatusOK, products)
}
