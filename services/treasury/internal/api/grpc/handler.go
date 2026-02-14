package grpc

import (
	"context"
	"log/slog"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/treasury"
	treasury_pb "github.com/versoit/diploma/services/treasury/api/proto/pb"
	"github.com/versoit/diploma/services/treasury/usecase"
	"google.golang.org/grpc"
)

type TreasuryHandler struct {
	treasury_pb.UnimplementedPaymentServiceServer
	uc  *usecase.TreasuryUseCase
	log *slog.Logger
}

func NewTreasuryHandler(uc *usecase.TreasuryUseCase, log *slog.Logger) *TreasuryHandler {
	return &TreasuryHandler{
		uc:  uc,
		log: log,
	}
}

func (h *TreasuryHandler) Register(server *grpc.Server) {
	treasury_pb.RegisterPaymentServiceServer(server, h)
}

func (h *TreasuryHandler) ProcessPayment(ctx context.Context, req *treasury_pb.PaymentRequest) (*treasury_pb.PaymentResponse, error) {
	h.log.Info("Processing payment", slog.String("order_id", req.OrderId), slog.Float64("amount", req.Amount))
	
	payment, err := h.uc.InitiatePayment(ctx, req.OrderId, common.NewMoney(req.Amount), treasury.PaymentMethod(req.Method))
	if err != nil {
		h.log.Error("Failed to initiate payment", slog.String("order_id", req.OrderId), slog.Any("error", err))
		return nil, err
	}

	return &treasury_pb.PaymentResponse{
		PaymentId: payment.ID(),
		Status:    payment.Status().String(),
	}, nil
}
