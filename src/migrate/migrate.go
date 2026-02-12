package migrate

import (
	"log"

	"gorm.io/gorm"
)

func Migrate(DB *gorm.DB) {
	DB.AutoMigrate()

	log.Println("Database migrated successfully")
}
