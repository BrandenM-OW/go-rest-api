package controllers

import (
    "encoding/json"
    "github.com/BrandenM-PM/go-rest-api/initializers"
    "github.com/BrandenM-PM/go-rest-api/models"
    "github.com/go-redis/redis/v8"
    "github.com/gofiber/fiber/v2"
    "strconv"
    "time"
    "os"
)

// CreateArticle godoc
// @Summary Creates an Article
// @Produce json
// @Param title formData string true "Title"
// @Param content formData string true "Content"
// @Success 200 {object} models.Article
// @Router /articles [post]
func CreateArticle(c *fiber.Ctx) error {
    article := models.Article{
        Title:   c.FormValue("title"),
        Content: c.FormValue("content"),
    }

    result := initializers.DB.Create(&article)
    if result.Error != nil {
        return result.Error
    }

    rdb := redis.NewClient(&redis.Options{
        Addr: os.Getenv("REDIS_HOST"),
        DB:   0,
    })
    // Invalidate the Redis cache for all articles
    err := rdb.Del(rdb.Context(), "all_articles").Err()
    if err != nil {
        return err
    }

    c.JSON(&article)
    return nil
}

// GetAllArticle godoc
// @Summary Retrieves all Articles
// @Produce json
// @Success 200 {object} []models.Article
// @Router /articles [get]
func GetAllArticles(c *fiber.Ctx) error {
    rdb := redis.NewClient(&redis.Options{
        Addr: os.Getenv("REDIS_HOST"),
        DB:   0,
    })
    // Check if the data is cached in Redis
    cachedData, err := rdb.Get(rdb.Context(), "all_articles").Result()
    if err != redis.Nil && err != nil {
        return err
    }
    if cachedData != "" {
        // Serve the cached data
        var articles []models.Article
        err = json.Unmarshal([]byte(cachedData), &articles)
        if err != nil {
            return err
        }
        return c.JSON(&articles)
    }

    // Fetch the data from your PostgreSQL database using GORM
    var articles []models.Article
    result := initializers.DB.Find(&articles)
    if result.Error != nil {
        return result.Error
    }

    // Store the data in the Redis cache
    jsonData, err := json.Marshal(&articles)
    if err != nil {
        return err
    }
    err = rdb.Set(rdb.Context(), "all_articles", string(jsonData), time.Minute).Err()
    if err != nil {
        return err
    }

    // Return the fetched data as JSON
    return c.JSON(&articles)
}

// GetArticle godoc
// @Summary Retrieves an Article based on given ID
// @Produce json
// @Param id path integer true "Article ID"
// @Success 200 {object} models.Article
// @Router /articles/{id} [get]
func GetArticle(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return err
    }

    rdb := redis.NewClient(&redis.Options{
        Addr: os.Getenv("REDIS_HOST"),
        DB:   0,
    })
    // Check if the data is cached in Redis
    cachedData, err := rdb.Get(rdb.Context(), "article_"+strconv.Itoa(id)).Result()
    if err != redis.Nil && err != nil {
        return err
    }
    if cachedData != "" {
        // Serve the cached data
        var article models.Article
        err = json.Unmarshal([]byte(cachedData), &article)
        if err != nil {
            return err
        }
        return c.JSON(&article)
    }

    // Fetch the data from your PostgreSQL database using GORM
    var article models.Article
    result := initializers.DB.First(&article, id)
    if result.Error != nil {
        return result.Error
    }

    // Store the data in the Redis cache
    jsonData, err := json.Marshal(&article)
    if err != nil {
        return err
    }
    err = rdb.Set(rdb.Context(), "article_"+strconv.Itoa(id), string(jsonData), time.Minute).Err()
    if err != nil {
        return err
    }

    // Return the fetched data as JSON
    return c.JSON(&article)
}

// UpdateArticle godoc
// @Summary Updates an Article based on given ID
// @Produce json
// @Param id path integer true "Article ID"
// @Param title formData string true "Title"
// @Param content formData string true "Content"
// @Success 200 {object} models.Article
// @Router /articles/{id} [patch]
func UpdateArticle(c *fiber.Ctx) error {
    var article models.Article
    result := initializers.DB.First(&article, c.Params("id"))
    if result.Error != nil {
        return result.Error
    }

    article.Title = c.FormValue("title")
    article.Content = c.FormValue("content")
    result = initializers.DB.Save(&article)
    if result.Error != nil {
        return result.Error
    }


    rdb := redis.NewClient(&redis.Options{
        Addr: os.Getenv("REDIS_HOST"),
        DB:   0,
    })
    err := rdb.Del(rdb.Context(), "articles").Err()
    if err != nil {
        return err
    }
    err2 := rdb.Del(rdb.Context(), "article_"+c.Params("id")).Err()
    if err2 != nil { return err2 }

    c.JSON(&article)
    return nil
}

// DeleteArticle godoc
// @Summary Deletes an Article based on given ID
// @Produce json
// @Param id path integer true "Article ID"
// @Success 200 {object} models.Article
// @Router /articles/{id} [delete]
func DeleteArticle(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return err
    }

    rdb := redis.NewClient(&redis.Options{
        Addr: os.Getenv("REDIS_HOST"),
        DB:   0,
    })
    // Check if the data is cached in Redis
    cachedData, err := rdb.Get(rdb.Context(), "article_"+strconv.Itoa(id)).Result()
    if err != redis.Nil && err != nil {
        return err
    }
    if cachedData != "" {
        // Invalidate the Redis cache for the deleted article
        err = rdb.Del(rdb.Context(), "article_"+strconv.Itoa(id)).Err()
        if err != nil {
            return err
        }
    }

    // Delete the article from the database
    result := initializers.DB.Delete(&models.Article{}, id)
    if result.Error != nil {
        return result.Error
    }

    // Return a success message
    return c.JSON(&fiber.Map{
        "message": "Article deleted successfully",
    })
}

