package controllers

import (
    "testing"
    "net/http/httptest"
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    _ "github.com/BrandenM-PM/go-rest-api/docs"
    "github.com/BrandenM-PM/go-rest-api/initializers"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
    "github.com/gofiber/template/handlebars"
    "strings"
    "time"

) 

func TestWebController(t *testing.T) {
    // test cases
    tests := []struct {
        description string
        route       string
        expected    int
    }{
        {
            description: "Index route",
            route:       "/",
            expected:    200,
        },
        {
            description: "Fail route",
            route:       "/fail",
            expected:    404,
        },
    }

    // Setup
    initializers.LoadEnvVars("../")
    initializers.ConnectToPostgresDB()

    engine := handlebars.New("../views", ".hbs")
    app := fiber.New(fiber.Config{
        ErrorHandler: initializers.CustomErrorHandler,
        Views:        engine,
    })
    app.Use(cors.New())
    app.Use(recover.New())
    app.Use(limiter.New(limiter.Config{
        Next: func(c *fiber.Ctx) bool {
            return strings.Contains(c.Path(), "/swagger") // don't limit swagger
        },
        Expiration: 10 * time.Second,
        Max:        300,
    }))
    app.Use(logger.New())
    app.Static("/", "../public")
    app.Static("/purecss", "../node_modules/purecss/build")

    app.Get("/", Index)

    // Test
    for _, test := range tests {
        req := httptest.NewRequest("GET", test.route, nil)
        res, err := app.Test(req)
        if err != nil { t.Fatal(err) }
        assert.Equal(t, test.expected, res.StatusCode, test.description)
    }
}
