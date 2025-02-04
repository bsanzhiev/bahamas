package events

import (
	"encoding/json"
	"fmt"
	"time"
)

type LoanApprovedEvent struct {
	CustomerID string  `json:"customer_id"`
	LoanID     string  `json:"loan_id"`
	Amount     float64 `json:"amount"`
	Timestamp  int64   `json:"timestump"`
}

func NewLoanApprovedEvent(customerID, LoanID string, amount float64) (LoanApprovedEvent, error) {
	return LoanApprovedEvent{
		CustomerID: customerID,
		LoanID:     LoanID,
		Amount:     amount,
		Timestamp:  time.Now().Unix(),
	}, nil
}

// Serialize event to JSON for sending throgh Kafka
func (e LoanApprovedEvent) Serialize() ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize LoanApprovedEvent: %w", err)
	}
	return data, nil
}

// JSON to event
func DeserializeLoanApprovedEvent(data []byte) (LoanApprovedEvent, error) {
	var event LoanApprovedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return LoanApprovedEvent{}, fmt.Errorf("failed to deserialize msg: %w", err)
	}
	return event, nil
}
