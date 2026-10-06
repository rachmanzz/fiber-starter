package bootstrap_test

import (
	"net/http/httptest"
	"testing"


	"github.com/gofiber/fiber/v3"
	"github.com/rachmanzz/fiber-starter/app/routes"
	"github.com/rachmanzz/fiber-starter/bootstrap"
	"github.com/rachmanzz/fiber-starter/cores"
)

// newTestApp builds an app following the same flow as Application.Bootstrap:
// CreateApp -> RegisterMiddleware -> routes.
func newTestApp(t *testing.T) *fiber.App {
	t.Helper()
	core := cores.CreateContract().
		CreateApp().
		RegisterMiddleware(bootstrap.RegisterMiddleware).
		RegisterRoute(routes.ApiRoute)
	return core.App
}

func TestRecoverMiddleware_PanicReturns500(t *testing.T) {
	core := cores.CreateContract().CreateApp()
	bootstrap.RegisterMiddleware(core.App)
	core.App.Get("/boom", func(c fiber.Ctx) error {
		panic("boom")
	})

	resp, err := core.App.Test(httptest.NewRequest("GET", "/boom", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("expected status %d after panic, got %d", fiber.StatusInternalServerError, resp.StatusCode)
	}
}

func TestRecoverMiddleware_ServerSurvivesPanic(t *testing.T) {
	core := cores.CreateContract().CreateApp()
	bootstrap.RegisterMiddleware(core.App)
	core.App.Get("/boom", func(c fiber.Ctx) error {
		panic("boom")
	})
	core.App.Get("/alive", func(c fiber.Ctx) error {
		return c.SendString("alive")
	})

	// Trigger the panic first.
	resp, err := core.App.Test(httptest.NewRequest("GET", "/boom", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("expected status %d after panic, got %d", fiber.StatusInternalServerError, resp.StatusCode)
	}

	// The server must still serve subsequent requests.
	resp, err = core.App.Test(httptest.NewRequest("GET", "/alive", nil))
	if err != nil {
		t.Fatalf("server did not survive the panic: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestApiRoute(t *testing.T) {
	app := newTestApp(t)

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestApiRouteUnknownPath(t *testing.T) {
	app := newTestApp(t)

	resp, err := app.Test(httptest.NewRequest("GET", "/does-not-exist", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}
