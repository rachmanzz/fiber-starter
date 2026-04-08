# Delivery Layer (HTTP Handlers) Documentation

The delivery layer, often referred to as the handler or controller layer, is responsible for receiving incoming requests (e.g., HTTP requests), parsing them, delegating business logic to the service layer, and sending back appropriate responses. This layer focuses on protocol-specific concerns and data marshaling.

## Handler Location

Handler files are typically organized within the `app/delivery/` directory. Each set of related handlers is usually named following the pattern `[resource_name]_handler.go`, for example, `app/delivery/user_handler.go`.

## Handler Structure Example

A typical handler includes a struct to hold its dependencies (like services) and methods that correspond to specific API endpoints.

Here's a simplified example of a handler for a `User` resource:

```go
package delivery

import (
	"strconv" // Required for parsing IDs
	"github.com/gofiber/fiber/v3"
	"github.com/rachmanzz/fiber-starter/app/services" // Assuming your services are here
	"github.com/rachmanzz/fiber-starter/cores"      // For standard responses
)

// UserHandler defines the structure for our user HTTP handlers.
type UserHandler struct {
	userService service.UserServiceInterface // Dependency on the user service
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

// GetUserByID handles GET requests to /users/:id
func (h *UserHandler) GetUserByID(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64) // Example for uint ID
	if err != nil {
		return cores.RespBadRequest(c, "Invalid user ID format", nil)
	}

	// Delegate to the service layer for business logic
	user, err := h.userService.GetByID(c.Context(), id)
	if err != nil {
		// Handle different service errors (e.g., not found)
		return cores.RespInternalServerError(c, "Failed to fetch user", err.Error())
	}

	return cores.RespSuccess(c, "User fetched successfully", user)
}

// CreateUser handles POST requests to /users
func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var req struct { // Define request body structure
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.Bind().Body(&req); err != nil {
		return cores.RespBadRequest(c, "Invalid request body", err.Error())
	}

	// Delegate to the service layer
	newUser, err := h.userService.Create(c.Context(), req.Name, req.Email)
	if err != nil {
		return cores.RespInternalServerError(c, "Failed to create user", err.Error())
	}

	return cores.RespCreated(c, "User created successfully", newUser)
}
```

In this example:
- `UserHandler` holds an instance of `service.UserServiceInterface`, injected through its constructor `NewUserHandler`.
- Methods like `GetUserByID` and `CreateUser` are responsible for handling specific HTTP routes.
- They parse request parameters/body, call the relevant service method, and format the response using `cores` helper functions.

## Flexibility in Structure

This documentation outlines a recommended and common practice for structuring delivery handlers within this boilerplate. However, this structure is not absolute. Developers are encouraged to adapt and evolve their handler architecture based on their project's specific needs, team conventions, and industry best practices, as long as it maintains clarity, testability, and maintainability.