package usecases

import (
	"context"

	pb "github.com/bsanzhiev/bahamas/libs/pb/customers"
	shared "github.com/bsanzhiev/bahamas/libs/pb/shared"
)

// CustomerUseCase defines the use case interface for customer operations.
type CustomerUseCase interface {
	GetCustomerByID(ctx context.Context, req *pb.GetCustomerRequest) (*pb.Customer, error)
}

// CustomerUseCaseImpl is an example implementation of CustomerUseCase.
// Add any dependencies here, such as a repository or service client.
type customerUseCaseImpl struct {
}

// NewCustomerUseCase creates a new instance of customerUseCase.
func NewCustomerUseCase() CustomerUseCase {
	return &customerUseCaseImpl{}
}

// GetCustomerByID retrieves a customer by ID.
func (uc *customerUseCaseImpl) GetCustomerByID(ctx context.Context, req *pb.GetCustomerRequest) (*pb.Customer, error) {
	// Implement the logic to fetch the customer from a database or external service.
	// For now, we'll return a mock customer.

	// Example mock data
	mockCustomer := &pb.Customer{
		Id:   "123",
		Name: "John Doe",
		Address: &shared.Address{
			Street: "123 Main St",
			City:   "New York",
		},
	}

	return mockCustomer, nil
}
