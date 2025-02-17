package main

import (
	"fmt"
	"os"
	"studentapifiber/db"
	"studentapifiber/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(fmt.Errorf("error loading .env file: %v", err))
	}

	err = utils.ValidateEnv()
	if err != nil {
		log.Fatal(fmt.Errorf("error validating environment variables: %v", err))
	}

	err = db.InitDbPool()
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing database pool: %v", err))
	}

	pool, _ := db.GetDbPool()
	defer pool.Close()

	app := fiber.New()
	err = registerMiddlewares(app)
	if err != nil {
		log.Fatal(fmt.Errorf("error registering middlewares: %v", err))
	}
	err = registerRoutes(app)
	if err != nil {
		log.Fatal(fmt.Errorf("error registering routes: %v", err))
	}

	port, _ := os.LookupEnv("PORT")

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
