# Dual-Format Support: JSON & MessagePack

`fiber-starter` provides end-to-end support for both **JSON** and **MessagePack** across the entire request/response lifecycle:
- **Input (Request Body Parsing)** via `c.Bind().Body(&req)` or `c.Bind().MsgPack(&req)`.
- **Output (Response Serialization)** via [`cores/response.go`](../cores/response.go) with automatic HTTP Content Negotiation.

---

## 🎯 Features

- **Standardized Response Envelope**: Unified schema (`success`, `message`, `data`, and `error`) with `json` and `msgpack` struct tags.
- **Bi-Directional MessagePack Support**: Parse incoming MessagePack request bodies and respond with MessagePack seamlessly.
- **Supported MIME Types**:
  - `application/x-msgpack` (common de-facto standard)
  - `application/msgpack`
  - `application/vnd.msgpack` (Fiber v3 standard)
- **Automatic Content Negotiation**: Honors the client's `Accept` header and quality factors (`q`, e.g. `Accept: application/json;q=0.5, application/x-msgpack;q=0.9`).
- **Validation Integration**: `StructValidator` (`go-playground/validator`) executes automatically on decoded MessagePack payloads.

---

## 📥 Inbound Request Body Binding (Input)

Fiber v3 parses incoming request bodies based on the request's `Content-Type` header. `MsgPackBinder` is registered out-of-the-box in `cores.CreateApp()`.

### Handler Usage
In your handler, simply call `c.Bind().Body(&req)`:

```go
type CreateUserRequest struct {
	Name  string `json:"name"  msgpack:"name"  validate:"required"`
	Email string `json:"email" msgpack:"email" validate:"required,email"`
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return err // GlobalErrorHandler converts this to 400 Bad Request
	}

	return cores.RespCreated(c, "User created", req)
}
```

If the client sends:
- `Content-Type: application/json` ➡️ Fiber decodes JSON.
- `Content-Type: application/x-msgpack` (or `application/msgpack` / `application/vnd.msgpack`) ➡️ `MsgPackBinder` decodes MessagePack.
- Struct validation is automatically triggered for both formats.

You can also explicitly bind MessagePack using `c.Bind().MsgPack(&req)`.

---

## 📦 Envelope Schema (`BaseResponse`)

All API responses follow the standard `BaseResponse` struct:

```go
type BaseResponse struct {
	Success bool   `json:"success" msgpack:"success"`
	Message string `json:"message" msgpack:"message"`
	Data    any    `json:"data,omitempty" msgpack:"data,omitempty"`
	Error   any    `json:"error,omitempty" msgpack:"error,omitempty"`
}
```

---

## 🚀 Available Response Helpers

Use these helper functions inside your handlers:

| Helper Function | HTTP Status | Description |
| :--- | :--- | :--- |
| `cores.RespSuccess(c, msg, data)` | `200 OK` | Success response with payload |
| `cores.RespCreated(c, msg, data)` | `201 Created` | Resource creation response |
| `cores.RespBadReq(c, msg, err)` | `400 Bad Request` | Bad request or validation failure |
| `cores.RespUnauthorized(c, msg)` | `401 Unauthorized` | Authentication error |
| `cores.RespForbidden(c, msg)` | `403 Forbidden` | Authorization error |
| `cores.RespNotFound(c, msg)` | `404 Not Found` | Resource not found |
| `cores.RespConflict(c, msg, err)` | `409 Conflict` | Resource conflict or duplicate state |
| `cores.RespInternalError(c, msg, err)`| `500 Internal Server Error` | Internal error (logs error via Zap, hides stack trace from client) |

---

## 🔄 Content Negotiation (JSON vs MessagePack)

Content negotiation is handled automatically by `sendResponse`:

```go
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
```

---

## 💡 Example Requests & Responses

### 1. Default JSON Request
```http
POST /users HTTP/1.1
Content-Type: application/json
Accept: application/json

{"name": "John Doe", "email": "john@example.com"}
```
**Response (`Content-Type: application/json`):**
```json
{
  "success": true,
  "message": "User created",
  "data": { "name": "John Doe", "email": "john@example.com" }
}
```

### 2. MessagePack Request & Response (Full Binary Roundtrip)
```http
POST /users HTTP/1.1
Content-Type: application/x-msgpack
Accept: application/x-msgpack

<Binary MessagePack Payload>
```
**Response (`Content-Type: application/x-msgpack`):**
Binary MessagePack payload encoding `BaseResponse`.
