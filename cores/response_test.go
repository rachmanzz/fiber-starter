package cores_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestRespForbidden(t *testing.T) {
	app := newTestApp(func(c fiber.Ctx) error {
		return cores.RespForbidden(c, "access denied")
	})

	resp, err := doTest(t, app, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}

	res := decodeJSON(t, resp.Body)
	if res.Success || res.Message != "access denied" {
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

	mimes := []string{
		"application/x-msgpack",
		"application/msgpack",
		"application/vnd.msgpack",
	}

	for _, mime := range mimes {
		t.Run("Accept_"+mime, func(t *testing.T) {
			resp, err := doTest(t, app, mime)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.StatusCode != fiber.StatusOK {
				t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
			}
			if got := resp.Header.Get("Content-Type"); got != mime {
				t.Fatalf("expected content type %q, got %q", mime, got)
			}

			var res cores.BaseResponse
			if err := msgpack.Unmarshal(readAll(t, resp.Body), &res); err != nil {
				t.Fatalf("failed to decode msgpack response: %v", err)
			}
			if !res.Success || res.Message != "msgpack" {
				t.Fatalf("unexpected payload: %+v", res)
			}
		})
	}
}

func TestContentNegotiation(t *testing.T) {
	app := newTestApp(func(c fiber.Ctx) error {
		return cores.RespSuccess(c, "negotiated", fiber.Map{"id": 1})
	})

	t.Run("Weighted MsgPack Priority", func(t *testing.T) {
		resp, err := doTest(t, app, "application/json;q=0.5, application/x-msgpack;q=0.9")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := resp.Header.Get("Content-Type"); got != "application/x-msgpack" {
			t.Fatalf("expected application/x-msgpack, got %q", got)
		}
	})

	t.Run("Weighted JSON Priority", func(t *testing.T) {
		resp, err := doTest(t, app, "application/x-msgpack;q=0.5, application/json;q=0.9")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := resp.Header.Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
			t.Fatalf("expected application/json, got %q", got)
		}
	})
}

func readAll(t *testing.T, r io.Reader) []byte {
	t.Helper()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	return b
}
