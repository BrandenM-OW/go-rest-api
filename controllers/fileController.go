package controllers

import (
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/BrandenM-PM/go-rest-api/initializers"
	"github.com/BrandenM-PM/go-rest-api/models"
	"github.com/gofiber/fiber/v2"
)

// CreateFile godoc
// @Summary Creates an File
// @Produce json
// @Param title formData string true "Title"
// @Param content formData string true "Content"
// @Success 200 {object} models.File
// @Router /files [post]
func CreateFile(c *fiber.Ctx) error {
	name := c.FormValue("name")
	contentType := c.FormValue("contentType")
	size := c.FormValue("size")

	var intSize int64
	if size != "" {
		intSize, _ = strconv.ParseInt(size, 10, 64)
	}

	new_file := models.File{
		Name:        name,
		ContentType: contentType,
		Size:        intSize,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Owner:       "Branden",
	}

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
	start, _ := strconv.Atoi(c.FormValue("start"))
	length, _ := strconv.Atoi(c.FormValue("length"))
	draw, _ := strconv.Atoi(c.FormValue("draw"))
	intStart := int(start)
	intLength := int(length)
	intDraw := int(draw)

	var count int64
	initializers.DB.Model(&models.File{}).Count(&count)

	var files []models.File
	initializers.DB.Offset(intStart).Limit(intLength).Find(&files)

	c.JSON(&fiber.Map{
		"data":  files,
		"count": count,
		"draw":  intDraw,
	})

	return nil
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
// @Router /get-presigned [get]
func GetPresigned(c *fiber.Ctx) error {

	name := c.Query("name")
	contentType := c.Query("type")
	bucketName := os.Getenv("BUCKET_NAME")
	expiration := time.Now().Add(5 * time.Minute)

	signedURL, err := storage.SignedURL(bucketName, name, &storage.SignedURLOptions{
		GoogleAccessID: os.Getenv("GOOGLE_ACCESS_ID"),
		PrivateKey:     []byte(os.Getenv("GOOGLE_PRIVATE_KEY")),
		Method:         "PUT",
		Expires:        expiration,
		ContentType:    contentType,
	})
	if err != nil {
		log.Fatalf("Failed to create signed URL: %v", err)
	}

	return c.JSON(&fiber.Map{
		"url":  signedURL,
		"type": contentType,
	})
}
