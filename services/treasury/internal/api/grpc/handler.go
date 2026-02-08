package grpc

import (
	"context"

	"github.com/versoit/diploma/services/treasury/api/proto/pb"
	"github.com/versoit/diploma/services/treasury/usecase"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/treasury"
	"google.golang.org/grpc"
)

type TreasuryHandler struct {
	pb.UnimplementedPaymentServiceServer
	uc *usecase.TreasuryUseCase
}

func NewTreasuryHandler(uc *usecase.TreasuryUseCase) *TreasuryHandler {
	return &TreasuryHandler{uc: uc}
}

func (h *TreasuryHandler) Register(server *grpc.Server) {
	pb.RegisterPaymentServiceServer(server, h)
}

func (h *TreasuryHandler) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	payment, err := h.uc.InitiatePayment(ctx, req.OrderId, common.Money(req.Amount), treasury.PaymentMethod(req.Method))
	if err != nil {
		return nil, err
	}

	return &pb.PaymentResponse{
		PaymentId: payment.ID(),
		Status:    payment.Status().String(),
	}, nil
}
