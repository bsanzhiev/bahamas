package client

type CustomerClient struct {
	// proto is packet from customers/proto
	grpcClient pb.CustomerServiceClient
}
