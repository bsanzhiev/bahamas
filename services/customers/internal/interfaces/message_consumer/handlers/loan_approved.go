package kafka_handlers

import (
	"context"
	"customers_service/internal/application/events"
	"log"
)

func HandleLoanApproved(ctx context.Context, msg []byte) {
	event, err := events.DeserializeLoanApprovedEvent(msg)
	if err != nil {
		log.Printf("Failed to deserialize LoanApprovedEvent: %v", err)
		return
	}

	// Logic
	log.Printf("handling LoanApproved event: %+v", event)
	// Например, обновить данные клиента в базе данных
	// updateCustomerData(event.CustomerID, event.LoanID, event.Amount) и т.д.
}
