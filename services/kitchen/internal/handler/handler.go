package handler

import (
	"context"
	"log/slog"
	"time"

	kitchen_pb "github.com/versoit/diploma/services/kitchen/api/proto/pb"
	"github.com/versoit/diploma/services/kitchen/internal/domain"
	"github.com/versoit/diploma/services/kitchen/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	h.log.Info("Creating kitchen ticket", slog.String("order_id", req.OrderId), slog.String("order_number", req.OrderNumber))

	items := make([]domain.KitchenItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.KitchenItem{
			ProductID: item.ProductId,
			Name:      item.ProductName,
			Quantity:  int(item.Quantity),
		}
	}

	ticket, err := h.uc.AcceptOrder(ctx, req.OrderId, req.OrderNumber, items)
	if err != nil {
		h.log.Error("Failed to create ticket", slog.String("order_id", req.OrderId), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "accept order: %v", err)
	}

	return &kitchen_pb.TicketResponse{
		TicketId: ticket.ID(),
		Status:   ticket.Status().String(),
	}, nil
}

func (h *KitchenHandler) ListTickets(ctx context.Context, req *kitchen_pb.ListTicketsRequest) (*kitchen_pb.ListTicketsResponse, error) {
	tickets, err := h.uc.ListTickets(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list tickets: %v", err)
	}

	pbTickets := make([]*kitchen_pb.Ticket, len(tickets))
	for i, t := range tickets {
		items := make([]*kitchen_pb.KitchenItem, len(t.Items()))
		for j, it := range t.Items() {
			items[j] = &kitchen_pb.KitchenItem{
				ProductId:   it.ProductID,
				ProductName: it.Name,
				Quantity:    int32(it.Quantity),
			}
		}
		pbTickets[i] = &kitchen_pb.Ticket{
			Id:          t.ID(),
			OrderId:     t.OrderID(),
			OrderNumber: t.OrderNumber(),
			Status:      t.Status().String(),
			Items:       items,
			CreatedAt:   t.CreatedAt().Format(time.RFC3339),
		}
	}

	return &kitchen_pb.ListTicketsResponse{Tickets: pbTickets}, nil
}

func (h *KitchenHandler) UpdateTicketStatus(ctx context.Context, req *kitchen_pb.UpdateTicketStatusRequest) (*kitchen_pb.TicketResponse, error) {
	h.log.Info("Updating ticket status", slog.String("ticket_id", req.TicketId), slog.String("status", req.Status))

	var err error
	var orderID string
	switch req.Status {
	case "cooking":
		orderID, err = h.uc.StartCooking(ctx, req.TicketId)
	case "ready":
		orderID, err = h.uc.MarkReady(ctx, req.TicketId)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid status: %s", req.Status)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "update status: %v", err)
	}

	return &kitchen_pb.TicketResponse{
		TicketId: req.TicketId,
		OrderId:  orderID,
		Status:   req.Status,
	}, nil
}
