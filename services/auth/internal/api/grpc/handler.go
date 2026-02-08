package grpc

import (
	"context"

	"github.com/versoit/diploma/services/auth"
	auth_pb "github.com/versoit/diploma/services/auth/api/proto/pb"
	"github.com/versoit/diploma/services/auth/usecase"
	"google.golang.org/grpc"
)

type AuthHandler struct {
	auth_pb.UnimplementedUserServiceServer
	uc *usecase.AuthUseCase
}

func NewAuthHandler(uc *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(server *grpc.Server) {
	auth_pb.RegisterUserServiceServer(server, h)
}

func (h *AuthHandler) CreateUser(ctx context.Context, req *auth_pb.CreateUserRequest) (*auth_pb.UserResponse, error) {
	user, err := h.uc.Register(ctx, req.Email, req.Password, auth.RoleClient)
	if err != nil {
		return nil, err
	}

	return &auth_pb.UserResponse{
		Id:    user.ID(),
		Email: user.Email(),
		Role:  user.Role().String(),
	}, nil
}

func (h *AuthHandler) GetUser(ctx context.Context, req *auth_pb.GetUserRequest) (*auth_pb.UserResponse, error) {
	return &auth_pb.UserResponse{Id: req.Id}, nil
}
