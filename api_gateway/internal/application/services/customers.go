package gateway_services

import (
	"context"

	"github.com/bsanzhiev/bahamas/services/customers/client"
	"google.golang.org/grpc"
)

type CustomerService struct {
	customerClient *client.CustomerClient
}

func NewCustomerService(conn *grpc.ClientConn) *CustomerService {
	return &CustomerService{
		customerClient: client.NewCustomerService(conn),
	}
}

func (s *CustomerService) GetCustomer(id string) (interface{}, error) {
	return s.customerClient.GetCustomer(context.Background(), id)
}
