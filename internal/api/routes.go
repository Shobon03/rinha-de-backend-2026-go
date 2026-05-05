package api

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App, state *ApiState) {
	app.Get("/ready", state.IsReadyHandler)
	app.Post("/fraud-score", state.ProcessPaymentFraudHandler)
}
