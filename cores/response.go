package cores

import (
	"github.com/gofiber/fiber/v3"
	"github.com/shamaton/msgpack/v3"
	"go.uber.org/zap"
)

type BaseResponse struct {
	Success bool   `json:"success" msgpack:"success"`
	Message string `json:"message" msgpack:"message"`
	Data    any    `json:"data,omitempty" msgpack:"data,omitempty"`
	Error   any    `json:"error,omitempty" msgpack:"error,omitempty"`
}

func sendResponse(c fiber.Ctx, status int, payload BaseResponse) error {
	if c.Get("Accept") == "application/x-msgpack" {
		b, err := msgpack.Marshal(payload)
		if err != nil {
			zap.L().Error("failed to marshal msgpack", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}
		c.Set("Content-Type", "application/x-msgpack")
		return c.Status(status).Send(b)
	}

	// Default  JSON
	return c.Status(status).JSON(payload)
}

// --- List Helper Functions ---

func RespSuccess(c fiber.Ctx, message string, data any) error {
	return sendResponse(c, fiber.StatusOK, BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func RespCreated(c fiber.Ctx, message string, data any) error {
	return sendResponse(c, fiber.StatusCreated, BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func RespBadReq(c fiber.Ctx, message string, err any) error {
	return sendResponse(c, fiber.StatusBadRequest, BaseResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func RespUnauthorized(c fiber.Ctx, message string) error {
	return sendResponse(c, fiber.StatusUnauthorized, BaseResponse{
		Success: false,
		Message: message,
	})
}

func RespNotFound(c fiber.Ctx, message string) error {
	return sendResponse(c, fiber.StatusNotFound, BaseResponse{
		Success: false,
		Message: message,
	})
}

func RespInternalError(c fiber.Ctx, message string, err error) error {
	zap.L().Error(message, zap.Error(err))
	return sendResponse(c, fiber.StatusInternalServerError, BaseResponse{
		Success: false,
		Message: "Internal Server Error",
	})
}
