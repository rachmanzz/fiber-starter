package cores

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
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
	app.appMu.Lock()
	mappers := append([]ErrorMapperFn(nil), app.errorMappers...)
	app.appMu.Unlock()

	for _, mapper := range mappers {
		if status, payload, handled := mapper(c, err); handled {
			msg := "Error"
			switch m := payload.(type) {
			case fiber.Map:
				if customMsg, ok := m["message"].(string); ok && customMsg != "" {
					msg = customMsg
				}
			case map[string]any:
				if customMsg, ok := m["message"].(string); ok && customMsg != "" {
					msg = customMsg
				}
			}
			return sendResponse(c, status, BaseResponse{
				Success: false,
				Message: msg,
				Error:   payload,
			})
		}
	}

	// 2. Handle go-playground/validator errors
	var invalidValErr *validator.InvalidValidationError
	if errors.As(err, &invalidValErr) {
		return RespInternalError(c, "Invalid struct validation target", invalidValErr)
	}

	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		errMap := make(map[string]string)
		for _, e := range valErrs {
			errMap[e.Field()] = e.Tag()
		}
		return RespBadReq(c, "Validation failed", errMap)
	}

	// 3. Handle Fiber v3 Bind errors (e.g. invalid JSON or MsgPack payload)
	var bindErr *fiber.BindError
	if errors.As(err, &bindErr) {
		return RespBadReq(c, "Invalid request payload", "malformed request body or parameters")
	}

	// 4. Handle Fiber v3 native HTTP errors (e.g. 404, 400, 405)
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return sendResponse(c, fiberErr.Code, BaseResponse{
			Success: false,
			Message: fiberErr.Message,
		})
	}

	// 5. Default fallback: Internal Server Error
	return RespInternalError(c, "Unhandled server error", err)
}
