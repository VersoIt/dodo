package grpc

import (
	"context"
	"time"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/versoit/diploma/services/auth"
	auth_pb "github.com/versoit/diploma/services/auth/api/proto/pb"
	"github.com/versoit/diploma/services/auth/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var jwtKey = []byte("super-secret-key")

type AuthHandler struct {
	auth_pb.UnimplementedUserServiceServer
	uc  *usecase.AuthUseCase
	log *slog.Logger
}

func NewAuthHandler(uc *usecase.AuthUseCase, log *slog.Logger) *AuthHandler {
	return &AuthHandler{
		uc:  uc,
		log: log,
	}
}

func (h *AuthHandler) Register(server *grpc.Server) {
	auth_pb.RegisterUserServiceServer(server, h)
}

func (h *AuthHandler) CreateUser(ctx context.Context, req *auth_pb.CreateUserRequest) (*auth_pb.UserResponse, error) {
	h.log.Info("Creating user", slog.String("email", req.Email))
	
	user, err := h.uc.Register(ctx, req.Email, req.Password, req.Name, auth.RoleClient)
	if err != nil {
		h.log.Error("Failed to register user", slog.String("email", req.Email), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &auth_pb.UserResponse{
		Id:    user.ID(),
		Email: user.Email(),
		Role:  user.Role().String(),
		Name:  user.Name(),
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *auth_pb.LoginRequest) (*auth_pb.LoginResponse, error) {
	h.log.Info("User login", slog.String("email", req.Email))
	
	user, err := h.uc.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.log.Warn("Login failed", slog.String("email", req.Email), slog.Any("error", err))
		return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID(),
		"role":    user.Role().String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		h.log.Error("Token generation failed", slog.String("user_id", user.ID()), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &auth_pb.LoginResponse{
		Token:  tokenString,
		UserId: user.ID(),
	}, nil
}

func (h *AuthHandler) GetUser(ctx context.Context, req *auth_pb.GetUserRequest) (*auth_pb.UserResponse, error) {
	user, err := h.uc.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &auth_pb.UserResponse{
		Id:    user.ID(),
		Email: user.Email(),
		Role:  user.Role().String(),
		Name:  user.Name(),
	}, nil
}
