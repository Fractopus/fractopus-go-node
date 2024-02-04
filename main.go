package main

import (
	"com.fractopus/fractopus-node/storage"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	storage.RedisInit()
}

func main() {

	//go gql.Process()
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
