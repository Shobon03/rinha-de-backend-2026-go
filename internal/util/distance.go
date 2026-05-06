package util

import "ivf-golang/internal/models"

func CalculateEuclidianDistance(a models.Vector, b models.Vector) float32 {
	var sum float32

	// Manual unrolling for 14 dims
	d0 := a[0] - b[0]
	sum += d0 * d0

	d1 := a[1] - b[1]
	sum += d1 * d1

	d2 := a[2] - b[2]
	sum += d2 * d2

	d3 := a[3] - b[3]
	sum += d3 * d3

	d4 := a[4] - b[4]
	sum += d4 * d4

	d5 := a[5] - b[5]
	sum += d5 * d5

	d6 := a[6] - b[6]
	sum += d6 * d6

	d7 := a[7] - b[7]
	sum += d7 * d7

	d8 := a[8] - b[8]
	sum += d8 * d8

	d9 := a[9] - b[9]
	sum += d9 * d9

	d10 := a[10] - b[10]
	sum += d10 * d10

	d11 := a[11] - b[11]
	sum += d11 * d11

	d12 := a[12] - b[12]
	sum += d12 * d12

	d13 := a[13] - b[13]
	sum += d13 * d13

	return sum
}

func QuantizeVector(v models.Vector) [14]uint16 {
	var res [14]uint16
	for i := range 14 {
		val := v[i]
		if val < 0 {
			val = 0
		}
		if val > 1 {
			val = 1
		}
		res[i] = uint16(val * 65535.0)
	}
	return res
}

func CalculateEuclidianDistanceInteger(a [14]uint16, b [14]uint16, threshold uint64) uint64 {
	var sum uint64

	d0 := int32(a[0]) - int32(b[0])
	sum += uint64(int64(d0) * int64(d0))
	if sum >= threshold {
		return sum
	}

	d1 := int32(a[1]) - int32(b[1])
	sum += uint64(int64(d1) * int64(d1))
	if sum >= threshold {
		return sum
	}

	d2 := int32(a[2]) - int32(b[2])
	sum += uint64(int64(d2) * int64(d2))
	if sum >= threshold {
		return sum
	}

	d3 := int32(a[3]) - int32(b[3])
	sum += uint64(int64(d3) * int64(d3))
	if sum >= threshold {
		return sum
	}

	d4 := int32(a[4]) - int32(b[4])
	sum += uint64(int64(d4) * int64(d4))
	if sum >= threshold {
		return sum
	}

	d5 := int32(a[5]) - int32(b[5])
	sum += uint64(int64(d5) * int64(d5))
	if sum >= threshold {
		return sum
	}

	d6 := int32(a[6]) - int32(b[6])
	sum += uint64(int64(d6) * int64(d6))
	if sum >= threshold {
		return sum
	}

	d7 := int32(a[7]) - int32(b[7])
	sum += uint64(int64(d7) * int64(d7))
	if sum >= threshold {
		return sum
	}

	d8 := int32(a[8]) - int32(b[8])
	sum += uint64(int64(d8) * int64(d8))
	if sum >= threshold {
		return sum
	}

	d9 := int32(a[9]) - int32(b[9])
	sum += uint64(int64(d9) * int64(d9))
	if sum >= threshold {
		return sum
	}

	d10 := int32(a[10]) - int32(b[10])
	sum += uint64(int64(d10) * int64(d10))
	if sum >= threshold {
		return sum
	}

	d11 := int32(a[11]) - int32(b[11])
	sum += uint64(int64(d11) * int64(d11))
	if sum >= threshold {
		return sum
	}

	d12 := int32(a[12]) - int32(b[12])
	sum += uint64(int64(d12) * int64(d12))
	if sum >= threshold {
		return sum
	}

	d13 := int32(a[13]) - int32(b[13])
	sum += uint64(int64(d13) * int64(d13))

	return sum
}

func CalculateEuclidianDistanceQuantized(a models.Vector, b [14]uint16, threshold float32) float32 {
	const inv65535 = 1.0 / 65535.0
	var sum float32

	for i := range 14 {
		diff := a[i] - float32(b[i])*inv65535
		sum += diff * diff
		if sum >= threshold {
			return sum
		}
	}

	return sum
}
