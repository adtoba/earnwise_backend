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

	// Ensure array columns use Postgres text[] type; AutoMigrate won't alter types.
	DB.Exec("ALTER TABLE IF EXISTS expert_profiles ALTER COLUMN categories TYPE text[] USING categories::text[]")
	DB.Exec("ALTER TABLE IF EXISTS expert_profiles ALTER COLUMN faq TYPE text[] USING faq::text[]")

	log.Println("Database migrated successfully")
}
