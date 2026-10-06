package cores_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/rachmanzz/fiber-starter/cores"
)

type SampleRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

func TestStructValidator(t *testing.T) {
	appContract := cores.CreateContract()
	appContract.CreateApp()
	app := appContract.App

	app.Post("/test-bind", func(c fiber.Ctx) error {
		var req SampleRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"message": "success",
			"data":    req,
		})
	})

	t.Run("Valid Payload", func(t *testing.T) {
		payload := SampleRequest{Name: "John Doe", Email: "john@example.com"}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/test-bind", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to perform request: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid Payload Validation Error", func(t *testing.T) {
		payload := SampleRequest{Name: "Jo", Email: "invalid-email"}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/test-bind", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to perform request: %v", err)
		}
		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}
