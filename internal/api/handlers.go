package api

import (
	"ivf-golang/internal/models"
	"ivf-golang/internal/vector"

	"github.com/gofiber/fiber/v3"
)

type ApiState struct {
	Normalization models.Normalization
	MccRisk       models.MccRisk
	IVF           *models.IVFIndex
}

func (state *ApiState) IsReadyHandler(c fiber.Ctx) error {
	if state.IVF == nil {
		return c.SendStatus(fiber.StatusServiceUnavailable)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (state *ApiState) ProcessPaymentFraudHandler(c fiber.Ctx) error {
	if state.IVF == nil {
		return c.SendStatus(fiber.StatusServiceUnavailable)
	}

	var transactionRequest models.FraudScoreRequest
	if err := c.Bind().JSON(&transactionRequest); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Normalize the request
	normalized := vector.NormalizeTransaction(state.Normalization, state.MccRisk, transactionRequest)

	// Manipulate data using IVF
	bestCentroids := vector.GetCentroids(normalized, state.IVF.Centroids)
	bestRecords := vector.GetBuckets(normalized, bestCentroids, state.IVF)

	// Calculate fraud score
	fraudScore := vector.SearchAndCheckFraudScore(bestRecords)

	return c.JSON(models.FraudScoreResponse{
		Approved:   fraudScore < 0.6,
		FraudScore: fraudScore,
	})
}
