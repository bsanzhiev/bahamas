package client

import pb "github.com/bsanzhiev/bagamas/libs/pb/customers"

type CustomerClient struct {
	// proto is packet from customers/proto
	grpcClient pb.CustomerServiceClient
}
