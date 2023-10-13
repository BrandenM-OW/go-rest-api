// @title Example Article API
// @version 1.0
// @description This is a sample swagger for Fiber
// @contact.name Branden Morin
// @contact.email brandenmorin14@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost
// @BasePath /
package main

import (
	"fmt"

	"github.com/BrandenM-PM/go-rest-api/controllers"
	_ "github.com/BrandenM-PM/go-rest-api/docs"
	"github.com/BrandenM-PM/go-rest-api/initializers"
	"github.com/BrandenM-PM/go-rest-api/migrations"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/handlebars"

	// "github.com/ansrivas/fiberprometheus/v2"
	"log"
	"os"
	"strings"
	"time"
)

func init() {
	initializers.LoadEnvVars("./")
	initializers.ConnectToPostgresDB()
}

func main() {
	engine := handlebars.New("./views", ".hbs")
	app := fiber.New(fiber.Config{
		ErrorHandler: initializers.CustomErrorHandler,
		Views:        engine,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://localhost:8080",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PATCH, DELETE",
	}))

	app.Use(recover.New())
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return strings.Contains(c.Path(), "/swagger") // don't limit swagger
		},
		Expiration: 10 * time.Second,
		Max:        300,
	}))
	app.Use(logger.New())
	// TODO: Add authentication middleware

	app.Static("/", "./public")
	app.Static("/purecss", "./node_modules/purecss/build")

	app.Get("/", controllers.Index)
	app.Get("/articles", controllers.GetAllArticles)
	app.Get("/articles/:id", controllers.GetArticle)
	app.Post("/articles", controllers.CreateArticle)
	app.Patch("/articles/:id", controllers.UpdateArticle)
	app.Delete("/articles/:id", controllers.DeleteArticle)

	app.Get("/swagger/*", swagger.HandlerDefault)

	// prometheus := fiberprometheus.New("go-rest-api")
	// prometheus.RegisterAt(app, "/metrics")
	// app.Use(prometheus.Middleware)
	migrations.Run()

	port := os.Getenv("PORT")
	fmt.Println("Server is running on port", port)
	log.Fatal(app.Listen(port))
}
