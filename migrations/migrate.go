package migrations

import "github.com/BrandenM-PM/go-rest-api/initializers"
import "github.com/BrandenM-PM/go-rest-api/models"

func Run() {
    initializers.DB.AutoMigrate(&models.Article{})
}
