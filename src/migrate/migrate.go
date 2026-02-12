package migrate

import (
	"log"

	"github.com/adtoba/earnwise_backend/src/models"
	"gorm.io/gorm"
)

func Migrate(DB *gorm.DB) {
	DB.AutoMigrate(
		&models.User{},
	)

	log.Println("Database migrated successfully")
}
