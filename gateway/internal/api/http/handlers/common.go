package handlers

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"go.uber.org/zap"
)

func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}

func HandleGrpcError(c *fiber.Ctx, log *zap.Logger, err error, defaultMsg string) error {
	st, ok := status.FromError(err)
	if !ok {
		return ErrorResponse(c, fiber.StatusInternalServerError, defaultMsg)
	}

	log.Error("gRPC Error",
		zap.String("code", st.Code().String()),
		zap.String("msg", st.Message()),
		zap.String("request_id", c.Get("X-Request-ID")))

	switch st.Code() {
	case codes.NotFound:
		return ErrorResponse(c, fiber.StatusNotFound, st.Message())
	case codes.InvalidArgument:
		return ErrorResponse(c, fiber.StatusBadRequest, st.Message())
	case codes.AlreadyExists:
		return ErrorResponse(c, fiber.StatusConflict, st.Message())
	case codes.Unauthenticated:
		return ErrorResponse(c, fiber.StatusUnauthorized, st.Message())
	case codes.Unavailable:
		return ErrorResponse(c, fiber.StatusServiceUnavailable, "service temporarily unavailable")
	default:
		return ErrorResponse(c, fiber.StatusInternalServerError, st.Message())
	}
}
