package controllers

import (
	"github.com/BrandenM-PM/go-rest-api/initializers"
	"github.com/BrandenM-PM/go-rest-api/models"
	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	var files []models.File
	result := initializers.DB.Find(&files)
	if result.Error != nil {
		return result.Error
	}
	c.Render("index", fiber.Map{
		"Title": "Home Page",
		"Files": files,
	}, "layouts/main")
	return nil
}

func Ping(c *fiber.Ctx) error {
	return c.SendString("pong")
}
