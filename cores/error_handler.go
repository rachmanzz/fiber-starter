package cores

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// ErrorMapperFn allows users to map custom error types to HTTP response status codes and payloads.
// It returns (statusCode, payload, handled). If handled is false, fallback processing continues.
type ErrorMapperFn func(c fiber.Ctx, err error) (int, any, bool)

// GlobalErrorHandler handles error responses centrally and supports extensible custom error mappers.
func (app *AppContracts) GlobalErrorHandler(c fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// 1. Check custom error mappers registered by user (e.g. from app/errors)
	for _, mapper := range app.errorMappers {
		if status, payload, handled := mapper(c, err); handled {
			return sendResponse(c, status, BaseResponse{
				Success: false,
				Message: "Error",
				Error:   payload,
			})
		}
	}

	// 2. Handle go-playground/validator errors
	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		errMap := make(map[string]string)
		for _, e := range valErrs {
			errMap[e.Field()] = e.Tag()
		}
		return RespBadReq(c, "Validation failed", errMap)
	}

	// 3. Handle Fiber v3 native HTTP errors (e.g. 404, 400, 405)
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return sendResponse(c, fiberErr.Code, BaseResponse{
			Success: false,
			Message: fiberErr.Message,
		})
	}

	// 4. Default fallback: Internal Server Error
	zap.L().Error("Unhandled server error", zap.Error(err))
	return RespInternalError(c, "Internal Server Error", err)
}
