package cores_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/rachmanzz/fiber-starter/cores"
)

func TestAppContracts_RegisterRoute_MultipleRoutes(t *testing.T) {
	routeA := func(app *fiber.App) {
		app.Get("/route-a", func(c fiber.Ctx) error {
			return c.SendString("route-a-response")
		})
	}

	routeB := func(app *fiber.App) {
		app.Get("/route-b", func(c fiber.Ctx) error {
			return c.SendString("route-b-response")
		})
	}

	core := cores.CreateContract().
		CreateApp().
		RegisterRoute(routeA, routeB)

	// Test Route A
	respA, err := core.App.Test(httptest.NewRequest("GET", "/route-a", nil))
	if err != nil {
		t.Fatalf("unexpected error for route-a: %v", err)
	}
	if respA.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, respA.StatusCode)
	}

	// Test Route B
	respB, err := core.App.Test(httptest.NewRequest("GET", "/route-b", nil))
	if err != nil {
		t.Fatalf("unexpected error for route-b: %v", err)
	}
	if respB.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, respB.StatusCode)
	}
}

func TestAppContracts_RegisterHook_And_Middleware(t *testing.T) {
	hookExecuted := false
	middlewareExecuted := false

	hookFn := func(core *cores.AppContracts) {
		hookExecuted = true
	}

	middlewareFn := func(app *fiber.App) {
		app.Use(func(c fiber.Ctx) error {
			middlewareExecuted = true
			return c.Next()
		})
	}

	routeFn := func(app *fiber.App) {
		app.Get("/test", func(c fiber.Ctx) error {
			return c.SendString("ok")
		})
	}

	core := cores.CreateContract().
		CreateApp().
		RegisterHook(hookFn).
		RegisterMiddleware(middlewareFn).
		RegisterRoute(routeFn)

	if !hookExecuted {
		t.Fatalf("expected hook to be executed")
	}

	resp, err := core.App.Test(httptest.NewRequest("GET", "/test", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
	if !middlewareExecuted {
		t.Fatalf("expected middleware to be executed")
	}
}
