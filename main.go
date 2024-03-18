package main

import (
	"com.fractopus/fractopus-node/services"
	"com.fractopus/fractopus-node/storage"
	"com.fractopus/fractopus-node/storage/db_dao"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	db_dao.SqliteDbInit()
	storage.RedisInit()
}

/*
• “p": "fractopus" / 分形章鱼协议
• "uri" / URI
• shrL:[{"uri":"xxx","shr":"0.1"}] 如果上游的uri没有上链过，在分润的时候，就跳过
*/
func main() {
	go services.ProcessOnChainedUri()
	go services.ProcessWaitOnChainUri()

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
