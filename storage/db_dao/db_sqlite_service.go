package db_dao

import (
	"com.fractopus/fractopus-node/storage/model"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var gormDb *gorm.DB

func SqliteDbInit() {
	db, err := gorm.Open(sqlite.Open("./fractopus.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	gormDb = db
	err = gormDb.AutoMigrate(&model.OpusUri{})
	err = gormDb.AutoMigrate(&model.OpusStream{})
	if err != nil {
		return
	}
}

func SaveMany() {
	log.Println("开始")
	for j := 0; j < 10000; j++ {
		var list []model.OpusUri
		for i := 1; i < 1000; i++ {
			temp := model.OpusUri{
				Uri: fmt.Sprintf("test %v", i)}
			list = append(list, temp)
		}
		log.Println("完成组装")
		gormDb.Save(&list)
		log.Println("完成")
	}
}
