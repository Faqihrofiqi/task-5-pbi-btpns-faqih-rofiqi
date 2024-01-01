package database

import (
	"github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PG *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=Kopisusu1212 dbname=finalpro port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{}, &models.Photo{})

	PG = db

}
