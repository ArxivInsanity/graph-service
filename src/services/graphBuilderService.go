package services

import (
	. "github.com/ArxivInsanity/graph-service/src/db"
	. "github.com/ArxivInsanity/graph-service/src/util"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
)

func GetGraph() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get ctx and session from gin context
		dbContext, dbSession := GetDBConnectionFromContext(ctx)
		result, err := dbSession.Run(dbContext, "MATCH (n) RETURN n", map[string]any{})
		paper, err := neo4j.CollectTWithContext[neo4j.Node](dbContext, result,
			// Extract the single record and transform it with a function
			func(record *neo4j.Record) (neo4j.Node, error) {
				// Extract the record value by the specified key
				// and map it to the specified generic type constraint
				paper, _, err := neo4j.GetRecordValue[neo4j.Node](record, "n")
				return paper, err
			})
		PanicOnErr(err)
		log.Printf("%v", paper)
	}
}
