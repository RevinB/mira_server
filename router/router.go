package router

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3"
	"runtime/debug"
)

var jwtMiddleware fiber.Handler

func NewRouter(secretKey string) *fiber.App {
	r := fiber.New(fiber.Config{ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		if e, ok := err.(*fiber.Error); ok {
			return ctx.SendStatus(e.Code)
		}

		stackString := string(debug.Stack())
		sentry.CaptureMessage(fmt.Sprintf("ERROR HANDLER\nCTX: %v\nError Message: %v\nStack Trace: %v",
			ctx.String(), err, stackString))

		_ = ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error. Please try again later.\n\nThis issue has been reported.")

		return nil
	}})

	r.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx *fiber.Ctx, e interface{}) {
			stackString := string(debug.Stack())
			sentry.CaptureMessage(fmt.Sprintf("RECOVER MIDDLEWARE\nCTX: %v\nError Message: %v\nStack Trace: %v",
				ctx.String(), e, stackString))
		},
	}))

	r.Use(csrf.New())

	jwtMiddleware = jwtware.New(jwtware.Config{
		SigningKey: secretKey,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.SendStatus(fiber.StatusForbidden)
		},
	})

	return r
}

func GetJwtHandler() fiber.Handler {
	return jwtMiddleware
}
