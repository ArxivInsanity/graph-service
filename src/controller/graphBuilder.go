package controller

import (
	"context"
	. "github.com/ArxivInsanity/graph-service/src/common"
	. "github.com/ArxivInsanity/graph-service/src/db"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/viper"
	"log"
)

func GetGraph() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbContextInterface, _ := ctx.Get(Neo4jContextKey)
		dbSessionInterface, _ := ctx.Get(Neo4jSessionKey)
		dbContext, dbSession := dbContextInterface.(context.Context), dbSessionInterface.(neo4j.SessionWithContext)
		result, err := dbSession.Run(dbContext, "MATCH (n) RETURN n", map[string]any{})
		PanicOnClosureError(err, dbContext, dbSession)
		log.Printf("Result %v", result)
		log.Printf("Printing NEO4J : %s\n", viper.Get("neo4j.connectionUri"))
		log.Printf("request handler: %s\n", ctx.Param("paperId"))
	}
}
