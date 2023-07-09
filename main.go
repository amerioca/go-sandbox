package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/antoniopapa/go-admin/database"
	"github.com/antoniopapa/go-admin/routes"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Loaded .env file")
}

func main() {
	database.Connect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		// AllowOrigins:     "*",
		// AllowHeaders:     "Origin, Content-Type, Accept",
		// AllowMethods:     "GET, POST, PUT, DELETE",
	}))

	routes.Setup(app)

	// http Port

	// fmt.Println("CERT:", os.Getenv("CERT"), os.Getenv("PRIVKEY"))
	// https://github.com/gofiber/recipes/blob/master/https-tls/main.go
	// Create tls certificate
	if false {
		log.Fatal(app.Listen(":8000"))
	} else {
		cert, err := tls.LoadX509KeyPair(os.Getenv("CERT"), os.Getenv("PRIVKEY"))
		if err != nil {
			log.Fatal(err)
		}
		cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
		ln, err := tls.Listen("tcp", os.Getenv("PORT"), cfg)
		if err != nil {
			panic(err)
		}
		// Start server with https/ssl enabled on http://localhost:443
		log.Fatal(app.Listener(ln))
		// log.Fatal(app.Listen(":8000"))
	}
}
