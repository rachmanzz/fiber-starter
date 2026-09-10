package cores_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/rachmanzz/fiber-starter/cores"
	"github.com/shamaton/msgpack/v3"
)

func newTestApp(handler func(c fiber.Ctx) error) *fiber.App {
	app := fiber.New()
	app.Get("/", handler)
	return app
}

func doTest(t *testing.T, app *fiber.App, accept string) (*http.Response, error) {
	t.Helper()
	req := httptest.NewRequest("GET", "/", nil)
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	return app.Test(req)
}

func decodeJSON(t *testing.T, body io.Reader) cores.BaseResponse {
	t.Helper()
	var res cores.BaseResponse
	if err := json.NewDecoder(body).Decode(&res); err != nil {
		t.Fatalf("failed to decode json response: %v", err)
	}
	return res
}

func TestRespSuccess(t *testing.T) {
	app := newTestApp(func(c fiber.Ctx) error {
		return cores.RespSuccess(c, "fetched", fiber.Map{"id": 1})
	})

	resp, err := doTest(t, app, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	res := decodeJSON(t, resp.Body)
	if !res.Success || res.Message != "fetched" || res.Error != nil {
		t.Fatalf("unexpected payload: %+v", res)
	}
	data, ok := res.Data.(map[string]any)
	if !ok || data["id"] != float64(1) {
		t.Fatalf("unexpected data: %+v", res.Data)
	}
}

func TestRespCreated(t *testing.T) {
	app := newTestApp(func(c fiber.Ctx) error {
		return cores.RespCreated(c, "created", nil)
	})

	resp, err := doTest(t, app, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}

	res := decodeJSON(t, resp.Body)
	if !res.Success || res.Message != "created" {
		t.Fatalf("unexpected payload: %+v", res)
	}
}

func TestRespBadReq(t *testing.T) {
	app := newTestApp(func(c fiber.Ctx) error {
		return cores.RespBadReq(c, "invalid payload", "field is required")
	})

	resp, err := doTest(t, app, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}

	res := decodeJSON(t, resp.Body)
	if res.Success || res.Message != "invalid payload" || res.Error != "field is required" {
		t.Fatalf("unexpected payload: %+v", res)
	}
}

func TestRespUnauthorized(t *testing.T) {
	app := newTestApp(func(c fiber.Ctx) error {
		return cores.RespUnauthorized(c, "missing token")
	})

	resp, err := doTest(t, app, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}

	res := decodeJSON(t, resp.Body)
	if res.Success || res.Message != "missing token" {
		t.Fatalf("unexpected payload: %+v", res)
	}
}

func TestRespNotFound(t *testing.T) {
	app := newTestApp(func(c fiber.Ctx) error {
		return cores.RespNotFound(c, "not found")
	})

	resp, err := doTest(t, app, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}

	res := decodeJSON(t, resp.Body)
	if res.Success || res.Message != "not found" {
		t.Fatalf("unexpected payload: %+v", res)
	}
}

func TestRespInternalError(t *testing.T) {
	app := newTestApp(func(c fiber.Ctx) error {
		return cores.RespInternalError(c, "something failed", io.ErrUnexpectedEOF)
	})

	resp, err := doTest(t, app, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
	}

	// err must not leak to the client
	res := decodeJSON(t, resp.Body)
	if res.Success || res.Message != "Internal Server Error" || res.Error != nil {
		t.Fatalf("unexpected payload: %+v", res)
	}
}

func TestMsgPackResponse(t *testing.T) {
	app := newTestApp(func(c fiber.Ctx) error {
		return cores.RespSuccess(c, "msgpack", fiber.Map{"id": 1})
	})

	resp, err := doTest(t, app, "application/x-msgpack")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/x-msgpack" {
		t.Fatalf("expected msgpack content type, got %q", got)
	}

	var res cores.BaseResponse
	if err := msgpack.Unmarshal(readAll(t, resp.Body), &res); err != nil {
		t.Fatalf("failed to decode msgpack response: %v", err)
	}
	if !res.Success || res.Message != "msgpack" {
		t.Fatalf("unexpected payload: %+v", res)
	}
}

func readAll(t *testing.T, r io.Reader) []byte {
	t.Helper()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	return b
}
