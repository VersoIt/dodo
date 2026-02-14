package grpc

import (
	"context"
	"log/slog"

	"github.com/versoit/diploma/services/kitchen"
	kitchen_pb "github.com/versoit/diploma/services/kitchen/api/proto/pb"
	"github.com/versoit/diploma/services/kitchen/usecase"
	"google.golang.org/grpc"
)

type KitchenHandler struct {
	kitchen_pb.UnimplementedTicketServiceServer
	uc  *usecase.KitchenUseCase
	log *slog.Logger
}

func NewKitchenHandler(uc *usecase.KitchenUseCase, log *slog.Logger) *KitchenHandler {
	return &KitchenHandler{
		uc:  uc,
		log: log,
	}
}

func (h *KitchenHandler) Register(server *grpc.Server) {
	kitchen_pb.RegisterTicketServiceServer(server, h)
}

func (h *KitchenHandler) CreateTicket(ctx context.Context, req *kitchen_pb.CreateTicketRequest) (*kitchen_pb.TicketResponse, error) {
	h.log.Info("Creating kitchen ticket", slog.String("order_id", req.OrderId))
	
	items := make([]kitchen.KitchenItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = kitchen.KitchenItem{
			ProductID: item.ProductId,
			Name:      item.ProductName,
			Quantity:  int(item.Quantity),
		}
	}

	ticket, err := h.uc.AcceptOrder(ctx, req.OrderId, items)
	if err != nil {
		h.log.Error("Failed to create ticket", slog.String("order_id", req.OrderId), slog.Any("error", err))
		return nil, err
	}

	return &kitchen_pb.TicketResponse{
		TicketId: ticket.ID(),
		Status:   ticket.Status().String(),
	}, nil
}

func (h *KitchenHandler) UpdateTicketStatus(ctx context.Context, req *kitchen_pb.UpdateTicketStatusRequest) (*kitchen_pb.TicketResponse, error) {
	h.log.Info("Updating ticket status", slog.String("ticket_id", req.TicketId), slog.Int("status", int(req.Status)))
	return &kitchen_pb.TicketResponse{TicketId: req.TicketId}, nil
}
