package db_dao

import (
	"com.fractopus/fractopus-node/storage/model"
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
	err = gormDb.AutoMigrate(&model.ConfigParam{})
	err = gormDb.AutoMigrate(&model.OpusNode{})
	err = gormDb.AutoMigrate(&model.OpusStream{})
	if err != nil {
		return
	}
}

func SaveLatestCursor(cursor string) {
	param := model.ConfigParam{}
	gormDb.Where("name=?", "latestCursor").First(&param)
	if param.ID == 0 {
		param = model.ConfigParam{Name: "latestCursor", Value: cursor}
	} else {
		param.Value = cursor
	}
	gormDb.Save(&param)
}

func GetLatestCursor() string {
	param := model.ConfigParam{}
	err := gormDb.Where("name=?", "latestCursor").First(&param).Error
	if err != nil {
		log.Println(err)
		return ""
	}
	return param.Value
}

func CheckUriExist(uri string) bool {
	var count int64 = -1
	err := gormDb.Model(model.OpusNode{}).Where("uri=?", uri).Count(&count).Error
	if err != nil {
		log.Println(err)
		return false
	}
	return count > 0
}

func SaveUris(list []model.OpusNode) {
	err := gormDb.Save(&list).Error
	if err != nil {
		log.Println(err)
	}
}
