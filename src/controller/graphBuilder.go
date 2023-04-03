package controller

import (
	. "github.com/ArxivInsanity/graph-service/src/db"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

func GetGraph() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := GetNeo4jSession()
		log.Printf("Session %v", session)
		log.Printf("Printing NEO4J : %s\n", viper.Get("neo4j.connectionUri"))
		log.Printf("request handler: %s\n", context.Param("paperId"))
	}
}
