package cores_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/rachmanzz/fiber-starter/cores"
)

type CustomDomainError struct {
	Msg string
}

func (e *CustomDomainError) Error() string {
	return e.Msg
}

func TestGlobalErrorHandler(t *testing.T) {
	appContract := cores.CreateContract()

	// Register custom error mapper (representing user mapping from app/errors)
	appContract.RegisterErrorMapper(func(c fiber.Ctx, err error) (int, any, bool) {
		var customErr *CustomDomainError
		if errors.As(err, &customErr) {
			return fiber.StatusTeapot, fiber.Map{
				"code":    "CUSTOM_TEAPOT",
				"detail":  customErr.Error(),
			}, true
		}
		return 0, nil, false
	})

	appContract.CreateApp()
	app := appContract.App

	// Route that returns a custom domain error
	app.Get("/test-custom-err", func(c fiber.Ctx) error {
		return &CustomDomainError{Msg: "custom domain error occurred"}
	})

	// Route that returns standard Fiber 404/400 error
	app.Get("/test-fiber-err", func(c fiber.Ctx) error {
		return fiber.ErrBadRequest
	})

	t.Run("Custom User Error Mapper", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test-custom-err", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to perform request: %v", err)
		}
		if resp.StatusCode != fiber.StatusTeapot {
			t.Errorf("Expected status %d, got %d", fiber.StatusTeapot, resp.StatusCode)
		}
	})

	t.Run("Fiber Native Error Fallback", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test-fiber-err", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to perform request: %v", err)
		}
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}
