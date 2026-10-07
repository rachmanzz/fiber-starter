package cores_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/rachmanzz/fiber-starter/cores"
	"github.com/shamaton/msgpack/v3"
)

type sampleUserRequest struct {
	Name  string `json:"name"  msgpack:"name"  validate:"required"`
	Email string `json:"email" msgpack:"email" validate:"required,email"`
}

func setupTestContractApp(handler func(c fiber.Ctx) error) *fiber.App {
	contract := cores.CreateContract().CreateApp()
	contract.App.Post("/users", handler)
	return contract.App
}

func TestMsgPackBinder_Metadata(t *testing.T) {
	b := cores.NewMsgPackBinder()
	if b.Name() != "msgpack" {
		t.Fatalf("expected binder name 'msgpack', got %q", b.Name())
	}

	expectedMIMEs := []string{
		"application/x-msgpack",
		"application/msgpack",
		"application/vnd.msgpack",
	}
	mimes := b.MIMETypes()
	if len(mimes) != len(expectedMIMEs) {
		t.Fatalf("expected %d mime types, got %d", len(expectedMIMEs), len(mimes))
	}
	for i, m := range expectedMIMEs {
		if mimes[i] != m {
			t.Errorf("expected mime %q at index %d, got %q", m, i, mimes[i])
		}
	}
}

func TestMsgPackInput_SupportedMIMETypes(t *testing.T) {
	app := setupTestContractApp(func(c fiber.Ctx) error {
		var req sampleUserRequest
		if err := c.Bind().Body(&req); err != nil {
			return err
		}
		return cores.RespSuccess(c, "created", fiber.Map{
			"name":  req.Name,
			"email": req.Email,
		})
	})

	payload := sampleUserRequest{
		Name:  "Alice",
		Email: "alice@example.com",
	}
	bodyBytes, err := msgpack.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal msgpack: %v", err)
	}

	contentTypes := []string{
		"application/x-msgpack",
		"application/msgpack",
		"application/vnd.msgpack",
	}

	for _, ct := range contentTypes {
		t.Run("ContentType_"+ct, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", ct)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status 200, got %d", resp.StatusCode)
			}

			var res cores.BaseResponse
			res = decodeJSON(t, resp.Body)
			if !res.Success || res.Message != "created" {
				t.Fatalf("unexpected response payload: %+v", res)
			}
			data := res.Data.(map[string]any)
			if data["name"] != "Alice" || data["email"] != "alice@example.com" {
				t.Fatalf("unexpected data: %+v", data)
			}
		})
	}
}

func TestMsgPackInput_ValidationFailure(t *testing.T) {
	app := setupTestContractApp(func(c fiber.Ctx) error {
		var req sampleUserRequest
		if err := c.Bind().Body(&req); err != nil {
			return err
		}
		return cores.RespSuccess(c, "created", req)
	})

	// Invalid email and missing name
	invalidPayload := sampleUserRequest{
		Name:  "",
		Email: "not-an-email",
	}
	bodyBytes, err := msgpack.Marshal(invalidPayload)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/x-msgpack")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for validation failure, got %d", resp.StatusCode)
	}

	res := decodeJSON(t, resp.Body)
	if res.Success {
		t.Fatalf("expected failure, got success: %+v", res)
	}
	if res.Message != "Validation failed" {
		t.Fatalf("expected 'Validation failed', got %q", res.Message)
	}
}

func TestMsgPackInput_CorruptedPayload(t *testing.T) {
	app := setupTestContractApp(func(c fiber.Ctx) error {
		var req sampleUserRequest
		if err := c.Bind().Body(&req); err != nil {
			return err
		}
		return cores.RespSuccess(c, "created", req)
	})

	corruptBytes := []byte{0xc1, 0xff, 0x00, 0x12} // invalid msgpack byte sequence
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(corruptBytes))
	req.Header.Set("Content-Type", "application/x-msgpack")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for corrupt msgpack, got %d", resp.StatusCode)
	}

	res := decodeJSON(t, resp.Body)
	if res.Success || res.Message != "Invalid request payload" {
		t.Fatalf("expected 'Invalid request payload', got %+v", res)
	}
}

func TestMsgPack_FullRoundtrip(t *testing.T) {
	app := setupTestContractApp(func(c fiber.Ctx) error {
		var req sampleUserRequest
		if err := c.Bind().Body(&req); err != nil {
			return err
		}
		return cores.RespSuccess(c, "echo", fiber.Map{
			"name":  req.Name,
			"email": req.Email,
		})
	})

	payload := sampleUserRequest{
		Name:  "Bob",
		Email: "bob@example.com",
	}
	bodyBytes, err := msgpack.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/x-msgpack")
	req.Header.Set("Accept", "application/x-msgpack")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/x-msgpack" {
		t.Fatalf("expected Content-Type application/x-msgpack, got %q", got)
	}

	var res cores.BaseResponse
	if err := msgpack.Unmarshal(readAll(t, resp.Body), &res); err != nil {
		t.Fatalf("failed to unmarshal msgpack response: %v", err)
	}

	if !res.Success || res.Message != "echo" {
		t.Fatalf("unexpected payload: %+v", res)
	}
}

func TestMsgPack_DirectBindMsgPack(t *testing.T) {
	app := setupTestContractApp(func(c fiber.Ctx) error {
		var req sampleUserRequest
		if err := c.Bind().MsgPack(&req); err != nil {
			return err
		}
		return cores.RespSuccess(c, "ok", req)
	})

	payload := sampleUserRequest{
		Name:  "Charlie",
		Email: "charlie@example.com",
	}
	bodyBytes, err := msgpack.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestMsgPackInput_EmptyBody(t *testing.T) {
	app := setupTestContractApp(func(c fiber.Ctx) error {
		var req sampleUserRequest
		if err := c.Bind().Body(&req); err != nil {
			return err
		}
		return cores.RespSuccess(c, "ok", req)
	})

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/x-msgpack")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for empty body, got %d", resp.StatusCode)
	}
}

func TestMsgPack_ConcurrentSafety(t *testing.T) {
	app := setupTestContractApp(func(c fiber.Ctx) error {
		var req sampleUserRequest
		if err := c.Bind().Body(&req); err != nil {
			return err
		}
		return cores.RespSuccess(c, "ok", req)
	})

	payload := sampleUserRequest{
		Name:  "Concurrent Tester",
		Email: "tester@example.com",
	}
	bodyBytes, err := msgpack.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal msgpack: %v", err)
	}

	const concurrency = 50
	errChan := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/x-msgpack")
			req.Header.Set("Accept", "application/x-msgpack")

			resp, err := app.Test(req)
			if err != nil {
				errChan <- err
				return
			}
			if resp.StatusCode != http.StatusOK {
				errChan <- err
				return
			}
			var res cores.BaseResponse
			if err := msgpack.Unmarshal(readAll(t, resp.Body), &res); err != nil {
				errChan <- err
				return
			}
			if !res.Success {
				errChan <- err
				return
			}
			errChan <- nil
		}()
	}

	for i := 0; i < concurrency; i++ {
		if err := <-errChan; err != nil {
			t.Fatalf("concurrent test failed: %v", err)
		}
	}
}

type binaryDataRequest struct {
	Data []byte `msgpack:"data"`
}

func TestSafeUnmarshal_BufferIsolation(t *testing.T) {
	originalData := []byte("confidential-payload")
	payload := binaryDataRequest{Data: originalData}

	encoded, err := msgpack.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var req binaryDataRequest
	if err := cores.SafeUnmarshal(encoded, &req); err != nil {
		t.Fatalf("SafeUnmarshal error: %v", err)
	}

	// Mutate the original encoded buffer (simulating fasthttp buffer reuse)
	for i := range encoded {
		encoded[i] = 0xff
	}

	// Verify that decoded req.Data remains uncorrupted
	if !bytes.Equal(req.Data, originalData) {
		t.Fatalf("memory corruption detected! expected %q, got %q", originalData, req.Data)
	}
}

func TestSafeUnmarshal_NilDestination(t *testing.T) {
	err := cores.SafeUnmarshal([]byte{0x90}, nil)
	if err == nil {
		t.Fatal("expected error when destination is nil, got nil")
	}
}


