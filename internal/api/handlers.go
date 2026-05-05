package api

import (
	"ivf-golang/internal/models"
	"ivf-golang/internal/vector"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type ApiState struct {
	Normalization models.Normalization
	MccRisk       models.MccRisk
	IVF           *models.IVFIndex
}

func (state *ApiState) IsReadyHandler(c fiber.Ctx) error {
	c.Status(http.StatusNoContent)

	if state.IVF == nil {
		c.Status(http.StatusServiceUnavailable)
	}

	return nil
}

func (state *ApiState) ProcessPaymentFraudHandler(c fiber.Ctx) error {
	if state.IVF == nil {
		c.Status(http.StatusServiceUnavailable)
		return nil
	}

	var transactionRequest models.FraudScoreRequest
	if err := c.Bind().JSON(&transactionRequest); err != nil {
		return err
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
