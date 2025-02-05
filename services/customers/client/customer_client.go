package client

import (
	"context"

	pb "github.com/bsanzhiev/bahamas/libs/pb/customers"
	"google.golang.org/grpc"
)

type CustomerClient struct {
	grpcClient pb.CustomerServiceClient
}

func NewCustomerService(conn *grpc.ClientConn) *CustomerClient {
	return &CustomerClient{
		grpcClient: pb.NewCustomerServiceClient(conn),
	}
}

func (c *CustomerClient) GetCustomerByID(ctx context.Context, id string) (*pb.Customer, error) {
	return c.grpcClient.GetCustomer(ctx, &pb.GetCustomerRequest{Id: id})
}
