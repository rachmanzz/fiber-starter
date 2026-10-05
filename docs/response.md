# Standard API Response & Content Negotiation

`fiber-starter` provides standardized API response helpers in [`cores/response.go`](../cores/response.go) with built-in **HTTP Content Negotiation** for **JSON** and **MessagePack**.

---

## 🎯 Features

- **Standardized Response Schema**: Unified envelope containing `success`, `message`, `data`, and `error` fields.
- **Native Content Negotiation**: Automatically serves **JSON** or **MessagePack** based on the client's `Accept` header using Fiber v3's `c.Accepts("application/json", "application/x-msgpack")`.
- **Quality Factor (`q`) Support**: Honors weighted client headers (e.g. `Accept: application/json;q=0.5, application/x-msgpack;q=0.9`).

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
| `cores.RespNotFound(c, msg)` | `404 Not Found` | Resource not found |
| `cores.RespInternalError(c, msg, err)`| `500 Internal Server Error` | Internal error (logs error via Zap, hides stack trace from client) |

---

## 🔄 Content Negotiation (JSON vs MessagePack)

Content negotiation is handled automatically by `sendResponse`:

```go
func sendResponse(c fiber.Ctx, status int, payload BaseResponse) error {
	match := c.Accepts("application/json", "application/x-msgpack")
	if match == "application/x-msgpack" {
		b, err := msgpack.Marshal(payload)
		if err != nil {
			zap.L().Error("failed to marshal msgpack", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}
		c.Set("Content-Type", "application/x-msgpack")
		return c.Status(status).Send(b)
	}

	return c.Status(status).JSON(payload)
}
```

### Example Requests & Responses

1. **Default JSON Request**:
   ```http
   GET /users HTTP/1.1
   Accept: application/json
   ```
   **Response (`Content-Type: application/json`):**
   ```json
   {
     "success": true,
     "message": "User fetched successfully",
     "data": { "id": 1, "name": "John Doe" }
   }
   ```

2. **MessagePack Request**:
   ```http
   GET /users HTTP/1.1
   Accept: application/x-msgpack
   ```
   **Response (`Content-Type: application/x-msgpack`):**
   Binary MessagePack payload encoding `BaseResponse`.
