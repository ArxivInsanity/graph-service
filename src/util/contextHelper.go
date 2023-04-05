package util

import (
	"context"
	. "github.com/ArxivInsanity/graph-service/src/db"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetDBConnectionFromContext(ginCtx *gin.Context) (context.Context, neo4j.SessionWithContext) {
	dbContextInterface, _ := ginCtx.Get(Neo4jContextKey)
	dbSessionInterface, _ := ginCtx.Get(Neo4jSessionKey)
	dbContext, dbSession := dbContextInterface.(context.Context), dbSessionInterface.(neo4j.SessionWithContext)
	return dbContext, dbSession
}
