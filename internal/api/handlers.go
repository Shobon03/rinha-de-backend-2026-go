package api

import (
	"io"
	"ivf-golang/internal/models"
	"ivf-golang/internal/vector"
	"net/http"
	"sync"

	"github.com/goccy/go-json"
)

type ApiState struct {
	Normalization models.Normalization
	MccRisk       models.MccRisk
	IVF           *models.IVFIndex
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 4096) // 4KB for the fraud-score payload
	},
}

func (state *ApiState) IsReadyHandler(w http.ResponseWriter, r *http.Request) {
	if state.IVF == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (state *ApiState) ProcessPaymentFraudHandler(w http.ResponseWriter, r *http.Request) {
	if state.IVF == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Get a buffer from the pool to avoid heap allocation per request
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)

	// Read body into the pre-allocated buffer
	n, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var transactionRequest models.FraudScoreRequest
	if err := json.Unmarshal(buf[:n], &transactionRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Normalize the request
	normalized := vector.NormalizeTransaction(state.Normalization, state.MccRisk, transactionRequest)

	// Manipulate data using IVF
	bestCentroids := vector.GetCentroids(normalized, state.IVF.Centroids)
	bestRecords := vector.GetBuckets(normalized, bestCentroids, state.IVF)

	// Calculate fraud score
	fraudScore := vector.SearchAndCheckFraudScore(bestRecords)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.FraudScoreResponse{
		Approved:   fraudScore < 0.6,
		FraudScore: fraudScore,
	})
}
