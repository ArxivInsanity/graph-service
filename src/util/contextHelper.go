package util

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const Neo4jContextKey string = "neo4jContext"
const Neo4jSessionKey string = "neo4jSession"

func GetDBConnectionFromContext(ginCtx *gin.Context) (context.Context, neo4j.SessionWithContext) {
	dbContextInterface, _ := ginCtx.Get(Neo4jContextKey)
	dbSessionInterface, _ := ginCtx.Get(Neo4jSessionKey)
	dbContext, dbSession := dbContextInterface.(context.Context), dbSessionInterface.(neo4j.SessionWithContext)
	return dbContext, dbSession
}
