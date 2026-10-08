package cores

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type BaseResponse struct {
	Success bool   `json:"success" msgpack:"success"`
	Message string `json:"message" msgpack:"message"`
	Data    any    `json:"data,omitempty" msgpack:"data,omitempty"`
	Error   any    `json:"error,omitempty" msgpack:"error,omitempty"`
}

func sendResponse(c fiber.Ctx, status int, payload BaseResponse) error {
	c.Vary(fiber.HeaderAccept)

	match := c.Accepts("application/json", "application/x-msgpack", "application/msgpack", "application/vnd.msgpack")
	if match != "" && match != "application/json" {
		if err := c.Status(status).MsgPack(payload, match); err != nil {
			zap.L().Error("failed to encode msgpack response", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}
		return nil
	}

	return c.Status(status).JSON(payload)
}

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

func RespForbidden(c fiber.Ctx, message string) error {
	return sendResponse(c, fiber.StatusForbidden, BaseResponse{
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
