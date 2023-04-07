package controllers

import (
    "github.com/BrandenM-PM/go-rest-api/initializers"
    "github.com/BrandenM-PM/go-rest-api/models"
    "github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
    // get all articles
    var articles []models.Article
    result := initializers.DB.Find(&articles)
    if result.Error != nil { return result.Error }
    c.Render("index", fiber.Map{
        "Title": "Home Page",
        "Articles": articles,
    }, "layouts/main")
    return nil
}
