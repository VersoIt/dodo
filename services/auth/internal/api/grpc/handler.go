package grpc

import (
	"context"

	"github.com/versoit/diploma/services/auth"
	"github.com/versoit/diploma/services/auth/api/proto/pb"
	"github.com/versoit/diploma/services/auth/usecase"
	"google.golang.org/grpc"
)

type AuthHandler struct {
	pb.UnimplementedUserServiceServer
	uc *usecase.AuthUseCase
}

func NewAuthHandler(uc *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(server *grpc.Server) {
	pb.RegisterUserServiceServer(server, h)
}

func (h *AuthHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := h.uc.Register(ctx, req.Email, req.Password, auth.RoleClient)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:    user.ID(),
		Email: user.Email(),
		Role:  user.Role().String(),
	}, nil
}

func (h *AuthHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	// Logic to get user
	return &pb.UserResponse{Id: req.Id}, nil
}
