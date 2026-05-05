package util

import "ivf-golang/internal/models"

func CalculateEuclidianDistance(a models.Vector, b models.Vector) float32 {
	var sum float32
	for i := range a {
		dist := a[i] - b[i]
		sum += dist * dist
	}

	return sum
}

func CalculateEuclidianDistanceQuantized(a models.Vector, b [14]uint16) float32 {
	var sum float32
	for i := range a {
		bFloat := float32(b[i]) / 65535.0
		dist := a[i] - bFloat
		sum += dist * dist
	}

	return sum
}
