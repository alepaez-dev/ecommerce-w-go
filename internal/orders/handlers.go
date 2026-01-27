package orders

import (
	"errors"
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

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var tempOrder createOrderParams

	if err := json.Read(r, &tempOrder); err != nil {
		log.Printf("error reading PlaceOrder payload %s", err)
		http.Error(w, "Invalid Payload", http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)

	if err != nil {
		log.Printf("error reading creating order %s", err)

		if errors.Is(err, ErrRequiredValue) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if errors.Is(err, ErrProductNoStock) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, "There was an error creating the order", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, createdOrder)
}
