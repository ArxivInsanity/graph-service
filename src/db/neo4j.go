package db

import (
	"context"
	. "github.com/ArxivInsanity/graph-service/src/common"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/viper"
	"log"
)

type Neo4jConnection struct {
	Uri        string
	Username   string
	Credential string
}

// global neo4j session
var neo4jSession neo4j.SessionWithContext = nil

func getNeo4jConnectionDetails() Neo4jConnection {
	log.Println(viper.AllKeys())
	return Neo4jConnection{
		Uri:        viper.Get("neo4j.connectionUri").(string),
		Username:   viper.Get("neo4j.username").(string),
		Credential: viper.Get("neo4j.credential").(string),
	}
}

// Init GetNeo4JConnection function to ne04j connection
func Init() {
	neo4jConnection := getNeo4jConnectionDetails()
	neo4jContext := context.Background()
	neo4jDriver, err := neo4j.NewDriverWithContext(neo4jConnection.Uri,
		neo4j.BasicAuth(neo4jConnection.Username, neo4jConnection.Credential, ""))
	PanicOnErr(err)
	defer PanicOnClosureError(neo4jContext, neo4jDriver)
	neo4jSession = neo4jDriver.NewSession(neo4jContext, neo4j.SessionConfig{})
}

// GetNeo4jSession Exposes neo4j session created with default pooling config
func GetNeo4jSession() neo4j.SessionWithContext {
	return neo4jSession
}
