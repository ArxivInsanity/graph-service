package services

import (
	. "github.com/ArxivInsanity/graph-service/src/util"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
)

func IsSeedPaper() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbContext, dbSession := GetDBConnectionFromContext(ctx)
		paperId := ctx.Param("paperId")
		log.Printf("URL Params: %v", paperId)
		isSeedCypher := "RETURN EXISTS( (:SEED_PAPER)-[:SEED]-(:PAPER {title: $title})) as isSeed"
		cypherParam := map[string]any{
			"title": paperId,
		}
		result, err := dbSession.Run(dbContext, isSeedCypher, cypherParam)
		record, err := neo4j.CollectWithContext(dbContext, result, err)
		isSeed, _ := record[0].Get("isSeed")
		log.Printf("Paper: %s is Seed Paper: %v", paperId, isSeed)
		PanicOnErr(err)
		ctx.JSON(http.StatusOK, isSeed)
	}
}
