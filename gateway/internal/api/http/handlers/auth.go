package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	auth_pb "github.com/versoit/diploma/services/auth/api/proto/pb"
	"log/slog"
)

type AuthHandler struct {
	log    *slog.Logger
	client auth_pb.UserServiceClient
}

func NewAuthHandler(log *slog.Logger, client auth_pb.UserServiceClient) *AuthHandler {
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

	h.log.Info("Registering new user via gRPC", "email", req.Email)

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

	h.log.Info("User login attempt", "email", req.Email)

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

func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	if userID == "" {
		return ErrorResponse(c, fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetUser(ctx, &auth_pb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to get user info")
	}

	return SuccessResponse(c, resp)
}
