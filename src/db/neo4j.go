package db

import (
	"context"
	. "github.com/ArxivInsanity/graph-service/src/common"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/viper"
	"os"
)

const Neo4jContextKey string = "neo4jContext"
const Neo4jSessionKey string = "neo4jSession"

type Neo4jConnectionConfig struct {
	Uri        string
	Username   string
	Credential string
}

func getNeo4jConnectionConfig() Neo4jConnectionConfig {
	return Neo4jConnectionConfig{
		Uri:        viper.Get("neo4j.connectionUri").(string),
		Username:   viper.Get("neo4j.username").(string),
		Credential: os.Getenv("NEO4J_CRED"), // read credential from os env directly
	}
}

// GetNeo4jContextAndSession Exposes neo4j context and session created with default pooling config
func GetNeo4jContextAndSession() (context.Context, neo4j.SessionWithContext) {
	neo4jConnectionConfig := getNeo4jConnectionConfig()
	neo4jContext := context.Background()

	neo4jDriver, err := neo4j.NewDriverWithContext(neo4jConnectionConfig.Uri,
		neo4j.BasicAuth(neo4jConnectionConfig.Username, neo4jConnectionConfig.Credential, ""))
	PanicOnClosureError(err, neo4jContext, neo4jDriver)

	neo4jSession := neo4jDriver.NewSession(neo4jContext, neo4j.SessionConfig{})
	PanicOnClosureError(err, neo4jContext, neo4jSession)

	return neo4jContext, neo4jSession
}

// InjectNeo4jContextAndSession Is a closure function that injects context and session
func InjectNeo4jContextAndSession(neo4jContext context.Context, neo4jSession neo4j.SessionWithContext) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// setting in gin context
		ginContext.Set(Neo4jContextKey, neo4jContext)
		ginContext.Set(Neo4jSessionKey, neo4jSession)
	}
}
