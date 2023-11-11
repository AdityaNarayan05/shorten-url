package main

import(
	"fmt"
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App){
	app.Get("/:url",routes.ResolveURL)
	app.post("/api/v1",routes.ShortenURL)
}

func main(){
	err := godotenv.Load()

	if err != nil{
		fmt.Println(err)
	}

	// Initialize the app with some configs
	app := fiber.New()

	app.Use(logger.New())

	setupRoutes(app)

	log.Fatal(app.listen(os.Getenv("APP_PORT")))
}