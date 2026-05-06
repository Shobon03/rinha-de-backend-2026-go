package main

import (
	"ivf-golang/internal/api"
	"ivf-golang/internal/models"
	"ivf-golang/internal/util"
	"ivf-golang/internal/vector"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Open Unix Socket Domain and register routes
	socketPath := os.Getenv("SOCKET_PATH")
	if socketPath == "" {
		socketPath = "/sockets/go-api.sock"
	}

	// Ensure the directory exists
	os.MkdirAll(filepath.Dir(socketPath), 0777)
	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	os.Chmod(socketPath, 0777)
	log.Println("Socket UDS created in:", socketPath)

	state := &api.ApiState{
		Normalization: util.LoadFileJson[models.Normalization]("normalization.json"),
		MccRisk:       util.LoadFileJson[models.MccRisk]("mcc_risk.json"),
		IVF:           vector.LoadIVF(),
	}
	log.Println("Data loaded into memory!")

	api.RegisterRoutes(state)

	log.Println("Server listening on UDS:", socketPath)
	if err := http.Serve(listener, nil); err != nil {
		log.Fatal(err)
	}
}
