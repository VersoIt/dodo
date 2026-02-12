package handlers

import (
	"context"
	"time"
	"github.com/gofiber/fiber/v2"
	auth_pb "github.com/versoit/diploma/services/auth/api/proto/pb"
	"go.uber.org/zap"
)

type AuthHandler struct {
	log    *zap.Logger
	client auth_pb.UserServiceClient
}

func NewAuthHandler(log *zap.Logger, client auth_pb.UserServiceClient) *AuthHandler {
	return &AuthHandler{
		log:    log,
		client: client,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request")
	}

	h.log.Info("Registering new user via gRPC", zap.String("email", req.Email))
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.CreateUser(ctx, &auth_pb.CreateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "registration failed")
	}

	return SuccessResponse(c, resp)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request")
	}

	h.log.Info("User login attempt", zap.String("email", req.Email))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.Login(ctx, &auth_pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "login failed")
	}

	return SuccessResponse(c, fiber.Map{
		"token":   resp.Token,
		"user_id": resp.UserId,
	})
}