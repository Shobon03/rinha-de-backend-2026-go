package vector

import (
	"encoding/binary"
	"ivf-golang/internal/models"
	"ivf-golang/internal/util"
	"os"
	"path/filepath"
	"syscall"
)

const NumCentroids = 1732

func LoadIVF() *models.IVFIndex {
	path := filepath.Join("resources", "ivf.bin")
	file, err := os.Open(path)
	util.Check(err)

	fileInfo, err := file.Stat()
	util.Check(err)
	size := fileInfo.Size()

	// Load centroids (first bytes)
	centroids := make(models.Centroids, NumCentroids)
	binary.Read(file, binary.LittleEndian, &centroids)

	// Load bucket indexes
	indexes := make([]models.FlatIndex, NumCentroids)
	binary.Read(file, binary.LittleEndian, &indexes)

	// Load reference data using mmap
	data, err := syscall.Mmap(int(file.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
	util.Check(err)

	return &models.IVFIndex{
		Centroids:     centroids,
		BucketIndexes: indexes,
		Data:          data,
	}
}
