package controllers

import (
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/BrandenM-PM/go-rest-api/initializers"
	"github.com/BrandenM-PM/go-rest-api/models"
	"github.com/gofiber/fiber/v2"
)

const (
	projectID = "nimble-charmer-403517"
)

// CreateFile godoc
// @Summary Creates an File
// @Produce json
// @Param title formData string true "Title"
// @Param content formData string true "Content"
// @Success 200 {object} models.File
// @Router /files [post]
func CreateFile(c *fiber.Ctx) error {
	new_file := models.File{}

	result := initializers.DB.Create(&new_file)
	if result.Error != nil {
		return result.Error
	}

	c.JSON(&new_file)
	return nil
}

// GetAllFiles godoc
// @Summary Retrieves all Files
// @Produce json
// @Success 200 {object} []models.File
// @Router /Files [get]
func GetAllFiles(c *fiber.Ctx) error {
	var files []models.File
	result := initializers.DB.Find(&files)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(&files)
}

// GetFile godoc
// @Summary Retrieves a File based on given Name
// @Produce json
// @Param name path string true "File Name"
// @Success 200 {object} models.File
// @Router /files/{name} [get]
func GetFile(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("name"))
	if err != nil {
		return err
	}

	var file models.File
	result := initializers.DB.First(&file, id)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(&file)
}

// DeleteFile godoc
// @Summary Deletes a File based on given ID
// @Produce json
// @Param id path integer true "File ID"
// @Success 200 {object} models.File
// @Router /files/{id} [delete]
func DeleteFile(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	result := initializers.DB.Delete(&models.File{}, id)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(&fiber.Map{
		"message": "File deleted successfully",
	})
}

// GetPresigned godoc
// @Summary Gets a presigned URL for uploading a file to GCS
// @Produce json
// @Success 200
// @Router /getpresigned [get]
func GetPresigned(c *fiber.Ctx) error {
	bucketName := os.Getenv("BUCKET_NAME")
	fileName := "test.txt"
	method := "PUT"
	expires := time.Now().Add(time.Second * 60)

	opts := &storage.SignedURLOptions{
		GoogleAccessID: os.Getenv("GOOGLE_ACCESS_ID"),
		PrivateKey:     []byte(os.Getenv("GOOGLE_PRIVATE_KEY")),
		Method:         method,
		Expires:        expires,
	}

	url, err := storage.SignedURL(bucketName, fileName, opts)
	if err != nil {
		return err
	}

	return c.JSON(&fiber.Map{
		"url": url,
	})
}
