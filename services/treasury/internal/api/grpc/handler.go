package grpc

import (
	"context"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/treasury"
	treasury_pb "github.com/versoit/diploma/services/treasury/api/proto/pb"
	"github.com/versoit/diploma/services/treasury/usecase"
	"google.golang.org/grpc"
)

type TreasuryHandler struct {
	treasury_pb.UnimplementedPaymentServiceServer
	uc *usecase.TreasuryUseCase
}

func NewTreasuryHandler(uc *usecase.TreasuryUseCase) *TreasuryHandler {
	return &TreasuryHandler{uc: uc}
}

func (h *TreasuryHandler) Register(server *grpc.Server) {
	treasury_pb.RegisterPaymentServiceServer(server, h)
}

func (h *TreasuryHandler) ProcessPayment(ctx context.Context, req *treasury_pb.PaymentRequest) (*treasury_pb.PaymentResponse, error) {
	payment, err := h.uc.InitiatePayment(ctx, req.OrderId, common.NewMoney(req.Amount), treasury.PaymentMethod(req.Method))
	if err != nil {
		return nil, err
	}

	return &treasury_pb.PaymentResponse{
		PaymentId: payment.ID(),
		Status:    payment.Status().String(),
	}, nil
}