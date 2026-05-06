package api

import "net/http"

func RegisterRoutes(state *ApiState) {
	http.HandleFunc("/ready", state.IsReadyHandler)
	http.HandleFunc("/fraud-score", state.ProcessPaymentFraudHandler)
}
