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
	state := &api.ApiState{
		Normalization: util.LoadFileJson[models.Normalization]("normalization.json"),
		MccRisk:       util.LoadFileJson[models.MccRisk]("mcc_risk.json"),
		IVF:           vector.LoadIVF(),
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	api.RegisterRoutes(app, state)

	log.Fatal(app.Listen(":9999", fiber.ListenConfig{
		EnablePrefork:         false,
		DisableStartupMessage: true,
	}))
}
