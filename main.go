// @title Files API
// @version 1.0
// @description This is a server for uploading and downloading files.
// @contact.name Branden Morin
// @contact.email brandenmorin14@gmail.com
// @host localhost
// @BasePath /
package main

import (
	"fmt"

	"github.com/BrandenM-PM/go-rest-api/controllers"
	_ "github.com/BrandenM-PM/go-rest-api/docs"
	"github.com/BrandenM-PM/go-rest-api/initializers"
	"github.com/BrandenM-PM/go-rest-api/middleware"
	"github.com/BrandenM-PM/go-rest-api/migrations"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/handlebars"

	"log"
	"os"
)

func init() {
	initializers.LoadEnvVars("./")
	initializers.ConnectToSqliteDB()
}

func main() {
	engine := handlebars.New("./views", ".hbs")

	app := fiber.New(fiber.Config{
		ErrorHandler: initializers.CustomErrorHandler,
		Views:        engine,
	})

	middleware.FiberMiddleware(app)

	app.Static("/", "./public")

	app.Get("/", controllers.Index)
	app.Get("/ping", controllers.Ping)
	app.Get("/files", controllers.GetAllFiles)
	app.Get("/files/:name", controllers.GetFile)
	app.Get("/get-presigned", controllers.GetPresigned)

	app.Group("/admin", middleware.JWTProtected(), middleware.IsAdmin).Get("/files", controllers.GetAllFiles)

	app.Post("/store-file/:name", controllers.CreateFile)
	app.Delete("/files/:id", controllers.DeleteFile)

	app.Get("/swagger/*", swagger.HandlerDefault)

	migrations.Run()

	port := os.Getenv("PORT")
	fmt.Println("Server is running on port", port)
	log.Fatal(app.Listen(port))
}
