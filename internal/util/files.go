package util

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func LoadFileJson[T any](filename string) T {
	path := filepath.Join("resources", filename)
	f, err := os.Open(path)
	Check(err)
	defer f.Close()

	var target T
	err = json.NewDecoder(f).Decode(&target)
	Check(err)

	return target
}
