package handler

import (
	"context"
	"log/slog"

	"github.com/versoit/diploma/pkg/common"
	treasury_pb "github.com/versoit/diploma/be/treasury/api/proto/pb"
	"github.com/versoit/diploma/be/treasury/internal/domain"
	"github.com/versoit/diploma/be/treasury/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	payment, err := h.uc.InitiatePayment(ctx, req.OrderId, common.NewMoney(req.Amount), domain.PaymentMethod(req.Method))
	if err != nil {
		h.log.Error("Failed to initiate payment", slog.String("order_id", req.OrderId), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "initiate payment: %v", err)
	}

	// For demonstration purposes, auto-confirm payment
	if err := h.uc.ConfirmPayment(ctx, req.OrderId, "demo-tx-id-"+payment.ID()); err != nil {
		h.log.Error("Failed to confirm payment", slog.String("order_id", req.OrderId), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "confirm payment: %v", err)
	}

	return &treasury_pb.PaymentResponse{
		PaymentId: payment.ID(),
		Status:    "success",
	}, nil
}
