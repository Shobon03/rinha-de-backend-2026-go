package models

/**
 * Internal data manipulation
 */
type FlatRecord struct {
	Vector  [14]uint16
	Label   uint8
	Padding [3]uint8 // 28 + 1 + 3 = 32 bytes for cache alignment
}

type FlatIndex struct {
	StartIndex uint32
	Count      uint32
}

type Centroids [][14]float32

type IVFIndex struct {
	Centroids     Centroids
	BucketIndexes []FlatIndex
	Data          []byte
}

type Normalization struct {
	MaxAmount            int `json:"max_amount"`
	MaxInstallments      int `json:"max_installments"`
	AmountVsAvgRatio     int `json:"amount_vs_avg_ratio"`
	MaxMinutes           int `json:"max_minutes"`
	MaxKm                int `json:"max_km"`
	MaxTxCount24h        int `json:"max_tx_count_24h"`
	MaxMerchantAvgAmount int `json:"max_merchant_avg_amount"`
}
type MccRisk map[string]float32

type Vector [14]float32

type Transaction struct {
	Amount       float32 `json:"amount"`
	Installments int     `json:"installments"`
	RequestedAt  string  `json:"requested_at"`
}

type Customer struct {
	AvgAmount      float32  `json:"avg_amount"`
	TxCount24h     int      `json:"tx_count_24h"`
	KnownMerchants []string `json:"known_merchants"`
}

type Merchant struct {
	ID        string  `json:"id"`
	MCC       string  `json:"mcc"`
	AvgAmount float32 `json:"avg_amount"`
}

type Terminal struct {
	IsOnline    bool    `json:"is_online"`
	CardPresent bool    `json:"card_present"`
	KmFromHome  float32 `json:"km_from_home"`
}

type LastTransaction struct {
	Timestamp     string  `json:"timestamp"`
	KmFromCurrent float32 `json:"km_from_current"`
}

/**
 * Requests and responses
 */
type FraudScoreRequest struct {
	ID              string           `json:"id"`
	Transaction     Transaction      `json:"transaction"`
	Customer        Customer         `json:"customer"`
	Merchant        Merchant         `json:"merchant"`
	Terminal        Terminal         `json:"terminal"`
	LastTransaction *LastTransaction `json:"last_transaction"`
}

type FraudScoreResponse struct {
	Approved   bool    `json:"approved"`
	FraudScore float32 `json:"fraud_score"`
}

type ErrorMessageResponse struct {
	Message string `json:"message"`
}
