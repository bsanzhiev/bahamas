package gatewayservices

import "github.com/bsanzhiev/bahamas/services/customers/client"

type CustomerService struct {
	customerClient *client.CustomerClient
}
