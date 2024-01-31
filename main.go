package main

import (
	"com.fractopus/fractopus-node/gql"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	go gql.Process()
	r := gin.Default()
	err := r.Run(":9081")
	if err != nil {
		log.Println(err)
	}
}
