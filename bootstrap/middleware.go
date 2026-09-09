package bootstrap

import (
	"runtime/debug"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"go.uber.org/zap"
)

func RegisterMiddleware(app *fiber.App) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, e any) {
			zap.L().Error("panic recovered",
				zap.Any("panic", e),
				zap.ByteString("stack", debug.Stack()),
			)
		},
		PanicHandler: func(c fiber.Ctx, r any) error {
			return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
		},
	}))
}
