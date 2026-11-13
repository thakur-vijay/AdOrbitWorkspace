package main

// @title AdOrbit Workspace API
// @version 1.0
// @description Workspace service for AdOrbit.
// @host localhost:8081
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT token with Bearer prefix.

import (
	"log"

	swaggerDocs "github.com/adorbit/workspace/docs"
	"github.com/adorbit/workspace/internal/app"
	"github.com/adorbit/workspace/internal/config"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	setupSwagger(cfg.SWAGGER_HOST, cfg.Port)

	application, err := app.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	defer application.Close()

	if err := application.Run(cfg.Port); err != nil {
		log.Fatal(err)
	}

}

func setupSwagger(swaggerHost, authPort string) {
	if swaggerHost == "" {
		swaggerHost = "localhost:" + authPort
	}

	swaggerDocs.SwaggerInfo.Host = swaggerHost
	swaggerDocs.SwaggerInfo.BasePath = "/api/v1"

	if swaggerHost == "localhost:"+authPort {
		swaggerDocs.SwaggerInfo.Schemes = []string{"http"}
	} else {
		swaggerDocs.SwaggerInfo.Schemes = []string{"https"}
	}
}
