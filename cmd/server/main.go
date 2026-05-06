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
)

func main() {
	state := &api.ApiState{}

	// Establish a Unix Socket Domain connection
	socketPath := os.Getenv("SOCKET_PATH")
	if socketPath == "" {
		socketPath = "/tmp/go-api.sock"
	}
	os.Remove(socketPath)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("FATAL ASYNC INIT: %v", r)
			}
		}()

		state.Normalization = util.LoadFileJson[models.Normalization]("normalization.json")
		state.MccRisk = util.LoadFileJson[models.MccRisk]("mcc_risk.json")
		state.IVF = vector.LoadIVF()
		log.Println("Ready to serve requests!")
	}()

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	os.Chmod(socketPath, 0777)

	api.RegisterRoutes(state)

	log.Println("Server listening on UDS:", socketPath)
	if err := http.Serve(listener, nil); err != nil {
		log.Fatal(err)
	}
}
