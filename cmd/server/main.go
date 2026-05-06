package main

import (
	"ivf-golang/internal/api"
	"ivf-golang/internal/models"
	"ivf-golang/internal/util"
	"ivf-golang/internal/vector"
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	state := &api.ApiState{}
	api.RegisterRoutes(app, state)

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

	log.Println("Server listening on TCP :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
