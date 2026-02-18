package handler

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/google/uuid"
	pb "github.com/versoit/diploma/be/chat/api/proto/pb"
	"github.com/versoit/diploma/be/chat/internal/config"
	"github.com/versoit/diploma/be/chat/internal/domain"
	"github.com/versoit/diploma/be/chat/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatHandler struct {
	pb.UnimplementedChatServiceServer
	svc service.MessageService
	hub *service.Hub
	cfg *config.Config
	log *slog.Logger
}

func NewChatHandler(svc service.MessageService, hub *service.Hub, cfg *config.Config, log *slog.Logger) *ChatHandler {
	return &ChatHandler{
		svc: svc,
		hub: hub,
		cfg: cfg,
		log: log,
	}
}

func (h *ChatHandler) Register(server *grpc.Server) {
	pb.RegisterChatServiceServer(server, h)
}

func (h *ChatHandler) GetHistory(ctx context.Context, req *pb.GetHistoryRequest) (*pb.GetHistoryResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id")
	}

	messages, err := h.svc.GetHistory(ctx, orderID, int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get history: %v", err)
	}

	pbMessages := make([]*pb.Message, len(messages))
	for i, m := range messages {
		pbMessages[i] = &pb.Message{
			Id:        m.ID,
			OrderId:   m.OrderID.String(),
			SenderId:  m.SenderID.String(),
			Role:      string(m.Role),
			Text:      m.Text,
			IsRead:    m.IsRead,
			CreatedAt: timestamppb.New(m.CreatedAt),
		}
	}

	return &pb.GetHistoryResponse{Messages: pbMessages}, nil
}

func (h *ChatHandler) Connect(stream pb.ChatService_ConnectServer) error {
	// 1. Read the first message (Connect handshake)
	firstMsg, err := stream.Recv()
	if err != nil {
		return err
	}

	connectEvent := firstMsg.GetConnect()
	if connectEvent == nil {
		return status.Errorf(codes.InvalidArgument, "expected ConnectEvent as first message")
	}

	orderID, err := uuid.Parse(connectEvent.OrderId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid order_id")
	}
	userID, err := uuid.Parse(connectEvent.UserId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid user_id")
	}
	role := domain.Role(connectEvent.Role)

	h.log.Info("Client connected", slog.String("user_id", userID.String()), slog.String("order_id", orderID.String()))

	// 2. Register Client
	client := &service.Client{
		Hub:     h.hub,
		OrderID: orderID,
		Send:    make(chan []byte, 256),
	}
	h.hub.Register <- client

	// 3. Cleanup on exit
	defer func() {
		h.hub.Unregister <- client
	}()

	// 4. Goroutine to send messages from Hub -> Stream
	errChan := make(chan error, 1)
	go func() {
		for payload := range client.Send {
			var wsMsg domain.WSMessage
			if err := json.Unmarshal(payload, &wsMsg); err != nil {
				h.log.Error("Failed to unmarshal broadcast message", slog.Any("error", err))
				continue
			}

			// Convert domain.WSMessage to pb.ServerMessage
			var serverMsg *pb.ServerMessage
			switch wsMsg.Event {
			case domain.EventNewMessage:
				serverMsg = &pb.ServerMessage{
					Event: &pb.ServerMessage_NewMessage{
						NewMessage: &pb.NewMessageEvent{
							MessageId: wsMsg.MessageID,
							Text:      wsMsg.Text,
							Role:      string(wsMsg.Role),
							OrderId:   wsMsg.OrderID.String(),
							CreatedAt: timestamppb.New(wsMsg.CreatedAt),
						},
					},
				}
			case domain.EventMessageAck:
				serverMsg = &pb.ServerMessage{
					Event: &pb.ServerMessage_Ack{
						Ack: &pb.MessageAck{
							RequestId: wsMsg.RequestID,
							MessageId: wsMsg.MessageID,
						},
					},
				}
			case domain.EventError:
				strPayload, _ := wsMsg.Payload.(string)
				serverMsg = &pb.ServerMessage{
					Event: &pb.ServerMessage_Error{
						Error: &pb.ErrorEvent{
							Message: strPayload,
						},
					},
				}
			}

			if serverMsg != nil {
				if err := stream.Send(serverMsg); err != nil {
					h.log.Error("Stream send error", slog.Any("error", err))
					errChan <- err
					return
				}
			}
		}
	}()

	// 5. Loop read from stream
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			select {
			case sendErr := <-errChan:
				return sendErr
			default:
				return err
			}
		}

		switch event := in.Event.(type) {
		case *pb.ClientMessage_SendMessage:
			msg := &domain.Message{
				OrderID:  orderID,
				SenderID: userID,
				Role:     role,
				Text:     event.SendMessage.Text,
			}
			
			ctx, cancel := context.WithTimeout(stream.Context(), 5*time.Second)
			err := h.svc.Save(ctx, msg)
			cancel()

			if err != nil {
				h.log.Error("Failed to save message", slog.Any("error", err))
				_ = stream.Send(&pb.ServerMessage{
					Event: &pb.ServerMessage_Error{
						Error: &pb.ErrorEvent{Message: "Failed to save message"},
					},
				})
				continue
			}

			_ = stream.Send(&pb.ServerMessage{
				Event: &pb.ServerMessage_Ack{
					Ack: &pb.MessageAck{
						RequestId: event.SendMessage.RequestId,
						MessageId: msg.ID,
					},
				},
			})

			h.svc.Broadcast(domain.WSMessage{
				Event:     domain.EventNewMessage,
				MessageID: msg.ID,
				Text:      msg.Text,
				Role:      msg.Role,
				OrderID:   orderID,
				CreatedAt: msg.CreatedAt,
			})

		case *pb.ClientMessage_ReadMessage:
			if event.ReadMessage.MessageId != 0 {
				ctx, cancel := context.WithTimeout(stream.Context(), 2*time.Second)
				_ = h.svc.MarkAsRead(ctx, event.ReadMessage.MessageId)
				cancel()
			}
		}
	}
}
