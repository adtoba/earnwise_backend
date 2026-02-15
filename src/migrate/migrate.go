package migrate

import (
	"log"

	"github.com/adtoba/earnwise_backend/src/models"
	"gorm.io/gorm"
)

func Migrate(DB *gorm.DB) {
	DB.AutoMigrate(
		&models.User{},
		&models.ExpertProfile{},
		&models.Category{},
		&models.Wallet{},
		&models.Post{},
		&models.Comment{},
		&models.Review{},
	)

	log.Println("Database migrated successfully")
}
