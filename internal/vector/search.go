package vector

import (
	"ivf-golang/internal/models"
	"ivf-golang/internal/util"
	"slices"
	"unsafe"
)

var (
	N = 4 // Number of centroids to search
	K = 5 // Number of nearest neighbors
)

type CentroidDistance struct {
	Index    int
	Distance float32
}

type RecordDistance struct {
	Record   models.FlatRecord
	Distance float32
}

func GetCentroids(vector models.Vector, centroids models.Centroids) []CentroidDistance {
	// Maintain a sorted slice of the top N closest centroids.
	bestCentroids := make([]CentroidDistance, 0, N)

	for i, centroid := range centroids {
		dist := util.CalculateEuclidianDistance(vector, centroid)

		if len(bestCentroids) < N {
			bestCentroids = append(bestCentroids, CentroidDistance{Index: i, Distance: dist})
			slices.SortFunc(bestCentroids, func(a, b CentroidDistance) int {
				if a.Distance < b.Distance {
					return -1
				}
				if a.Distance > b.Distance {
					return 1
				}
				return 0
			})
		} else if dist < bestCentroids[N-1].Distance {
			bestCentroids[N-1] = CentroidDistance{Index: i, Distance: dist}
			slices.SortFunc(bestCentroids, func(a, b CentroidDistance) int {
				if a.Distance < b.Distance {
					return -1
				}
				if a.Distance > b.Distance {
					return 1
				}
				return 0
			})
		}
	}

	return bestCentroids
}

func GetBuckets(vector models.Vector, bestCentroids []CentroidDistance, ivf *models.IVFIndex) []RecordDistance {
	var bestK []RecordDistance

	// For each centroid, find the closest K vectors
	for _, centroid := range bestCentroids {
		bucketInfo := ivf.BucketIndexes[centroid.Index]
		if bucketInfo.Count == 0 {
			continue
		}

		// Offset: Centroids (1732*14*4) + Indexes (1732*8) = 110848
		// Record Size: 30 bytes
		offset := 110848 + (int64(bucketInfo.StartIndex) * 30)

		// Zero-copy: access the mmap'ed data directly from the byte slice
		// No ReadAt, no buffer allocation, no sync.Pool.
		records := unsafe.Slice((*models.FlatRecord)(unsafe.Pointer(&ivf.Data[offset])), bucketInfo.Count)

		for _, record := range records {
			dist := util.CalculateEuclidianDistanceQuantized(vector, record.Vector)

			if len(bestK) < K {
				bestK = append(bestK, RecordDistance{Record: record, Distance: dist})
				slices.SortFunc(bestK, func(a, b RecordDistance) int {
					if a.Distance < b.Distance {
						return -1
					}
					if a.Distance > b.Distance {
						return 1
					}
					return 0
				})
			} else if dist < bestK[K-1].Distance {
				bestK[K-1] = RecordDistance{Record: record, Distance: dist}
				slices.SortFunc(bestK, func(a, b RecordDistance) int {
					if a.Distance < b.Distance {
						return -1
					}
					if a.Distance > b.Distance {
						return 1
					}
					return 0
				})
			}
		}
	}

	return bestK
}

func SearchAndCheckFraudScore(bestRecords []RecordDistance) float32 {
	fraudVotes := float32(0.0)
	for _, record := range bestRecords {
		if record.Record.Label == 1 {
			fraudVotes++
		}
	}

	return fraudVotes / float32(K)
}
