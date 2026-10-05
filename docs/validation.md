# Request Validation

`fiber-starter` provides integrated struct binding and request validation out of the box using **Fiber v3 Native Binding** and **`go-playground/validator/v10`**.

---

## 🎯 Features

- **Automated Body & Query Parsing**: Bind incoming requests using Fiber v3 `c.Bind().Body(&req)` or `c.Bind().Query(&req)`.
- **Automatic Tag Validation**: Validates struct fields automatically using standard `validate:"..."` tags during binding.
- **Centralized Error Responses**: Validation failures automatically trigger structured error responses handled by the Global Error Handler.

---

## 🚀 How to Use

### 1. Define Request Struct with Validation Tags

Use standard `json` and `validate` struct tags:

```go
package dto

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"gte=18"`
}
```

---

### 2. Bind Request in Controller / Handler

Call `c.Bind().Body(&req)` (or `.Query()`, `.Params()`) directly. If validation fails or payload is malformed, `c.Bind()` returns an error that can be handled or returned directly to the Global Error Handler.

```go
package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rachmanzz/fiber-starter/cores"
)

func CreateUserHandler(c fiber.Ctx) error {
	var req CreateUserRequest

	// Automatically parses payload AND validates struct tags
	if err := c.Bind().Body(&req); err != nil {
		return err // Delegates error formatting to Global Error Handler
	}

	return cores.RespSuccess(c, "User created successfully", req)
}
```

---

## ⚙️ Architecture & Configuration

The struct validator is registered during application setup in [`cores/contract.go`](../cores/contract.go):

```go
cfg := fiber.Config{
    AppName:         Config().App.Name,
    StructValidator: NewStructValidator(), // Registered automatically
}
```

The adapter implementation [`cores/validator.go`](../cores/validator.go) satisfies Fiber v3's `fiber.StructValidator` interface:

```go
type StructValidator struct {
	validator *validator.Validate
}

func (v *StructValidator) Validate(out any) error {
	return v.validator.Struct(out)
}
```

---

## 📊 Default Validation Error Output Format

When a request fails validation, the Global Error Handler returns HTTP `400 Bad Request` with a field-by-field error map:

```json
{
  "success": false,
  "message": "Validation failed",
  "error": {
    "Email": "email",
    "Name": "min"
  }
}
```
