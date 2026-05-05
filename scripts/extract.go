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
	"math"
	"os"
	"path/filepath"

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
	fmt.Println("Selecting centroids")
	// Select centroids
	total := int(math.Sqrt(float64(len(references))))
	centroids = make([]Reference, total)

	for i := 0; i < total; i++ {
		// Divides references into total slices and selects the middle element as the centroid
		start := i * total
		end := min((i+1)*total, len(references))
		baseSlice := references[start:end]
		centroids[i] = baseSlice[len(baseSlice)/2]
	}
}

func AttributeBuckets() {
	fmt.Println("Starting to attribute buckets")

	buckets := make([][]Reference, len(centroids))

	// Calculate each point to its distance from each centroid and attribute it to the closest bucket
	for _, ref := range references {
		minDistance := float32(math.MaxFloat32)
		bucketIndex := 0
		for i, centroid := range centroids {
			distance := util.CalculateEuclidianDistance(ref.Vector, centroid.Vector)
			if distance < minDistance {
				minDistance = distance
				bucketIndex = i
			}
		}
		buckets[bucketIndex] = append(buckets[bucketIndex], ref)
	}

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
