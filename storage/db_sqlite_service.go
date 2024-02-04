package storage

import (
	"com.fractopus/fractopus-node/storage/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func SqliteDbInit() {
	db, err := gorm.Open(sqlite.Open("./fractopus.db"), &gorm.Config{})
	log.Println(err)

	err = db.AutoMigrate(&model.Product{})
	if err != nil {
		return
	}
}
