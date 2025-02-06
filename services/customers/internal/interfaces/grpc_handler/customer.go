package grpc_handler

import (
	"context"

	pb "github.com/bsanzhiev/bahamas/libs/pb/customers"
	"github.com/bsanzhiev/bahamas/services/customers/internal/application/usecases"
)

type CustomerHandler struct {
	useCase usecases.CustomerUseCase
}

func NewCustomerHandler(useCase usecases.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{
		useCase: useCase,
	}
}

func (h *CustomerHandler) GetCustomer(
	ctx context.Context,
	req *pb.GetCustomerRequest,
) (
	*pb.GetCustomerResponse,
	error,
) {
	customer, err := h.useCase.GetCustomerByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.GetCustomerResponse{
		Customer: customer,
	}, nil
}
