package main

import (
	"com.fractopus/fractopus-node/storage"
	"com.fractopus/fractopus-node/storage/db_dao"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	db_dao.SqliteDbInit()
	storage.RedisInit()
}

func main() {

	//go gql.Process()

	//go db_dao.SaveMany()

	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	err := router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	if err != nil {
		return
	}
	err = router.Run(":9081")
	if err != nil {
		log.Println(err)
	}
}
