package products

import (
	"log"
	"net/http"
	"strconv"

	"github.com/alepaez-dev/ecommerce/internal/json"
	"github.com/go-chi/chi"
)

type handler struct {
	service Service
}

func NewHandler(svc Service) *handler {
	return &handler{service: svc}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Printf("error fetching products: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (h *handler) FindProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("error parsing id param: %s", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	product, err := h.service.FindProduct(r.Context(), id)
	if err != nil {
		log.Printf("error fetching products: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, product)
}
