package vector

import (
	"ivf-golang/internal/models"
	"ivf-golang/internal/util"
	"slices"
)

func Clamp(x float32) float32 {
	if x < 0.0 {
		return 0.0
	}

	if x > 1.0 {
		return 1.0
	}

	return x
}

func NormalizeTransaction(
	normalization models.Normalization,
	mccRisk models.MccRisk,
	transaction models.FraudScoreRequest,
) models.Vector {
	amount := Clamp(float32(transaction.Transaction.Amount) / float32(normalization.MaxAmount))
	installments := Clamp(float32(transaction.Transaction.Installments) / float32(normalization.MaxInstallments))
	amountVsAvg := Clamp((float32(transaction.Transaction.Amount) / float32(transaction.Customer.AvgAmount)) / float32(normalization.AmountVsAvgRatio))

	reqHour, reqWeekday, reqAbsMinutes := util.ParseFastDate(transaction.Transaction.RequestedAt)

	hourOfDay := float32(reqHour) / 23.0
	dayOfWeek := float32(reqWeekday) / 6.0

	var minutesSinceLastTx float32 = -1.0
	var kmFromLastTx float32 = -1.0

	if transaction.LastTransaction != nil {
		_, _, lastAbsMinutes := util.ParseFastDate(transaction.LastTransaction.Timestamp)

		minutes := float32(reqAbsMinutes - lastAbsMinutes)
		if minutes < 0 {
			minutes = 0.0
		}

		minutesSinceLastTx = Clamp(minutes / float32(normalization.MaxMinutes))
		kmFromLastTx = Clamp(float32(transaction.LastTransaction.KmFromCurrent) / float32(normalization.MaxKm))
	}

	kmFromHome := Clamp(float32(transaction.Terminal.KmFromHome) / float32(normalization.MaxKm))
	txCount24h := Clamp(float32(transaction.Customer.TxCount24h) / float32(normalization.MaxTxCount24h))

	var isOnline float32 = 0.0
	if transaction.Terminal.IsOnline {
		isOnline = 1.0
	}

	var cardPresent float32 = 0.0
	if transaction.Terminal.CardPresent {
		cardPresent = 1.0
	}

	var unknownMerchant float32 = 0.0
	if !slices.Contains(transaction.Customer.KnownMerchants, transaction.Merchant.ID) {
		unknownMerchant = 1.0
	}

	risk, ok := mccRisk[transaction.Merchant.MCC]
	if !ok {
		risk = 0.5
	}

	merchantAvgAmount := Clamp(float32(transaction.Merchant.AvgAmount) / float32(normalization.MaxMerchantAvgAmount))

	return models.Vector{
		amount,
		installments,
		amountVsAvg,
		hourOfDay,
		dayOfWeek,
		minutesSinceLastTx,
		kmFromLastTx,
		kmFromHome,
		txCount24h,
		isOnline,
		cardPresent,
		unknownMerchant,
		risk,
		merchantAvgAmount,
	}
}
