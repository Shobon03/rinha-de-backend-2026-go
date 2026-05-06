package vector

import (
	"ivf-golang/internal/models"
	"ivf-golang/internal/util"
	"math"
	"unsafe"
)

var (
	N = 4 // Number of centroids to search (Optimized for K-Means)
	K = 5 // Number of nearest neighbors
)

type CentroidDistance struct {
	Index    int
	Distance float32
}

type RecordDistance struct {
	Record   models.FlatRecord
	Distance uint64
}

func GetCentroids(vector models.Vector, centroids models.Centroids) []CentroidDistance {
	var best [4]CentroidDistance
	for i := range N {
		best[i].Distance = math.MaxFloat32
		best[i].Index = -1
	}

	for i, centroid := range centroids {
		dist := util.CalculateEuclidianDistance(vector, centroid, best[3].Distance)

		if dist < best[3].Distance {
			best[3] = CentroidDistance{Index: i, Distance: dist}
			for j := 3; j > 0 && best[j].Distance < best[j-1].Distance; j-- {
				best[j], best[j-1] = best[j-1], best[j]
			}
		}
	}

	return best[:]
}

func GetBuckets(vector models.Vector, bestCentroids []CentroidDistance, ivf *models.IVFIndex) []RecordDistance {
	q := util.QuantizeVector(vector)
	bestK := make([]RecordDistance, 0, K)

	for _, centroid := range bestCentroids {
		if centroid.Index == -1 {
			continue
		}
		bucketInfo := ivf.BucketIndexes[centroid.Index]
		if bucketInfo.Count == 0 {
			continue
		}

		offset := int64(NumCentroids)*64 + (int64(bucketInfo.StartIndex) * 32)
		records := unsafe.Slice((*models.FlatRecord)(unsafe.Pointer(&ivf.Data[offset])), bucketInfo.Count)

		for _, record := range records {
			var threshold uint64 = math.MaxUint64
			if len(bestK) == K {
				threshold = bestK[K-1].Distance
			}

			dist := util.CalculateEuclidianDistanceInteger(q, record.Vector, threshold)

			if dist < threshold {
				if len(bestK) < K {
					bestK = append(bestK, RecordDistance{Record: record, Distance: dist})
				} else {
					bestK[K-1] = RecordDistance{Record: record, Distance: dist}
				}
				for j := len(bestK) - 1; j > 0 && bestK[j].Distance < bestK[j-1].Distance; j-- {
					bestK[j], bestK[j-1] = bestK[j-1], bestK[j]
				}
			}
		}
	}

	return bestK
}

func SearchAndCheckFraudScore(bestRecords []RecordDistance) float32 {
	if len(bestRecords) == 0 {
		return 0.0
	}

	fraudVotes := float32(0.0)
	for _, record := range bestRecords {
		if record.Record.Label == 1 {
			fraudVotes++
		}
	}

	return fraudVotes / float32(K)
}
