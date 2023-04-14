package services

import (
	"log"
	"net/http"

	. "github.com/ArxivInsanity/graph-service/src/util"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func IsSeedPaperHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		paperId := ctx.Param("paperId")
		log.Printf("Seed Paper handler URL Params: %v", paperId)
		isSeed := IsSeedPaper(paperId, ctx)
		ctx.JSON(http.StatusOK, isSeed)
	}
}

func IsSeedPaper(paperId string, ctx *gin.Context) bool {
	dbContext, dbSession := GetDBConnectionFromContext(ctx)
	isSeedCypher := "RETURN EXISTS( (:SEED_PAPER)-[:SEED]-(:PAPER {paperId: $paperId})) as isSeed"
	cypherParam := map[string]any{
		"paperId": paperId,
	}
	result, err := dbSession.Run(dbContext, isSeedCypher, cypherParam)
	record, err := neo4j.CollectWithContext(dbContext, result, err)
	isSeed, _ := record[0].Get("isSeed")
	log.Printf("Paper: %s is Seed Paper: %v", paperId, isSeed)
	PanicOnErr(err)
	return isSeed.(bool)
}
