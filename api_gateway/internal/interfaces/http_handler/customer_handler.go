package http_handler

// Вызывает Use Cases

import (
	"log"
	"net/http"

	gateway_services "github.com/bsanzhiev/bahamas/api_gateway/internal/application/services"
)

type CustomerHandler struct {
	customerService gateway_services.CustomerService
	logger          *log.Logger
}

func NewCustomerHandler(
	customerService gateway_services.CustomerService,
	logger *log.Logger,
) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
		logger:          logger,
	}
}

func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	// Get customer ID from query
	customerID := r.URL.Query().Get("id")

	// Call Use Case from customerService
	customer, err := h.customerService.GetCustomerByID(customerID)
	if err != nil {
		h.logger.Printf("Failed to get customer: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Temp
	type customer_contain struct {
		ID   string
		Name string
	}

	data, ok := customer.(customer_contain)
	if !ok {
		h.logger.Printf("Failed to convert customer to map")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data.Name))
}
