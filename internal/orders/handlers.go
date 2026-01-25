package orders

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

		if err == ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		if err == ErrProductNoStock {
			http.Error(w, err.Error(), http.StatusConflict)
		}

		http.Error(w, "There was an error creating the order", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, createdOrder)
}
