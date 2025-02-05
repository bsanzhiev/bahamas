package gateway_services

import (
	"context"

	"github.com/bsanzhiev/bahamas/services/customers/client"
	"google.golang.org/grpc"
)

type CustomerService struct {
	customerClient *client.CustomerClient
}

// type Customer struct {
// 	ID   string
// 	Name string
// }

func NewCustomerService(conn *grpc.ClientConn) *CustomerService {
	return &CustomerService{
		customerClient: client.NewCustomerService(conn),
	}
}

func (s *CustomerService) GetCustomerByID(id string) (interface{}, error) {
	return s.customerClient.GetCustomerByID(context.Background(), id)
}
