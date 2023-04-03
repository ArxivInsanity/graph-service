package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

func GetGraph() gin.HandlerFunc {
	return func(context *gin.Context) {
		log.Printf("Printing NEO4J : %s\n", viper.Get("neo4j.connectionUri"))
		log.Printf("request handler: %s\n", context.Request.Header)
	}
}
