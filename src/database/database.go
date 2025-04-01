package database

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"go_fiber_blogs/src/models"
	"go_fiber_blogs/src/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		utils.GetEnvConfig("DB_HOST", "localhost"),
		utils.GetEnvConfig("DB_USER", "postgres"),
		utils.GetEnvConfig("DB_PASSWORD", "postgres"),
		utils.GetEnvConfig("DB_NAME", "go_fiber_blogs"),
		utils.GetEnvConfig("DB_PORT", "5432"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&models.Blog{}, &models.User{})

	DB = Dbinstance{
		Db: db,
	}
}
