/**
 * Script to extract references from the JSON file and select centroids for the IVF index
 * Creates ivf.bin in resources/ folder
 */

package main

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"ivf-golang/internal/models"
	"ivf-golang/internal/util"
	"ivf-golang/internal/vector"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	json "github.com/goccy/go-json"
)

type Reference struct {
	Vector [14]float32 `json:"vector"`
	Label  string      `json:"label"`
}

var (
	references []Reference
	centroids  []Reference
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func ExtractReferences() {
	fmt.Println("Starting to extract references")

	// Loads file
	path := filepath.Join("resources", "references.json.gz")
	file, err := os.Open(path)
	Check(err)
	defer file.Close()

	// Unzips file
	fmt.Println("Unzipping and decoding file")
	gz, err := gzip.NewReader(file)
	Check(err)
	defer gz.Close()

	// Decodes via JSON
	decoder := json.NewDecoder(gz)

	err = decoder.Decode(&references)
	Check(err)

	fmt.Println("File decoded successfully")
}

func SelectCentroids() {
	fmt.Println("Selecting centroids with Parallel K-Means")
	total := vector.NumCentroids
	centroids = make([]Reference, total)

	// 1. Initial Selection: Regular intervals
	for i := 0; i < total; i++ {
		centroids[i] = references[(i*len(references))/total]
	}

	// 2. Iterations
	iterations := 30
	numWorkers := runtime.NumCPU()
	fmt.Printf("Using %d CPU cores for parallel processing\n", numWorkers)

	for iter := 0; iter < iterations; iter++ {
		fmt.Printf("K-Means Iteration %d/%d\n", iter+1, iterations)

		newCentroids := make([][14]float64, total)
		counts := make([]int, total)
		var mu sync.Mutex

		// Parallel Assignment Step
		var wg sync.WaitGroup
		chunkSize := (len(references) + numWorkers - 1) / numWorkers

		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				start := workerID * chunkSize
				end := start + chunkSize
				if end > len(references) {
					end = len(references)
				}
				if start >= end {
					return
				}

				localNewCentroids := make([][14]float64, total)
				localCounts := make([]int, total)

				for i := start; i < end; i++ {
					ref := references[i]
					minDist := float32(math.MaxFloat32)
					bestC := 0
					for j, c := range centroids {
						dist := util.CalculateEuclidianDistance(ref.Vector, c.Vector)
						if dist < minDist {
							minDist = dist
							bestC = j
						}
					}
					for j := 0; j < 14; j++ {
						localNewCentroids[bestC][j] += float64(ref.Vector[j])
					}
					localCounts[bestC]++
				}

				mu.Lock()
				for i := 0; i < total; i++ {
					for j := 0; j < 14; j++ {
						newCentroids[i][j] += localNewCentroids[i][j]
					}
					counts[i] += localCounts[i]
				}
				mu.Unlock()
			}(w)
		}
		wg.Wait()

		// Update step
		for i := 0; i < total; i++ {
			if counts[i] > 0 {
				var updatedVector [14]float32
				for j := 0; j < 14; j++ {
					updatedVector[j] = float32(newCentroids[i][j] / float64(counts[i]))
				}
				centroids[i] = Reference{Vector: updatedVector}
			}
		}
	}
}

func AttributeBuckets() {
	fmt.Println("Starting to attribute buckets")

	buckets := make([][]Reference, len(centroids))

	// Parallel bucket attribution
	numWorkers := runtime.NumCPU()
	var mu sync.Mutex
	var wg sync.WaitGroup
	chunkSize := (len(references) + numWorkers - 1) / numWorkers

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			start := workerID * chunkSize
			end := start + chunkSize
			if end > len(references) {
				end = len(references)
			}
			if start >= end {
				return
			}

			localBuckets := make([][]Reference, len(centroids))
			for i := start; i < end; i++ {
				ref := references[i]
				minDist := float32(math.MaxFloat32)
				bestC := 0
				for j, c := range centroids {
					dist := util.CalculateEuclidianDistance(ref.Vector, c.Vector)
					if dist < minDist {
						minDist = dist
						bestC = j
					}
				}
				localBuckets[bestC] = append(localBuckets[bestC], ref)
			}

			mu.Lock()
			for i := 0; i < len(centroids); i++ {
				buckets[i] = append(buckets[i], localBuckets[i]...)
			}
			mu.Unlock()
		}(w)
	}
	wg.Wait()

	fmt.Println("Buckets attributed, writing to file")

	// Create the binary file
	path := filepath.Join("resources", "ivf.bin")
	file, err := os.Create(path)
	Check(err)
	defer file.Close()

	// Block #1: Centroids [1732][14]float32
	centroidsVectors := make(models.Centroids, len(centroids))
	for i, c := range centroids {
		centroidsVectors[i] = c.Vector
	}
	err = binary.Write(file, binary.LittleEndian, centroidsVectors)
	Check(err)

	// Block #2: Centroid indexes (FlatIndex)
	bucketIndexes := make([]models.FlatIndex, len(centroids))
	var currentStartIndex uint32 = 0
	for i, bucket := range buckets {
		bucketIndexes[i] = models.FlatIndex{
			StartIndex: currentStartIndex,
			Count:      uint32(len(bucket)),
		}
		currentStartIndex += uint32(len(bucket))
	}
	err = binary.Write(file, binary.LittleEndian, bucketIndexes)
	Check(err)

	// Block #3: Records grouped by bucket (FlatRecord quantized)
	for _, bucket := range buckets {
		records := make([]models.FlatRecord, len(bucket))
		for i, ref := range bucket {
			label := uint8(0)
			if ref.Label == "fraud" {
				label = 1
			}

			var quantizedVector [14]uint16
			for j := 0; j < 14; j++ {
				val := ref.Vector[j]
				if val < 0 {
					val = 0
				}
				if val > 1 {
					val = 1
				}
				quantizedVector[j] = uint16(val * 65535)
			}

			records[i] = models.FlatRecord{
				Vector: quantizedVector,
				Label:  label,
			}
		}
		err = binary.Write(file, binary.LittleEndian, records)
		Check(err)
	}

	fmt.Println("IVF file written to", path)
}

func main() {
	ExtractReferences()
	SelectCentroids()
	AttributeBuckets()
}
